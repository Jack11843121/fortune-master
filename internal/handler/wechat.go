package handler

import (
	"encoding/json"
	"net/http"
	"net/url"

	"fortune/internal/webutil"
)

// 微信登录（小程序 jscode2session）：用前端传来的 code 换取 openid。
func (h *Handler) WechatLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Code string `json:"code"`
	}
	if err := webutil.ReadJSON(r, &req); err != nil || req.Code == "" {
		webutil.Error(w, http.StatusBadRequest, "缺少登录 code")
		return
	}

	appid := h.cfg.Payment.Wechat.AppID
	secret := h.cfg.Payment.Wechat.AppSecret
	if appid == "" || secret == "" {
		webutil.Error(w, http.StatusBadRequest, "微信登录未配置（缺少 appid / app_secret）")
		return
	}

	api := "https://api.weixin.qq.com/sns/jscode2session" +
		"?appid=" + url.QueryEscape(appid) +
		"&secret=" + url.QueryEscape(secret) +
		"&js_code=" + url.QueryEscape(req.Code) +
		"&grant_type=authorization_code"

	resp, err := http.Get(api)
	if err != nil {
		webutil.Error(w, http.StatusBadGateway, "微信登录请求失败")
		return
	}
	defer resp.Body.Close()

	var lr struct {
		OpenID     string `json:"openid"`
		SessionKey string `json:"session_key"`
		UnionID    string `json:"unionid"`
		ErrCode    int    `json:"errcode"`
		ErrMsg     string `json:"errmsg"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&lr); err != nil {
		webutil.Error(w, http.StatusBadGateway, "微信登录响应解析失败")
		return
	}
	if lr.ErrCode != 0 {
		webutil.Error(w, http.StatusBadRequest, "微信登录失败："+lr.ErrMsg)
		return
	}
	webutil.OK(w, map[string]string{"openid": lr.OpenID, "session_key": lr.SessionKey})
}