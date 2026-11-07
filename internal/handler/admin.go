package handler

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"fortune/internal/auth"
	"fortune/internal/store"
	"fortune/internal/webutil"
)

type ctxKey string

const kAdminUsername ctxKey = "adminUsername"

// AdminLogin 管理后台登录（登录接口有速率限制，密码走 bcrypt 校验）。
func (h *Handler) AdminLogin(w http.ResponseWriter, r *http.Request) {
	if !h.loginLimiter.Allow(clientIP(r)) {
		webutil.Error(w, http.StatusTooManyRequests, "登录尝试过于频繁，请稍后再试")
		return
	}
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := webutil.ReadJSON(r, &req); err != nil {
		webutil.Error(w, http.StatusBadRequest, "请求体格式错误")
		return
	}
	if req.Username == "" || req.Password == "" {
		webutil.Error(w, http.StatusBadRequest, "请填写账号与密码")
		return
	}
	var adm store.AdminUser
	if h.db.First(&adm, "username = ?", req.Username).Error != nil || !auth.CheckPassword(adm.PasswordHash, req.Password) {
		webutil.Error(w, http.StatusUnauthorized, "账号或密码错误")
		return
	}
	token, hash, err := auth.NewToken()
	if err != nil {
		webutil.Error(w, http.StatusInternalServerError, "生成会话失败")
		return
	}
	exp := time.Now().Add(time.Duration(h.cfg.Admin.SessionTTLHours) * time.Hour)
	h.db.Create(&store.AdminSession{AdminID: adm.ID, TokenHash: hash, ExpiresAt: exp})
	h.db.Where("expires_at < ?", time.Now()).Delete(&store.AdminSession{}) // 懒清理过期会话
	webutil.OK(w, map[string]any{"token": token, "username": req.Username, "expires_at": exp})
}

// adminAuth 管理后台鉴权中间件：Bearer token -> 会话哈希 -> 验期。
func (h *Handler) adminAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := ""
		if a := r.Header.Get("Authorization"); strings.HasPrefix(a, "Bearer ") {
			token = strings.TrimPrefix(a, "Bearer ")
		}
		if token == "" {
			webutil.Error(w, http.StatusUnauthorized, "未登录")
			return
		}
		var s store.AdminSession
		if h.db.First(&s, "token_hash = ?", auth.HashToken(token)).Error != nil {
			webutil.Error(w, http.StatusUnauthorized, "会话无效")
			return
		}
		if time.Now().After(s.ExpiresAt) {
			h.db.Delete(&s)
			webutil.Error(w, http.StatusUnauthorized, "会话已过期，请重新登录")
			return
		}
		ctx := context.WithValue(r.Context(), kAdminUsername, s.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) AdminLogout(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	h.db.Where("token_hash = ?", auth.HashToken(token)).Delete(&store.AdminSession{})
	webutil.OK(w, map[string]any{"ok": true})
}

func (h *Handler) AdminMe(w http.ResponseWriter, r *http.Request) {
	var adm store.AdminUser
	if h.db.First(&adm, r.Context().Value(kAdminUsername)).Error != nil {
		webutil.Error(w, http.StatusUnauthorized, "未登录")
		return
	}
	webutil.OK(w, map[string]any{"username": adm.Username})
}

// AdminDashboard 数据总览。
func (h *Handler) AdminDashboard(w http.ResponseWriter, _ *http.Request) {
	var users, vipUsers, orders, paidOrders, aiMsgs, posts, pendingPosts int64
	h.db.Model(&store.User{}).Count(&users)
	h.db.Model(&store.User{}).Where("vip_until > ?", time.Now()).Count(&vipUsers)
	h.db.Model(&store.Order{}).Count(&orders)
	h.db.Model(&store.Order{}).Where("status = ?", "paid").Count(&paidOrders)
	h.db.Model(&store.AiMessage{}).Count(&aiMsgs)
	h.db.Model(&store.CommunityPost{}).Count(&posts)
	h.db.Model(&store.CommunityPost{}).Where("status = ?", "pending").Count(&pendingPosts)
	var revenue int64
	h.db.Model(&store.Order{}).Where("status = ?", "paid").Select("COALESCE(SUM(amount),0)").Scan(&revenue)

	var recent []store.Order
	h.db.Order("id desc").Limit(10).Find(&recent)

	webutil.OK(w, map[string]any{
		"users":         users,
		"vip_users":     vipUsers,
		"orders":        orders,
		"paid_orders":   paidOrders,
		"revenue_cents": revenue,
		"ai_messages":   aiMsgs,
		"posts":         posts,
		"pending_posts": pendingPosts,
		"recent_orders": recent,
	})
}

