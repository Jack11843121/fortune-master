package payment

import (
	"bytes"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// 微信支付 API v3 —— 手写实现（签名 / 验签 / AES-256-GCM 解密），
// 仅依赖标准库，便于离线做签名加解密单测；拿到商户号与证书后即可启用。

const wechatAPIBase = "https://api.mch.weixin.qq.com"

// WechatConfig 微信支付配置。
type WechatConfig struct {
	AppID            string
	MchID            string
	MchSerialNo      string
	PrivateKeyPath   string // apiclient_key.pem（PKCS#8）
	PlatformCertPath string // 平台证书 wechatpay_xxx.pem（验签用；为空则跳过验签仅解密）
	APIv3Key         string // 32 位 APIv3 密钥
	NotifyURL        string
	PayType          string // jsapi | native | h5
}

type wechat struct {
	cfg        WechatConfig
	baseURL    string
	privateKey *rsa.PrivateKey
	platform   *rsa.PublicKey
	http       *http.Client
}

func newWechat(c Config) (*wechat, error) {
	w := &wechat{
		cfg:     c.Wechat,
		baseURL: c.BaseURL,
		http:    &http.Client{Timeout: 15 * time.Second},
	}
	priv, err := loadPrivateKey(c.Wechat.PrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("加载商户私钥失败: %w", err)
	}
	w.privateKey = priv
	if c.Wechat.PlatformCertPath != "" {
		pub, err := loadPublicKey(c.Wechat.PlatformCertPath)
		if err != nil {
			return nil, fmt.Errorf("加载平台证书失败: %w", err)
		}
		w.platform = pub
	}
	if c.Wechat.NotifyURL == "" {
		w.cfg.NotifyURL = strings.TrimRight(c.BaseURL, "/") + "/api/v1/pay/notify/wechat"
	}
	if w.cfg.PayType == "" {
		w.cfg.PayType = "native"
	}
	return w, nil
}

func (w *wechat) Name() string { return "wechat" }

func (w *wechat) CreateOrder(o OrderInfo) (*PayResult, error) {
	switch w.cfg.PayType {
	case "jsapi":
		if o.OpenID == "" {
			return nil, fmt.Errorf("JSAPI 支付需要 openid")
		}
		return w.prepayJSAPI(o)
	case "h5":
		if o.ClientIP == "" {
			o.ClientIP = "127.0.0.1"
		}
		return w.prepayH5(o)
	default:
		return w.prepayNative(o)
	}
}

// ---------- 三种下单 ----------

func (w *wechat) prepayNative(o OrderInfo) (*PayResult, error) {
	body := map[string]any{
		"appid":       w.cfg.AppID,
		"mchid":       w.cfg.MchID,
		"description": o.Desc,
		"out_trade_no": o.OrderNo,
		"notify_url":  w.cfg.NotifyURL,
		"amount":      map[string]any{"total": o.Amount, "currency": "CNY"},
	}
	var resp struct {
		CodeURL string `json:"code_url"`
	}
	if err := w.request("POST", "/v3/pay/transactions/native", body, &resp); err != nil {
		return nil, err
	}
	return &PayResult{Provider: "wechat", PayParams: map[string]any{"code_url": resp.CodeURL}}, nil
}

func (w *wechat) prepayH5(o OrderInfo) (*PayResult, error) {
	body := map[string]any{
		"appid":       w.cfg.AppID,
		"mchid":       w.cfg.MchID,
		"description": o.Desc,
		"out_trade_no": o.OrderNo,
		"notify_url":  w.cfg.NotifyURL,
		"amount":      map[string]any{"total": o.Amount, "currency": "CNY"},
		"scene_info": map[string]any{
			"payer_client_ip": o.ClientIP,
			"h5_info":         map[string]any{"type": "Wap"},
		},
	}
	var resp struct {
		H5URL string `json:"h5_url"`
	}
	if err := w.request("POST", "/v3/pay/transactions/h5", body, &resp); err != nil {
		return nil, err
	}
	return &PayResult{Provider: "wechat", PayParams: map[string]any{"h5_url": resp.H5URL}}, nil
}

func (w *wechat) prepayJSAPI(o OrderInfo) (*PayResult, error) {
	body := map[string]any{
		"appid":       w.cfg.AppID,
		"mchid":       w.cfg.MchID,
		"description": o.Desc,
		"out_trade_no": o.OrderNo,
		"notify_url":  w.cfg.NotifyURL,
		"amount":      map[string]any{"total": o.Amount, "currency": "CNY"},
		"payer":       map[string]any{"openid": o.OpenID},
	}
	var resp struct {
		PrepayID string `json:"prepay_id"`
	}
	if err := w.request("POST", "/v3/pay/transactions/jsapi", body, &resp); err != nil {
		return nil, err
	}
	// 拼接前端 WeixinJSBridge 调起支付所需参数
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	nonce := randNonce()
	pkg := "prepay_id=" + resp.PrepayID
	paySign, err := signRSA(w.privateKey, []byte(w.cfg.AppID+"\n"+ts+"\n"+nonce+"\n"+pkg+"\n"))
	if err != nil {
		return nil, err
	}
	return &PayResult{Provider: "wechat", PayParams: map[string]any{
		"appId":     w.cfg.AppID,
		"timeStamp": ts,
		"nonceStr":  nonce,
		"package":   pkg,
		"signType":  "RSA",
		"paySign":   paySign,
	}}, nil
}

// ---------- 通用 HTTP（带签名） ----------

func (w *wechat) request(method, path string, body any, target any) error {
	bs, _ := json.Marshal(body)
	req, err := http.NewRequest(method, wechatAPIBase+path, bytes.NewReader(bs))
	if err != nil {
		return err
	}
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	nonce := randNonce()
	msg := method + "\n" + path + "\n" + ts + "\n" + nonce + "\n" + string(bs) + "\n"
	sig, err := signRSA(w.privateKey, []byte(msg))
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "fortune/1.0")
	req.Header.Set("Authorization",
		fmt.Sprintf(`WECHATPAY2-SHA256-RSA2048 mchid="%s",nonce_str="%s",signature="%s",timestamp="%s",serial_no="%s"`,
			w.cfg.MchID, nonce, sig, ts, w.cfg.MchSerialNo))

	resp, err := w.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("微信支付请求失败(%d): %s", resp.StatusCode, string(data))
	}
	if target != nil && len(data) > 0 {
		if err := json.Unmarshal(data, target); err != nil {
			return err
		}
	}
	return nil
}

