package handler

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"fortune/internal/ai"
	"fortune/internal/config"
	"fortune/internal/fortune"
	"fortune/internal/payment"
	"fortune/internal/ratelimit"
	"fortune/internal/store"
	"fortune/internal/webutil"
)

// Handler 聚合各依赖，承载 HTTP 处理逻辑。
type Handler struct {
	db           *gorm.DB
	engine       *fortune.Engine
	pay          payment.Provider
	cfg          *config.Config
	ai           *ai.Client
	loginLimiter *ratelimit.Limiter
	aiLimiter    *ratelimit.Limiter
	postLimiter  *ratelimit.Limiter
}

func New(db *gorm.DB, engine *fortune.Engine, pay payment.Provider, cfg *config.Config) *Handler {
	return &Handler{
		db:           db,
		engine:       engine,
		pay:          pay,
		cfg:          cfg,
		ai:           ai.New(ai.Config{Enabled: cfg.AI.Enabled, BaseURL: cfg.AI.BaseURL, APIKey: cfg.AI.APIKey, Model: cfg.AI.Model}),
		loginLimiter: ratelimit.New(10, time.Minute), // 登录限流
		aiLimiter:    ratelimit.New(30, time.Minute), // AI 对话限流（另有会员额度控制）
		postLimiter:  ratelimit.New(10, time.Minute), // 社区发帖/评论限流（防灌水）
	}
}

// Routes 注册所有 API 路由。
func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Use(recoverMiddleware)
	r.Use(logMiddleware)
	r.Use(corsMiddleware)

	r.Route("/v1", func(r chi.Router) {
		r.Post("/readings", h.CreateReading)
		r.Get("/readings/{id}", h.GetReading)
		r.Post("/orders", h.CreateOrder)
		r.Get("/orders/{orderNo}", h.GetOrder)
		r.Post("/naming", h.Naming)
		r.Post("/wechat/login", h.WechatLogin)
		r.Get("/pay/mock/confirm", h.MockConfirm)
		r.Post("/pay/mock/confirm", h.MockConfirm)
		r.Post("/pay/notify/wechat", h.WechatNotify)
		r.Post("/pay/notify/alipay", h.AlipayNotify)
		r.Get("/meta", h.Meta)

		// 账户与登录
		r.Post("/auth/wechat", h.AuthWechat)
		r.Post("/auth/wechat/oauth", h.AuthWechatOAuth)
		r.Get("/auth/wechat/url", h.AuthWechatURL)
		r.Post("/auth/dev", h.AuthDev)
		r.Post("/auth/logout", h.AuthLogout)
		r.Get("/auth/profile", h.AuthProfile)
		r.Post("/auth/delete", h.AuthDelete)

		// 会员 / 用户
		r.Get("/me", h.Me)
		r.Get("/membership/plans", h.MembershipPlans)
		r.Post("/membership/subscribe", h.Subscribe)
		r.Post("/membership/quota", h.BuyQuota)

		// AI 算命师
		r.Post("/ai/chat", h.AIChat)
		r.Get("/ai/history", h.AIHistory)

		// 交流社区
		r.Get("/community/posts", h.ListPosts)
		r.Post("/community/posts", h.CreatePost)
		r.Get("/community/posts/{id}", h.GetPost)
		r.Post("/community/posts/{id}/like", h.LikePost)
		r.Post("/community/posts/{id}/comments", h.CreateComment)

		// 内容分享
		r.Post("/share", h.RecordShare)

		// 管理后台
		r.Route("/admin", func(r chi.Router) {
			r.Post("/login", h.AdminLogin)
			r.Group(func(r chi.Router) {
				r.Use(h.adminAuth)
				r.Get("/me", h.AdminMe)
				r.Post("/logout", h.AdminLogout)
				r.Get("/dashboard", h.AdminDashboard)
				r.Get("/users", h.AdminUsers)
				r.Patch("/users/{id}", h.AdminUpdateUser)
				r.Get("/memberships", h.AdminMemberships)
				r.Get("/orders", h.AdminOrders)
				r.Get("/consultations", h.AdminConsultations)
				r.Post("/consultations/{id}/reply", h.AdminReplyConsult)
				r.Get("/posts", h.AdminPosts)
				r.Patch("/posts/{id}", h.AdminModeratePost)
			})
		})
	})
	return r
}