// AdminUsers 会员列表（含搜索与分页）。
func (h *Handler) AdminUsers(w http.ResponseWriter, r *http.Request) {
	page, size := pageSize(r)
	q := h.db.Model(&store.User{})
	if kw := r.URL.Query().Get("q"); kw != "" {
		q = q.Where("nickname LIKE ? OR id LIKE ? OR open_id LIKE ?", "%"+kw+"%", "%"+kw+"%", "%"+kw+"%")
	}
	var total int64
	q.Count(&total)
	var list []store.User
	q.Order("created_at desc").Offset((page - 1) * size).Limit(size).Find(&list)
	webutil.OK(w, map[string]any{"users": list, "total": total, "page": page})
}

// AdminUpdateUser 会员管理：封禁/解封、调整次数、延长会员、调整成长值/等级。
func (h *Handler) AdminUpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var u store.User
	if h.db.First(&u, "id = ?", id).Error != nil {
		webutil.Error(w, http.StatusNotFound, "用户不存在")
		return
	}
	var req struct {
		Status       *string `json:"status"`
		AddCredits   *int    `json:"add_credits"`
		SetCredits   *int    `json:"set_credits"`
		AddVipDays   *int    `json:"add_vip_days"`
		Level        *int    `json:"level"`
		SetGrowth    *int    `json:"set_growth"`
		Note         string  `json:"note"`
	}
	if err := webutil.ReadJSON(r, &req); err != nil {
		webutil.Error(w, http.StatusBadRequest, "请求体格式错误")
		return
	}
	if req.Status != nil {
		if *req.Status != "active" && *req.Status != "banned" {
			webutil.Error(w, http.StatusBadRequest, "status 仅支持 active/banned")
			return
		}
		h.db.Model(&u).Update("status", *req.Status)
	}
	if req.AddCredits != nil && *req.AddCredits != 0 {
		h.db.Model(&u).Update("ai_credits", gorm.Expr("ai_credits + ?", *req.AddCredits))
		h.db.Create(&store.MembershipLog{UserID: id, Kind: "admin", Credits: *req.AddCredits, Note: "后台调整次数"})
	}
	if req.SetCredits != nil {
		h.db.Model(&u).Update("ai_credits", *req.SetCredits)
		h.db.Create(&store.MembershipLog{UserID: id, Kind: "admin", Credits: *req.SetCredits, Note: "后台设置次数"})
	}
	if req.AddVipDays != nil && *req.AddVipDays != 0 {
		base := time.Now()
		if u.VipUntil != nil && u.VipUntil.After(base) {
			base = *u.VipUntil
		}
		newUntil := base.AddDate(0, 0, *req.AddVipDays)
		h.db.Model(&u).Update("vip_until", newUntil)
		h.db.Create(&store.MembershipLog{UserID: id, Kind: "admin", Days: *req.AddVipDays, Note: "后台调整会员天数"})
	}
	if req.SetGrowth != nil {
		h.db.Model(&u).Update("growth", *req.SetGrowth)
		h.recalcLevel(id)
	}
	if req.Level != nil && *req.Level >= 1 && *req.Level <= 9 {
		h.db.Model(&u).Update("level", *req.Level)
	}
	note := sanitizeUserText(req.Note, 255)
	if note != "" {
		h.db.Create(&store.MembershipLog{UserID: id, Kind: "admin", Note: note})
	}
	var nu store.User
	h.db.First(&nu, "id = ?", id)
	webutil.OK(w, nu)
}

