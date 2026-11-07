package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 是全局配置的顶层结构。
type Config struct {
	App        App        `yaml:"app"`
	Server     Server     `yaml:"server"`
	Database   Database   `yaml:"database"`
	Payment    Payment    `yaml:"payment"`
	Unlock     Unlock     `yaml:"unlock"`
	Membership Membership `yaml:"membership"`
	AI         AI         `yaml:"ai"`
	Admin      Admin      `yaml:"admin"`
}

type App struct {
	Name       string `yaml:"name"`
	BaseURL    string `yaml:"base_url"`
	PriceCents int    `yaml:"price_cents"`
}

type Server struct {
	Addr    string `yaml:"addr"`
	TLSCert string `yaml:"tls_cert"` // 可选：HTTPS 证书文件路径（配了 tls_key 后直接走 https）
	TLSKey  string `yaml:"tls_key"`  // 可选：HTTPS 私钥文件路径
}

type Database struct {
	Driver string `yaml:"driver"`    // sqlite | mysql
	DSN    string `yaml:"dsn"`
}

type Payment struct {
	Provider string `yaml:"provider"` // mock | wechat | alipay
	Mock     Mock   `yaml:"mock"`
	Wechat   Wechat `yaml:"wechat"`
	Alipay   Alipay `yaml:"alipay"`
}

type Mock struct {
	Enabled bool `yaml:"enabled"`
}

type Wechat struct {
	AppID            string `yaml:"appid"`
	AppSecret        string `yaml:"app_secret"`
	MchID            string `yaml:"mchid"`
	MchSerialNo      string `yaml:"mch_serial_no"`
	MchPrivateKey    string `yaml:"mch_private_key_path"`
	PlatformCertPath string `yaml:"platform_cert_path"` // 平台证书（验签用，可为空则跳过验签）
	APIv3Key         string `yaml:"api_v3_key"`
	NotifyURL        string `yaml:"notify_url"`
	PayType          string `yaml:"pay_type"` // jsapi | native | h5
}

type Alipay struct {
	AppID           string `yaml:"app_id"`
	PrivateKeyPath  string `yaml:"private_key_path"`         // 应用私钥（PKCS8）
	AlipayPublicKey string `yaml:"alipay_public_key_path"`   // 支付宝公钥
	NotifyURL       string `yaml:"notify_url"`               // 异步回调
	ReturnURL       string `yaml:"return_url"`               // 同步跳转
	Sandbox         bool   `yaml:"sandbox"`                  // 沙箱环境
}

type Unlock struct {
	FreePreviewRatio float64 `yaml:"free_preview_ratio"`
}

// Membership 会员与 AI 用量配置。
type Membership struct {
	MonthlyPriceCents int `yaml:"monthly_price_cents"` // 月费（分），默认 1000 = 10 元
	MonthlyQuota      int `yaml:"monthly_quota"`       // 每月赠送 AI 调用次数，默认 1000
	QuotaPriceCents   int `yaml:"quota_price_cents"`   // 充值包价格（分），默认 1000
	QuotaAdd          int `yaml:"quota_add"`           // 充值包次数，默认 1000
}

// AI 大模型配置（OpenAI 兼容接口，DeepSeek/通义/Moonshot 均可）。
type AI struct {
	Enabled bool   `yaml:"enabled"`
	BaseURL string `yaml:"base_url"` // 如 https://api.deepseek.com
	APIKey  string `yaml:"api_key"`
	Model   string `yaml:"model"`
}

// Admin 管理后台配置。用户名/密码为「首次启动种子」，写入库后以库中哈希为准。
type Admin struct {
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	SessionTTLHours int    `yaml:"session_ttl_hours"`
}

// Load 读取并解析配置文件，填充默认值。
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}
	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}
	cfg.applyDefaults()
	return cfg, nil
}

func (c *Config) applyDefaults() {
	if c.App.Name == "" {
		c.App.Name = "玄机命理"
	}
	if c.App.BaseURL == "" {
		c.App.BaseURL = "http://localhost:8080"
	}
	if c.App.PriceCents <= 0 {
		c.App.PriceCents = 990
	}
	if c.Server.Addr == "" {
		c.Server.Addr = ":8080"
	}
	if c.Database.Driver == "" {
		c.Database.Driver = "sqlite"
	}
	if c.Database.DSN == "" {
		c.Database.DSN = "data/fortune.db"
	}
	if c.Payment.Provider == "" {
		c.Payment.Provider = "mock"
	}
	if c.Unlock.FreePreviewRatio <= 0 {
		c.Unlock.FreePreviewRatio = 0.35
	}
	if c.Membership.MonthlyPriceCents <= 0 {
		c.Membership.MonthlyPriceCents = 1000
	}
	if c.Membership.MonthlyQuota <= 0 {
		c.Membership.MonthlyQuota = 1000
	}
	if c.Membership.QuotaPriceCents <= 0 {
		c.Membership.QuotaPriceCents = 1000
	}
	if c.Membership.QuotaAdd <= 0 {
		c.Membership.QuotaAdd = 1000
	}
	if c.AI.Model == "" {
		c.AI.Model = "deepseek-chat"
	}
	if c.Admin.Username == "" {
		c.Admin.Username = "admin"
	}
	if c.Admin.Password == "" {
		c.Admin.Password = "admin123"
	}
	if c.Admin.SessionTTLHours <= 0 {
		c.Admin.SessionTTLHours = 24
	}
}