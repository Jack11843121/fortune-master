// 支付宝支付：RSA2（SHA256withRSA）手写实现，无第三方 SDK。
// 采用「手机网站支付 alipay.trade.wap.pay」，返回支付宝网关跳转地址，前端重定向拉起。
package payment

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
)

// AlipayConfig 支付宝配置。
type AlipayConfig struct {
	AppID             string
	PrivateKeyPath    string // 应用私钥（PKCS8，-----BEGIN PRIVATE KEY-----）
	AlipayPublicKey   string // 支付宝公钥（-----BEGIN PUBLIC KEY-----）
	NotifyURL         string
	ReturnURL         string
	Sandbox           bool
}

type alipay struct {
	cfg        AlipayConfig
	baseURL    string
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

const (
	alipayGateway     = "https://openapi.alipay.com/gateway.do"
	alipaySandboxGW   = "https://openapi-sandbox.dl.alipaydev.com/gateway.do"
	alipayMethod      = "alipay.trade.wap.pay"
	alipayProductCode = "QUICK_WAP_WAY"
)

// AlipayNotifier 支付宝异步回调接口（表单验签）。
type AlipayNotifier interface {
	VerifyNotify(values url.Values) (*AlipayNotifyResult, error)
}

func newAlipay(c Config) (*alipay, error) {
	priv, err := loadRSAPrivateKey(c.Alipay.PrivateKeyPath)
	if err != nil {
		return nil, err
	}
	a := &alipay{cfg: c.Alipay, baseURL: c.BaseURL, privateKey: priv}
	if c.Alipay.AlipayPublicKey != "" {
		pub, err := loadRSAPublicKey(c.Alipay.AlipayPublicKey)
		if err != nil {
			return nil, err
		}
		a.publicKey = pub
	}
	return a, nil
}

func (a *alipay) Name() string { return "alipay" }

func (a *alipay) gateway() string {
	if a.cfg.Sandbox {
		return alipaySandboxGW
	}
	return alipayGateway
}

// CreateOrder 生成支付宝手机网站支付跳转地址。
func (a *alipay) CreateOrder(o OrderInfo) (*PayResult, error) {
	totalYuan := fmt.Sprintf("%.2f", float64(o.Amount)/100.0)
	biz, _ := json.Marshal(map[string]any{
		"out_trade_no": o.OrderNo,
		"total_amount": totalYuan,
		"subject":      o.Desc,
		"product_code": alipayProductCode,
	})

	params := map[string]string{
		"app_id":      a.cfg.AppID,
		"method":      alipayMethod,
		"format":      "JSON",
		"charset":     "utf-8",
		"sign_type":   "RSA2",
		"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
		"version":     "1.0",
		"notify_url":  a.cfg.NotifyURL,
		"return_url":  a.cfg.ReturnURL,
		"biz_content": string(biz),
	}
	sig, err := a.sign(params)
	if err != nil {
		return nil, fmt.Errorf("支付宝签名失败: %w", err)
	}
	params["sign"] = sig

	q := url.Values{}
	for k, v := range params {
		q.Set(k, v)
	}
	payURL := a.gateway() + "?" + q.Encode()

	return &PayResult{
		Provider: "alipay",
		PayParams: map[string]any{
			"amount": o.Amount,
			"desc":   o.Desc,
		},
		PayURL: payURL,
	}, nil
}

// AlipayNotifyResult 支付宝异步通知结果。
type AlipayNotifyResult struct {
	OutTradeNo  string
	TradeStatus string // TRADE_SUCCESS / TRADE_FINISHED / ...
	TotalAmount string
	AppID       string
}

// VerifyNotify 校验支付宝异步通知签名并解析交易结果。
func (a *alipay) VerifyNotify(values url.Values) (*AlipayNotifyResult, error) {
	if a.publicKey == nil {
		return nil, errors.New("未配置支付宝公钥，无法验签")
	}
	sign := values.Get("sign")
	if sign == "" {
		return nil, errors.New("缺少 sign")
	}
	// 重建待签内容（排除 sign 与 sign_type）
	params := map[string]string{}
	for k := range values {
		if k == "sign" || k == "sign_type" {
			continue
		}
		params[k] = values.Get(k)
	}
	ok, err := a.verify(params, sign)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("支付宝回调验签失败")
	}
	return &AlipayNotifyResult{
		OutTradeNo:  values.Get("out_trade_no"),
		TradeStatus: values.Get("trade_status"),
		TotalAmount: values.Get("total_amount"),
		AppID:       values.Get("app_id"),
	}, nil
}

func (a *alipay) sign(params map[string]string) (string, error) {
	return rsa2Sign(a.privateKey, buildSignStr(params))
}

// SignRSA2 用支付宝应用私钥对参数做 RSA2 签名（供开放接口复用，如 OAuth2 登录）。
func SignRSA2(privateKeyPath string, params map[string]string) (string, error) {
	key, err := loadRSAPrivateKey(privateKeyPath)
	if err != nil {
		return "", err
	}
	return rsa2Sign(key, buildSignStr(params))
}

func (a *alipay) verify(params map[string]string, signature string) (bool, error) {
	return rsa2Verify(a.publicKey, buildSignStr(params), signature)
}

// buildSignStr 按支付宝规则：ASCII 升序、key=value（value URL 编码）、& 连接。
func buildSignStr(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b strings.Builder
	for i, k := range keys {
		if i > 0 {
			b.WriteByte('&')
		}
		b.WriteString(k)
		b.WriteByte('=')
		b.WriteString(alipayEscape(params[k]))
	}
	return b.String()
}

// alipayEscape 支付宝的 URL 编码（空格 %20、* %2A，其余 QueryEscape）。
func alipayEscape(s string) string {
	esc := url.QueryEscape(s)
	esc = strings.ReplaceAll(esc, "+", "%20")
	esc = strings.ReplaceAll(esc, "*", "%2A")
	return esc
}

func rsa2Sign(key *rsa.PrivateKey, data string) (string, error) {
	h := sha256.Sum256([]byte(data))
	sig, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, h[:])
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(sig), nil
}

func rsa2Verify(key *rsa.PublicKey, data, signature string) (bool, error) {
	sig, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, err
	}
	h := sha256.Sum256([]byte(data))
	if err := rsa.VerifyPKCS1v15(key, crypto.SHA256, h[:], sig); err != nil {
		return false, nil
	}
	return true, nil
}

func loadRSAPrivateKey(path string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取支付宝私钥失败: %w", err)
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("支付宝私钥 PEM 解析失败")
	}
	if key, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
		if rk, ok := key.(*rsa.PrivateKey); ok {
			return rk, nil
		}
		return nil, errors.New("私钥不是 RSA 类型")
	}
	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}
	return nil, errors.New("支付宝私钥需为 PKCS8/PKCS1 格式")
}

func loadRSAPublicKey(path string) (*rsa.PublicKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取支付宝公钥失败: %w", err)
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("支付宝公钥 PEM 解析失败")
	}
	if key, err := x509.ParsePKIXPublicKey(block.Bytes); err == nil {
		if rk, ok := key.(*rsa.PublicKey); ok {
			return rk, nil
		}
		return nil, errors.New("公钥不是 RSA 类型")
	}
	if key, err := x509.ParsePKCS1PublicKey(block.Bytes); err == nil {
		return key, nil
	}
	return nil, errors.New("支付宝公钥需为 PKIX/PKCS1 格式")
}