// Meta 返回前端需要的全局配置（价格、场景选项、会员方案等）。
func (h *Handler) Meta(w http.ResponseWriter, _ *http.Request) {
	webutil.OK(w, map[string]any{
		"app_name":     h.cfg.App.Name,
		"price_cents":  h.cfg.App.PriceCents,
		"pay_provider": h.cfg.Payment.Provider,
		"roles":        fortune.Roles,
		"membership": map[string]any{
			"monthly_price_cents": h.cfg.Membership.MonthlyPriceCents,
			"monthly_quota":       h.cfg.Membership.MonthlyQuota,
			"quota_price_cents":   h.cfg.Membership.QuotaPriceCents,
			"quota_add":           h.cfg.Membership.QuotaAdd,
		},
		"ai_enabled": h.cfg.AI.Enabled,
	})
}

// ---------- 请求 / 响应 DTO ----------

type createReadingReq struct {
	Scenario string          `json:"scenario"`
	PersonA  fortune.Person  `json:"person_a"`
	PersonB  *fortune.Person `json:"person_b"`
	Role     string          `json:"role"`
}

type matchView struct {
	Score      fortune.Score    `json:"score"`
	Aspects    []fortune.Aspect `json:"aspects"`
	Conclusion string           `json:"conclusion,omitempty"`
	Suggestion string           `json:"suggestion,omitempty"`
}

type readingResponse struct {
	ID           string             `json:"id"`
	Scenario     string             `json:"scenario"`
	Title        string             `json:"title"`
	Locked       bool               `json:"locked"`
	PriceCents   int                `json:"price_cents"`
	Bazi         *fortune.Bazi      `json:"bazi,omitempty"`
	BaziB        *fortune.Bazi      `json:"bazi_b,omitempty"`
	Match        *matchView         `json:"match,omitempty"`
	Sections     []fortune.Section  `json:"sections"`
	LockedTitles []string           `json:"locked_titles"`
	LockedCount  int                `json:"locked_count"`
	Teaser       string             `json:"teaser,omitempty"`
}

type orderReq struct {
	ReadingID string `json:"reading_id"`
	OpenID    string `json:"openid"`  // JSAPI 支付需传微信 openid
	UserID    string `json:"user_id"` // 归属用户（会员订单时必填）
}

type orderResponse struct {
	OrderNo     string         `json:"order_no"`
	Provider    string         `json:"provider"`
	Amount      int            `json:"amount"`
	PayParams   map[string]any `json:"pay_params,omitempty"`
	MockPayURL  string         `json:"mock_pay_url,omitempty"`
	PayURL      string         `json:"pay_url,omitempty"` // 支付宝/微信 H5 跳转地址
	AlreadyPaid bool           `json:"already_paid"`
}

// ---------- 读取与创建 ----------

