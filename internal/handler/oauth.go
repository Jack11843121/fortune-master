package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"

	"fortune/internal/auth"
	"fortune/internal/store"
	"fortune/internal/webutil"
)

// 账户与登录：微信小程序 / 微信公众号 OAuth2 / 支付宝 OAuth2 / 开发模式。
// 登录成功后签发 30 天会话令牌，前端以 Authorization: Bearer <token> 携带。

// currentUser 从请求中解析当前登录用户（仅令牌，不做 uid 降级）。
func (h *Handler) currentUser(r *http.Request) *store.User {
	uid := h.loggedUserID(r)
	if uid == "" {
		return nil
	}
	var u store.User
	if h.db.First(&u, "id = ?", uid).Error != nil {
		return nil
	}
	return &u
}

// loggedUserID 校验 Bearer 令牌，返回用户 ID（无令牌/过期返回空）。
func (h *Handler) loggedUserID(r *http.Request) string {
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	if token == "" || token == r.Header.Get("Authorization") {
		return ""
	}
	var s store.UserSession
	if h.db.First(&s, "token_hash = ?", auth.HashToken(token)).Error != nil {
		return ""
	}
	if time.Now().After(s.ExpiresAt) {
		h.db.Delete(&s)
		return ""
	}
	return s.UserID
}