// AdminMemberships 会员流水。
func (h *Handler) AdminMemberships(w http.ResponseWriter, r *http.Request) {
	page, size := pageSize(r)
	q := h.db.Model(&store.MembershipLog{})
	if uid := r.URL.Query().Get("user_id"); uid != "" {
		q = q.Where("user_id = ?", uid)
	}
	var total int64
	q.Count(&total)
	var list []store.MembershipLog
	q.Order("id desc").Offset((page - 1) * size).Limit(size).Find(&list)
	webutil.OK(w, map[string]any{"logs": list, "total": total, "page": page})
}

// AdminOrders 支付/订单管理。
func (h *Handler) AdminOrders(w http.ResponseWriter, r *http.Request) {
	page, size := pageSize(r)
	q := h.db.Model(&store.Order{})
	if s := r.URL.Query().Get("status"); s != "" {
		q = q.Where("status = ?", s)
	}
	if t := r.URL.Query().Get("type"); t != "" {
		q = q.Where("type = ?", t)
	}
	var total int64
	q.Count(&total)
	var list []store.Order
	q.Order("id desc").Offset((page - 1) * size).Limit(size).Find(&list)
	webutil.OK(w, map[string]any{"orders": list, "total": total, "page": page})
}

// AdminConsultations 日常咨询（AI 对话）管理。
func (h *Handler) AdminConsultations(w http.ResponseWriter, r *http.Request) {
	page, size := pageSize(r)
	q := h.db.Model(&store.AiMessage{})
	if uid := r.URL.Query().Get("user_id"); uid != "" {
		q = q.Where("user_id = ?", uid)
	}
	var total int64
	q.Count(&total)
	var list []store.AiMessage
	q.Order("id desc").Offset((page - 1) * size).Limit(size).Find(&list)
	webutil.OK(w, map[string]any{"messages": list, "total": total, "page": page})
}

// AdminReplyConsult 管理员回复某条咨询（写入同会话，role=admin）。
func (h *Handler) AdminReplyConsult(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	mid, _ := strconv.Atoi(idStr)
	var m store.AiMessage
	if h.db.First(&m, "id = ?", mid).Error != nil {
		webutil.Error(w, http.StatusNotFound, "咨询不存在")
		return
	}
	var req struct {
		Content string `json:"content"`
	}
	if err := webutil.ReadJSON(r, &req); err != nil {
		webutil.Error(w, http.StatusBadRequest, "请求体格式错误")
		return
	}
	content := sanitizeUserText(req.Content, 2000)
	if content == "" {
		webutil.Error(w, http.StatusBadRequest, "回复不能为空")
		return
	}
	h.db.Create(&store.AiMessage{UserID: m.UserID, SessionID: m.SessionID, Role: "admin", Content: content})
	webutil.OK(w, map[string]any{"ok": true})
}

// AdminPosts 社区内容管理（全部状态）。
func (h *Handler) AdminPosts(w http.ResponseWriter, r *http.Request) {
	page, size := pageSize(r)
	q := h.db.Model(&store.CommunityPost{})
	if s := r.URL.Query().Get("status"); s != "" {
		q = q.Where("status = ?", s)
	}
	var total int64
	q.Count(&total)
	var list []store.CommunityPost
	q.Order("id desc").Offset((page - 1) * size).Limit(size).Find(&list)
	webutil.OK(w, map[string]any{"posts": list, "total": total, "page": page})
}

// AdminModeratePost 审核帖子：置顶/上架/下架（published / pending / hidden）。
func (h *Handler) AdminModeratePost(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	pid, _ := strconv.Atoi(idStr)
	var req struct {
		Status string `json:"status"`
	}
	if err := webutil.ReadJSON(r, &req); err != nil {
		webutil.Error(w, http.StatusBadRequest, "请求体格式错误")
		return
	}
	if req.Status != "published" && req.Status != "pending" && req.Status != "hidden" {
		webutil.Error(w, http.StatusBadRequest, "status 仅支持 published/pending/hidden")
		return
	}
	h.db.Model(&store.CommunityPost{}).Where("id = ?", pid).Update("status", req.Status)
	webutil.OK(w, map[string]any{"ok": true})
}

// pageSize 解析分页参数。
func pageSize(r *http.Request) (page, size int) {
	page, _ = strconv.Atoi(r.URL.Query().Get("page"))
	size, _ = strconv.Atoi(r.URL.Query().Get("size"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	return page, size
}