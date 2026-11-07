package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"fortune/internal/payment"
	"fortune/internal/store"
	"fortune/internal/webutil"
)

// placeOrder 统一创建订单并调支付通道，返回前端拉起支付所需信息。
func (h *Handler) placeOrder(o store.Order, desc, openid string, r *http.Request) (orderResponse, error) {
	if o.OrderNo == "" {
		o.OrderNo = newOrderNo()
	}
	pr, err := h.pay.CreateOrder(payment.OrderInfo{
		OrderNo:  o.OrderNo,
		Amount:   o.Amount,
		Desc:     desc,
		OpenID:   openid,
		ClientIP: clientIP(r),
	})
	if err != nil {
		return orderResponse{}, err
	}
	if pr != nil {
		o.Provider = pr.Provider
		if bs, e := json.Marshal(pr); e == nil {
			o.PayParams = string(bs)
		}
		if err := h.db.Create(&o).Error; err != nil {
			return orderResponse{}, err
		}
		return orderResponse{
			OrderNo:    o.OrderNo,
			Provider:   pr.Provider,
			Amount:     o.Amount,
			PayParams:  pr.PayParams,
			MockPayURL: pr.MockPayURL,
			PayURL:     pr.PayURL,
		}, nil
	}
	if err := h.db.Create(&o).Error; err != nil {
		return orderResponse{}, err
	}
	return orderResponse{OrderNo: o.OrderNo, Provider: o.Provider, Amount: o.Amount}, nil
}

// CreateOrder 为某次算命创建解锁订单。
func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req orderReq
	if err := webutil.ReadJSON(r, &req); err != nil || req.ReadingID == "" {
		webutil.Error(w, http.StatusBadRequest, "缺少 reading_id")
		return
	}
	var reading store.Reading
	if err := h.db.First(&reading, "id = ?", req.ReadingID).Error; err != nil {
		webutil.Error(w, http.StatusNotFound, "未找到该算命结果")
		return
	}
	if h.isPaid(reading.ID) {
		webutil.OK(w, orderResponse{AlreadyPaid: true, Amount: h.cfg.App.PriceCents, Provider: h.pay.Name()})
		return
	}
	if h.userIsVIP(r) {
		// 会员全站内容免解锁，不创建付费订单
		webutil.OK(w, orderResponse{AlreadyPaid: true, Amount: 0, Provider: h.pay.Name()})
		return
	}
	resp, err := h.placeOrder(store.Order{
		Type:      "unlock",
		UserID:    req.UserID,
		ReadingID: reading.ID,
		Scenario:  reading.Scenario,
		Amount:    h.cfg.App.PriceCents,
		Status:    "pending",
	}, "解锁「"+reading.Title+"」全文", req.OpenID, r)
	if err != nil {
		webutil.Error(w, http.StatusInternalServerError, "创建支付失败："+err.Error())
		return
	}
	webutil.OK(w, resp)
}

// GetOrder 查询订单状态（供前端轮询，例如扫码支付后确认）。
func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	orderNo := chi.URLParam(r, "orderNo")
	var o store.Order
	if err := h.db.First(&o, "order_no = ?", orderNo).Error; err != nil {
		webutil.Error(w, http.StatusNotFound, "订单不存在")
		return
	}
	webutil.OK(w, map[string]any{
		"order_no": o.OrderNo,
		"status":   o.Status,
		"paid":     o.Status == "paid",
		"amount":   o.Amount,
		"type":     o.Type,
	})
}

// MockConfirm 模拟支付确认：把订单置为已支付。
func (h *Handler) MockConfirm(w http.ResponseWriter, r *http.Request) {
	if h.pay.Name() != "mock" {
		webutil.Error(w, http.StatusForbidden, "模拟支付仅限开发环境")
		return
	}
	orderNo := chi.URLParam(r, "orderNo")
	if orderNo == "" {
		orderNo = r.URL.Query().Get("order_no")
	}
	if r.Method == http.MethodPost && orderNo == "" {
		var body struct {
			OrderNo string `json:"order_no"`
		}
		_ = webutil.ReadJSON(r, &body)
		orderNo = body.OrderNo
	}
	if orderNo == "" {
		webutil.Error(w, http.StatusBadRequest, "缺少 order_no")
		return
	}
	if err := h.markPaid(orderNo, "mock"); err != nil {
		webutil.Error(w, http.StatusNotFound, err.Error())
		return
	}
	webutil.OK(w, map[string]any{"status": "paid", "order_no": orderNo})
}

// WechatNotify 微信支付异步回调：验签、解密、更新订单。
func (h *Handler) WechatNotify(w http.ResponseWriter, r *http.Request) {
	notifier, ok := h.pay.(payment.Notifier)
	if !ok {
		webutil.Error(w, http.StatusNotFound, "支付通道不支持回调")
		return
	}
	body, _ := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	ts := r.Header.Get("Wechatpay-Timestamp")
	nonce := r.Header.Get("Wechatpay-Nonce")
	sig := r.Header.Get("Wechatpay-Signature")

	res, err := notifier.VerifyAndDecrypt(ts, nonce, sig, body)
	if err != nil {
		webutil.JSON(w, http.StatusBadRequest, map[string]string{"code": "FAIL", "message": err.Error()})
		return
	}
	if res.TradeState != "SUCCESS" {
		webutil.JSON(w, http.StatusOK, map[string]string{"code": "SUCCESS", "message": "已收到，非成功状态"})
		return
	}
	if err := h.markPaid(res.OutTradeNo, "wechat"); err != nil {
		webutil.JSON(w, http.StatusInternalServerError, map[string]string{"code": "FAIL", "message": err.Error()})
		return
	}
	webutil.JSON(w, http.StatusOK, map[string]string{"code": "SUCCESS", "message": "成功"})
}

