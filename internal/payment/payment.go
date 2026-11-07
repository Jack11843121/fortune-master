package payment

import "fmt"

// OrderInfo 是创建支付所需的订单信息。
type OrderInfo struct {
	OrderNo  string
	Amount   int    // 单位：分
	Desc     string
	OpenID   string // JSAPI 支付需要（微信内），其余忽略
	ClientIP string
}

// PayResult 返回给前端、用于拉起支付。
type PayResult struct {
	Provider  string         `json:"provider"`
	PayParams map[string]any `json:"pay_params,omitempty"`
	// MockPayURL 仅模拟支付使用：前端访问该地址即可“完成支付”。
	// PayURL 为支付宝/微信 H5 的跳转地址（前端重定向）。
	MockPayURL string `json:"mock_pay_url,omitempty"`
	PayURL     string `json:"pay_url,omitempty"`
}

// Provider 是支付抽象，mock / wechat / alipay 都实现它。
type Provider interface {
	Name() string
	CreateOrder(o OrderInfo) (*PayResult, error)
}

// Notifier 用于处理支付异步回调（wechat/alipay 实现；mock 走同步确认，无需实现）。
type Notifier interface {
	VerifyAndDecrypt(timestamp, nonce, signature string, rawBody []byte) (*NotifyResult, error)
}

// Config 是支付模块的配置。
type Config struct {
	Provider string // mock | wechat | alipay
	BaseURL  string // 本服务对外地址（用于 H5 等）
	Wechat   WechatConfig
	Alipay   AlipayConfig
}

// New 依据配置创建支付 Provider；未知 provider 默认回退到 mock（便于本地跑通）。
func New(cfg Config) (Provider, error) {
	switch cfg.Provider {
	case "wechat":
		if cfg.Wechat.MchID == "" || cfg.Wechat.PrivateKeyPath == "" || cfg.Wechat.APIv3Key == "" {
			return nil, fmt.Errorf("微信支付配置不完整（需要 mchid / mch_private_key_path / api_v3_key）")
		}
		return newWechat(cfg)
	case "alipay":
		if cfg.Alipay.AppID == "" || cfg.Alipay.PrivateKeyPath == "" {
			return nil, fmt.Errorf("支付宝配置不完整（需要 app_id / private_key_path）")
		}
		return newAlipay(cfg)
	default:
		return NewMock(), nil
	}
}