// ensureUserProvider 按 provider+openid 查找或创建账户，并更新昵称/头像。
func (h *Handler) ensureUserProvider(openid, unionid, provider, nickname, avatar string) (*store.User, error) {
	var u store.User
	if openid != "" {
		if err := h.db.Where("provider = ? AND open_id = ?", provider, openid).First(&u).Error; err == nil {
			h.refreshProfile(&u, unionid, nickname, avatar)
			return &u, nil
		}
	}
	if unionid != "" {
		if err := h.db.Where("union_id = ?", unionid).First(&u).Error; err == nil {
			updates := map[string]any{"provider": provider}
			if u.OpenID == "" {
				updates["open_id"] = openid
			}
			h.db.Model(&u).Updates(updates)
			h.refreshProfile(&u, unionid, nickname, avatar)
			return &u, nil
		}
	}
	if nickname == "" {
		nickname = "命主"
	}
	u = store.User{ID: uuid.NewString(), OpenID: openid, UnionID: unionid, Provider: provider, Nickname: nickname, Avatar: avatar, Level: 1, Status: "active"}
	if err := h.db.Create(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (h *Handler) refreshProfile(u *store.User, unionid, nickname, avatar string) {
	updates := map[string]any{}
	if nickname != "" && nickname != u.Nickname {
		updates["nickname"] = nickname
	}
	if avatar != "" && avatar != u.Avatar {
		updates["avatar"] = avatar
	}
	if unionid != "" && u.UnionID == "" {
		updates["union_id"] = unionid
	}
	if len(updates) > 0 {
		h.db.Model(u).Updates(updates)
	}
}

// issueToken 签发会话并返回登录响应。
func (h *Handler) issueToken(w http.ResponseWriter, u *store.User) {
	token, hash, err := auth.NewToken()
	if err != nil {
		webutil.Error(w, http.StatusInternalServerError, "生成会话失败")
		return
	}
	exp := time.Now().Add(720 * time.Hour) // 30 天
	h.db.Create(&store.UserSession{UserID: u.ID, TokenHash: hash, ExpiresAt: exp})
	h.db.Where("expires_at < ?", time.Now()).Delete(&store.UserSession{})
	webutil.OK(w, map[string]any{"token": token, "user": h.userView(*u), "expires_at": exp})
}

// AuthWechat 小程序/公众号 code 登录（jscode2session，取 openid）。
func (h *Handler) AuthWechat(w http.ResponseWriter, r *http.Request) {
	var req struct{ Code string `json:"code"` }
	if err := webutil.ReadJSON(r, &req); err != nil || req.Code == "" {
		webutil.Error(w, http.StatusBadRequest, "缺少登录 code")
		return
	}
	appid, secret := h.cfg.Payment.Wechat.AppID, h.cfg.Payment.Wechat.AppSecret
	if appid == "" || secret == "" {
		webutil.Error(w, http.StatusBadRequest, "微信登录未配置")
		return
	}
	data, err := httpGet("https://api.weixin.qq.com/sns/jscode2session?appid=" + url.QueryEscape(appid) +
		"&secret=" + url.QueryEscape(secret) + "&js_code=" + url.QueryEscape(req.Code) + "&grant_type=authorization_code")
	if err != nil {
		webutil.Error(w, http.StatusBadGateway, "微信登录请求失败")
		return
	}
	var lr struct {
		OpenID  string `json:"openid"`
		UnionID string `json:"unionid"`
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
	json.Unmarshal(data, &lr)
	if lr.ErrCode != 0 {
		webutil.Error(w, http.StatusBadRequest, "微信登录失败："+lr.ErrMsg)
		return
	}
	u, err := h.ensureUserProvider(lr.OpenID, lr.UnionID, "wechat", "", "")
	if err != nil {
		webutil.Error(w, http.StatusInternalServerError, "创建账户失败")
		return
	}
	h.issueToken(w, u)
}

// AuthWechatOAuth 公众号网页授权登录（snsapi_userinfo，取昵称头像）。
func (h *Handler) AuthWechatOAuth(w http.ResponseWriter, r *http.Request) {
	var req struct{ Code string `json:"code"` }
	if err := webutil.ReadJSON(r, &req); err != nil || req.Code == "" {
		webutil.Error(w, http.StatusBadRequest, "缺少登录 code")
		return
	}
	appid, secret := h.cfg.Payment.Wechat.AppID, h.cfg.Payment.Wechat.AppSecret
	if appid == "" || secret == "" {
		webutil.Error(w, http.StatusBadRequest, "微信登录未配置")
		return
	}
	data, err := httpGet("https://api.weixin.qq.com/sns/oauth2/access_token?appid=" + url.QueryEscape(appid) +
		"&secret=" + url.QueryEscape(secret) + "&code=" + url.QueryEscape(req.Code) + "&grant_type=authorization_code")
	if err != nil {
		webutil.Error(w, http.StatusBadGateway, "微信授权请求失败")
		return
	}
	var at struct {
		AccessToken string `json:"access_token"`
		OpenID      string `json:"openid"`
		UnionID     string `json:"unionid"`
		ErrCode     int    `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
	}
	json.Unmarshal(data, &at)
	if at.ErrCode != 0 {
		webutil.Error(w, http.StatusBadRequest, "微信授权失败："+at.ErrMsg)
		return
	}
	nick, avatar := "", ""
	if at.AccessToken != "" {
		if info, ierr := httpGet("https://api.weixin.qq.com/sns/userinfo?access_token=" + url.QueryEscape(at.AccessToken) +
			"&openid=" + url.QueryEscape(at.OpenID) + "&lang=zh_CN"); ierr == nil {
			var ui struct {
				Nickname   string `json:"nickname"`
				HeadImgURL string `json:"headimgurl"`
				UnionID    string `json:"unionid"`
			}
			json.Unmarshal(info, &ui)
			nick, avatar = ui.Nickname, ui.HeadImgURL
			if at.UnionID == "" {
				at.UnionID = ui.UnionID
			}
		}
	}
	u, err := h.ensureUserProvider(at.OpenID, at.UnionID, "wechat", nick, avatar)
	if err != nil {
		webutil.Error(w, http.StatusInternalServerError, "创建账户失败")
		return
	}
	h.issueToken(w, u)
}

// AuthWechatURL 返回公众号网页授权跳转地址（scope=snsapi_userinfo）。
func (h *Handler) AuthWechatURL(w http.ResponseWriter, r *http.Request) {
	appid := h.cfg.Payment.Wechat.AppID
	if appid == "" {
		webutil.Error(w, http.StatusBadRequest, "微信登录未配置")
		return
	}
	redirect := r.URL.Query().Get("redirect")
	if redirect == "" {
		redirect = h.cfg.App.BaseURL
	}
	authURL := "https://open.weixin.qq.com/connect/oauth2/authorize?appid=" + url.QueryEscape(appid) +
		"&redirect_uri=" + url.QueryEscape(redirect) +
		"&response_type=code&scope=snsapi_userinfo&state=login#wechat_redirect"
	webutil.OK(w, map[string]string{"url": authURL})
}

// AuthDev 开发模式登录（仅在 mock 支付环境可用，便于本地测试）。
func (h *Handler) AuthDev(w http.ResponseWriter, r *http.Request) {
	if h.pay.Name() != "mock" {
		webutil.Error(w, http.StatusForbidden, "开发登录仅在本地/演示环境可用")
		return
	}
	var req struct {
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
	}
	_ = webutil.ReadJSON(r, &req)
	nick := strings.TrimSpace(req.Nickname)
	if nick == "" {
		nick = "测试用户" + uuid.NewString()[:4]
	}
	u, err := h.ensureUserProvider("dev-"+uuid.NewString()[:16], "", "dev", nick, req.Avatar)
	if err != nil {
		webutil.Error(w, http.StatusInternalServerError, "创建账户失败")
		return
	}
	h.issueToken(w, u)
}

// AuthLogout 退出登录。
func (h *Handler) AuthLogout(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	h.db.Where("token_hash = ?", auth.HashToken(token)).Delete(&store.UserSession{})
	webutil.OK(w, map[string]any{"ok": true})
}

// AuthProfile 返回当前登录用户资料。
func (h *Handler) AuthProfile(w http.ResponseWriter, r *http.Request) {
	u := h.currentUser(r)
	if u == nil {
		webutil.Error(w, http.StatusUnauthorized, "未登录")
		return
	}
	webutil.OK(w, h.userView(*u))
}

// AuthDelete 账号注销：删除账户及其关联数据（PIPL 个人信息删权）。
func (h *Handler) AuthDelete(w http.ResponseWriter, r *http.Request) {
	u := h.currentUser(r)
	if u == nil {
		webutil.Error(w, http.StatusUnauthorized, "未登录")
		return
	}
	uid := u.ID
	// 删除关联数据
	h.db.Where("user_id = ?", uid).Delete(&store.UserSession{})
	h.db.Where("user_id = ?", uid).Delete(&store.MembershipLog{})
	h.db.Where("user_id = ?", uid).Delete(&store.AiMessage{})
	h.db.Where("user_id = ?", uid).Delete(&store.ShareRecord{})
	h.db.Where("user_id = ?", uid).Delete(&store.CommunityPost{})
	h.db.Where("user_id = ?", uid).Delete(&store.CommunityComment{})
	h.db.Delete(&store.User{ID: uid})
	webutil.OK(w, map[string]any{"ok": true, "message": "账号已注销"})
}

func httpGet(uri string) ([]byte, error) {
	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(io.LimitReader(resp.Body, 1<<20))
}