func (h *Handler) CreateReading(w http.ResponseWriter, r *http.Request) {
	var req createReadingReq
	if err := webutil.ReadJSON(r, &req); err != nil {
		webutil.Error(w, http.StatusBadRequest, "请求体格式错误")
		return
	}

	var (
		res *fortune.Result
		err error
	)
	switch req.Scenario {
	case "marriage":
		if req.PersonB == nil {
			webutil.Error(w, http.StatusBadRequest, "合婚需要填写两人（person_b）信息")
			return
		}
		res, err = h.engine.Marriage(req.PersonA, *req.PersonB)
	case "hiring":
		res, err = h.engine.Hiring(req.PersonA, req.PersonB, req.Role)
	case "wealth":
		res, err = h.engine.Wealth(req.PersonA)
	default:
		res, err = h.engine.Personal(req.PersonA)
	}
	if err != nil {
		webutil.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	id := uuid.NewString()
	payload, _ := json.Marshal(res)
	reading := store.Reading{
		ID:       id,
		Scenario: res.Scenario,
		Title:    res.Title,
		Payload:  string(payload),
	}
	if err := h.db.Create(&reading).Error; err != nil {
		webutil.Error(w, http.StatusInternalServerError, "保存结果失败")
		return
	}

	webutil.JSON(w, http.StatusCreated, h.buildResponse(id, res, h.userIsVIP(r)))
}

func (h *Handler) GetReading(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var reading store.Reading
	if err := h.db.First(&reading, "id = ?", id).Error; err != nil {
		webutil.Error(w, http.StatusNotFound, "未找到该算命结果")
		return
	}
	var res fortune.Result
	if err := json.Unmarshal([]byte(reading.Payload), &res); err != nil {
		webutil.Error(w, http.StatusInternalServerError, "结果数据损坏")
		return
	}
	// 免费/演示定位：全部报告免费开放（不设付费锁定），配合合作广告引流。
	paid := true
	webutil.OK(w, h.buildResponse(id, &res, paid))
}

func (h *Handler) isPaid(readingID string) bool {
	var n int64
	h.db.Model(&store.Order{}).Where("reading_id = ? AND status = ?", readingID, "paid").Count(&n)
	return n > 0
}

// buildResponse 依付费状态切分免费预览与完整内容。
func (h *Handler) buildResponse(id string, res *fortune.Result, paid bool) readingResponse {
	resp := readingResponse{
		ID:           id,
		Scenario:     res.Scenario,
		Title:        res.Title,
		Locked:       !paid,
		PriceCents:   h.cfg.App.PriceCents,
		Bazi:         res.Bazi,
		BaziB:        res.BaziB,
		Sections:     []fortune.Section{},
		LockedTitles: []string{},
	}

	if res.Match != nil {
		mv := matchView{Score: res.Match.Score}
		const freeAspects = 2
		if paid {
			mv.Aspects = res.Match.Aspects
			mv.Conclusion = res.Match.Conclusion
			mv.Suggestion = res.Match.Suggestion
		} else {
			n := min(freeAspects, len(res.Match.Aspects))
			mv.Aspects = res.Match.Aspects[:n]
			for _, a := range res.Match.Aspects[n:] {
				resp.LockedTitles = append(resp.LockedTitles, a.Title)
			}
			suggestionLabel := "化解建议"
			if res.Scenario == "hiring" {
				suggestionLabel = "用人建议"
			}
			resp.LockedTitles = append(resp.LockedTitles, "综合结论", suggestionLabel)
			resp.LockedCount = len(res.Match.Aspects) - n + 2
			if n < len(res.Match.Aspects) {
				resp.Teaser = "「" + res.Match.Aspects[n].Title + "」" + truncateRunes(res.Match.Aspects[n].Body, 40) + "……"
			}
		}
		resp.Match = &mv
		resp.Sections = res.Sections // 概览段始终免费
	} else {
		total := len(res.Sections)
		freeCount := int(math.Round(float64(total) * h.cfg.Unlock.FreePreviewRatio))
		if freeCount < 2 {
			freeCount = 2
		}
		if freeCount > total {
			freeCount = total
		}
		if paid {
			freeCount = total
		}
		resp.Sections = res.Sections[:freeCount]
		for _, s := range res.Sections[freeCount:] {
			resp.LockedTitles = append(resp.LockedTitles, s.Title)
		}
		resp.LockedCount = total - freeCount
		if !paid && freeCount < total {
			resp.Teaser = "「" + res.Sections[freeCount].Title + "」" + truncateRunes(res.Sections[freeCount].Body, 40) + "……"
		}
	}
	return resp
}

func truncateRunes(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n])
}

// ---------- 中间件 ----------

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s (%s)", r.Method, r.URL.Path, time.Since(start))
	})
}

func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("handler panic: %v", rec)
				webutil.Error(w, http.StatusInternalServerError, "服务异常，请稍后重试")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-User-Id, Authorization, Wechatpay-Signature, Wechatpay-Timestamp, Wechatpay-Nonce, Wechatpay-Serial")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// newOrderNo 生成订单号。
func newOrderNo() string {
	return "FT" + strings.ReplaceAll(uuid.NewString(), "-", "")[:24] + time.Now().Format("150405")
}