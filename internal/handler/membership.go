package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	"fortune/internal/store"
	"fortune/internal/webutil"
)

// userID 从前端拿当前用户 uid（X-User-Id 头；本地会话，未绑定微信时作为轻量身份）。
func (h *Handler) userID(r *http.Request) string {
	return r.Header.Get("X-User-Id")
}

// ensureUser 读取或创建用户。
func (h *Handler) ensureUser(uid, openid string) *store.User {
	if uid == "" {
		uid = uuid.NewString()
	}
	var u store.User
	if err := h.db.First(&u, "id = ?", uid).Error; err != nil {
		u = store.User{ID: uid, OpenID: openid, Nickname: "命主", Level: 1, Status: "active"}
		h.db.Create(&u)
	}
	if openid != "" && u.OpenID == "" {
		h.db.Model(&u).Update("open_id", openid)
	}
	return &u
}

// userIsVIP 判断当前请求用户是否为有效会员（会员全站内容免解锁）。
func (h *Handler) userIsVIP(r *http.Request) bool {
	u := h.currentUser(r)
	return u != nil && u.IsVIP()
}

// userView 用户可见的会员状态视图（含资料与会员数据）。
func (h *Handler) userView(u store.User) map[string]any {
	vip := u.VipUntil != nil && u.VipUntil.After(time.Now())
	return map[string]any{
		"id":         u.ID,
		"nickname":   u.Nickname,
		"avatar":     u.Avatar,
		"provider":   u.Provider,
		"vip":        vip,
		"vip_until":  u.VipUntil,
		"ai_credits": u.AICredits,
		"growth":     u.Growth,
		"level":      u.Level,
		"status":     u.Status,
	}
}

// Me 返回当前登录用户的会员状态与 AI 用量。
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	u := h.currentUser(r)
	if u == nil {
		webutil.OK(w, map[string]any{"id": "", "nickname": "", "avatar": "", "vip": false, "ai_credits": 0, "growth": 0, "level": 1, "status": "active"})
		return
	}
	if u.Status == "banned" {
		webutil.Error(w, http.StatusForbidden, "该账号已被限制")
		return
	}
	webutil.OK(w, h.userView(*u))
}

// MembershipPlans 返回会员与充值的价位方案。
func (h *Handler) MembershipPlans(w http.ResponseWriter, _ *http.Request) {
	webutil.OK(w, map[string]any{
		"monthly_price_cents": h.cfg.Membership.MonthlyPriceCents,
		"monthly_quota":       h.cfg.Membership.MonthlyQuota,
		"quota_price_cents":   h.cfg.Membership.QuotaPriceCents,
		"quota_add":           h.cfg.Membership.QuotaAdd,
	})
}

// Subscribe 创建会员订阅订单（需登录）。
func (h *Handler) Subscribe(w http.ResponseWriter, r *http.Request) {
	u := h.currentUser(r)
	if u == nil {
		webutil.Error(w, http.StatusUnauthorized, "请先登录后再开通会员")
		return
	}
	resp, err := h.placeOrder(store.Order{
		Type:   "subscribe",
		UserID: u.ID,
		Amount: h.cfg.Membership.MonthlyPriceCents,
		Status: "pending",
	}, "算命大师 VIP 会员（月付 · 赠 1000 次 AI）", "", r)
	if err != nil {
		webutil.Error(w, http.StatusInternalServerError, "创建订阅订单失败："+err.Error())
		return
	}
	webutil.OK(w, resp)
}

// BuyQuota 购买 AI 次数充值包（需登录）。
func (h *Handler) BuyQuota(w http.ResponseWriter, r *http.Request) {
	u := h.currentUser(r)
	if u == nil {
		webutil.Error(w, http.StatusUnauthorized, "请先登录后再充值")
		return
	}
	resp, err := h.placeOrder(store.Order{
		Type:   "quota",
		UserID: u.ID,
		Amount: h.cfg.Membership.QuotaPriceCents,
		Status: "pending",
	}, fmt.Sprintf("AI 算命师次数充值包（%d 次）", h.cfg.Membership.QuotaAdd), "", r)
	if err != nil {
		webutil.Error(w, http.StatusInternalServerError, "创建充值订单失败："+err.Error())
		return
	}
	webutil.OK(w, resp)
}