// AlipayNotify 支付宝异步回调：验签、更新订单。
func (h *Handler) AlipayNotify(w http.ResponseWriter, r *http.Request) {
	notifier, ok := h.pay.(payment.AlipayNotifier)
	if !ok {
		webutil.Error(w, http.StatusNotFound, "支付通道不支持回调")
		return
	}
	if err := r.ParseForm(); err != nil {
		writeAlipayFail(w)
		return
	}
	res, err := notifier.VerifyNotify(r.PostForm)
	if err != nil {
		writeAlipayFail(w)
		return
	}
	if res.TradeStatus != "TRADE_SUCCESS" && res.TradeStatus != "TRADE_FINISHED" {
		writeAlipaySuccess(w)
		return
	}
	if err := h.markPaid(res.OutTradeNo, "alipay"); err != nil {
		writeAlipayFail(w)
		return
	}
	writeAlipaySuccess(w)
}

func writeAlipaySuccess(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte("success"))
}
func writeAlipayFail(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte("failure"))
}

// markPaid 把订单置为已支付并发放对应权益（幂等：已支付订单不重复发货）。
func (h *Handler) markPaid(orderNo, provider string) error {
	var o store.Order
	if err := h.db.First(&o, "order_no = ?", orderNo).Error; err != nil {
		return err
	}
	if o.Status == "paid" {
		return nil // 防重复回调
	}
	now := time.Now()
	updates := map[string]any{"status": "paid", "paid_at": &now, "provider": provider}
	if provider == "wechat" || provider == "alipay" {
		updates["wechat_txn_id"] = orderNo // 第三方交易号的兜底记录
	}
	if err := h.db.Model(&o).Updates(updates).Error; err != nil {
		return err
	}
	h.fulfill(&o)
	return nil
}

// fulfill 依据订单类型发放权益。
func (h *Handler) fulfill(o *store.Order) {
	switch o.Type {
	case "subscribe":
		h.grantSubscription(o.UserID, o.Amount)
	case "quota":
		h.addCredits(o.UserID, h.cfg.Membership.QuotaAdd, o.Amount)
	}
}

// grantSubscription 开通/续费会员：延长 30 天 + 赠送每月 AI 次数 + 成长值。
func (h *Handler) grantSubscription(userID string, amountCents int) {
	if userID == "" {
		return
	}
	var u store.User
	if h.db.First(&u, "id = ?", userID).Error != nil {
		return
	}
	now := time.Now()
	base := now
	if u.VipUntil != nil && u.VipUntil.After(now) {
		base = *u.VipUntil
	}
	newUntil := base.AddDate(0, 1, 0)
	quota := h.cfg.Membership.MonthlyQuota
	h.db.Model(&u).Updates(map[string]any{
		"vip_until":  newUntil,
		"ai_credits": gorm.Expr("ai_credits + ?", quota),
		"growth":     gorm.Expr("growth + ?", 100),
	})
	h.db.Create(&store.MembershipLog{
		UserID: userID, Kind: "subscribe", Amount: amountCents,
		Credits: quota, Days: 30, Note: "开通/续费会员（月付）",
	})
	h.recalcLevel(userID)
}

// addCredits 增加 AI 调用次数（充值）。
func (h *Handler) addCredits(userID string, credits, amountCents int) {
	if userID == "" {
		return
	}
	var u store.User
	if h.db.First(&u, "id = ?", userID).Error != nil {
		return
	}
	h.db.Model(&u).Update("ai_credits", gorm.Expr("ai_credits + ?", credits))
	h.db.Create(&store.MembershipLog{
		UserID: userID, Kind: "topup", Amount: amountCents,
		Credits: credits, Note: "AI 次数充值包",
	})
}

// levelFor 依据成长值计算等级（会员成长体系）。
func levelFor(growth int) int {
	switch {
	case growth >= 5000:
		return 9
	case growth >= 2000:
		return 8
	case growth >= 1000:
		return 7
	case growth >= 500:
		return 6
	case growth >= 250:
		return 5
	case growth >= 100:
		return 4
	case growth >= 50:
		return 3
	case growth >= 20:
		return 2
	default:
		return 1
	}
}

// recalcLevel 依成长值重算并更新用户等级。
func (h *Handler) recalcLevel(userID string) {
	var u store.User
	if h.db.First(&u, "id = ?", userID).Error != nil {
		return
	}
	if lv := levelFor(u.Growth); lv != u.Level {
		h.db.Model(&u).Update("level", lv)
	}
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return strings.TrimSpace(strings.Split(xff, ",")[0])
	}
	host := r.RemoteAddr
	if i := strings.LastIndex(host, ":"); i > 0 {
		return strings.TrimPrefix(host[:i], "[")
	}
	return host
}