// ---------- 回调验签与解密 ----------

// NotifyResource 微信回调中的 resource 字段。
type NotifyResource struct {
	OriginalType   string `json:"original_type"`
	Algorithm      string `json:"algorithm"`
	Ciphertext     string `json:"ciphertext"`
	AssociatedData string `json:"associated_data"`
	Nonce          string `json:"nonce"`
}

// NotifyResult 解密后的回调内容。
type NotifyResult struct {
	OutTradeNo string `json:"out_trade_no"`
	TradeState string `json:"trade_state"` // SUCCESS / NOTPAY / CLOSED ...
	Amount     struct {
		Total int `json:"total"`
	} `json:"amount"`
}

// VerifyAndDecrypt 验签并解密微信支付回调；返回交易结果。
// 若配置了平台证书则先验签，否则仅解密（生产务必配置平台证书）。
func (w *wechat) VerifyAndDecrypt(timestamp, nonce, signature string, rawBody []byte) (*NotifyResult, error) {
	if w.platform != nil {
		msg := timestamp + "\n" + nonce + "\n" + string(rawBody) + "\n"
		if !verifyRSA(w.platform, []byte(msg), signature) {
			return nil, fmt.Errorf("微信回调验签失败")
		}
	}
	var payload struct {
		Resource NotifyResource `json:"resource"`
	}
	if err := json.Unmarshal(rawBody, &payload); err != nil {
		return nil, err
	}
	plain, err := gcmDecrypt([]byte(w.cfg.APIv3Key), payload.Resource.Nonce, payload.Resource.Ciphertext, payload.Resource.AssociatedData)
	if err != nil {
		return nil, fmt.Errorf("回调解密失败: %w", err)
	}
	var res NotifyResult
	if err := json.Unmarshal(plain, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

// ---------- 加解密与签名 ----------

func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	bs, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(bs)
	if block == nil {
		return nil, fmt.Errorf("无法解析 PEM 私钥")
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("私钥不是 RSA 类型")
	}
	return rsaKey, nil
}

func loadPublicKey(path string) (*rsa.PublicKey, error) {
	bs, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(bs)
	if block == nil {
		return nil, fmt.Errorf("无法解析 PEM 证书")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	return cert.PublicKey.(*rsa.PublicKey), nil
}

func signRSA(key *rsa.PrivateKey, msg []byte) (string, error) {
	digest := sha256.Sum256(msg)
	sig, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, digest[:])
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(sig), nil
}

func verifyRSA(key *rsa.PublicKey, msg []byte, sigB64 string) bool {
	sig, err := base64.StdEncoding.DecodeString(sigB64)
	if err != nil {
		return false
	}
	digest := sha256.Sum256(msg)
	err = rsa.VerifyPKCS1v15(key, crypto.SHA256, digest[:], sig)
	return err == nil
}

func gcmDecrypt(key []byte, nonce, ciphertextB64, associatedData string) ([]byte, error) {
	ct, err := base64.StdEncoding.DecodeString(ciphertextB64)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return gcm.Open(nil, []byte(nonce), ct, []byte(associatedData))
}

func randNonce() string {
	const chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return strconv.FormatInt(time.Now().UnixNano(), 10)
	}
	for i := range b {
		b[i] = chars[int(b[i])%len(chars)]
	}
	return string(b)
}