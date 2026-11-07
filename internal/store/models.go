package store

import "time"

// Reading 保存一次算命生成的完整结果（结构化 JSON 存入 Payload）。
type Reading struct {
	ID        string    `gorm:"primaryKey;size:64" json:"id"`
	Scenario  string    `gorm:"size:32;index" json:"scenario"` // personal / marriage / hiring
	Title     string    `gorm:"size:255" json:"title"`
	Payload   string    `gorm:"type:longtext" json:"-"` // 完整结果 JSON（含 preview 与 full）
	CreatedAt time.Time `json:"created_at"`
}

// Order 表示一笔订单。Type 区分业务：unlock(解锁报告)/subscribe(会员订阅)/quota(AI 次数充值)。
type Order struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	OrderNo     string     `gorm:"uniqueIndex;size:64" json:"order_no"`
	Type        string     `gorm:"size:16;index;default:unlock" json:"type"`
	UserID      string     `gorm:"size:64;index" json:"user_id"`
	ReadingID   string     `gorm:"index;size:64" json:"reading_id"`
	Scenario    string     `gorm:"size:32" json:"scenario"`
	Amount      int        `json:"amount"`                       // 单位：分
	Status      string     `gorm:"size:16;index" json:"status"`  // pending / paid / closed
	Provider    string     `gorm:"size:16" json:"provider"`      // mock / wechat / alipay
	WechatTxnID string     `gorm:"size:64" json:"wechat_txn_id"` // 第三方交易号（微信/支付宝通用）
	PayParams   string     `gorm:"type:text" json:"-"`           // 支付跳转参数（二维码/表单等）
	PaidAt      *time.Time `json:"paid_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// User 会员用户。身份：本地设备 uid（前端 localStorage 生成），可后续绑定微信 openid。
type User struct {
	ID        string     `gorm:"primaryKey;size:64" json:"id"` // uid
	OpenID    string     `gorm:"size:64;index" json:"openid"`  // 微信/支付宝 openid
	UnionID   string     `gorm:"size:64;index" json:"unionid"` // 微信 unionid（跨应用唯一，可选）
	Provider  string     `gorm:"size:16;default:dev;index" json:"provider"` // wechat / alipay / dev
	Nickname  string     `gorm:"size:64" json:"nickname"`
	Avatar    string     `gorm:"size:500" json:"avatar"`
	VipUntil  *time.Time `json:"vip_until"` // 会员到期时间，空=非会员
	AICredits int        `json:"ai_credits"` // 剩余 AI 调用次数（余额）
	Growth    int        `json:"growth"`     // 成长值
	Level     int        `gorm:"default:1" json:"level"` // 成长等级
	Status    string     `gorm:"size:16;default:active" json:"status"` // active / banned
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// IsVIP 判断用户当前是否为有效会员。
func (u User) IsVIP() bool {
	return u.VipUntil != nil && u.VipUntil.After(time.Now())
}

// UserSession 用户登录会话（令牌哈希落库，令牌明文只回客户端一次）。
type UserSession struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    string    `gorm:"size:64;index" json:"user_id"`
	TokenHash string    `gorm:"uniqueIndex;size:64" json:"-"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// MembershipLog 会员/额度流水（订阅、续费、充值、后台调整），用于会员管理与成长审计。
type MembershipLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    string    `gorm:"size:64;index" json:"user_id"`
	Kind      string    `gorm:"size:16" json:"kind"` // subscribe / renew / topup / admin / consume
	Amount    int       `json:"amount"`              // 金额（分）
	Credits   int       `json:"credits"`             // 增/减的 AI 次数（可为负）
	Days      int       `json:"days"`                // 增加/调整的天数
	Note      string    `gorm:"size:255" json:"note"`
	CreatedAt time.Time `json:"created_at"`
}

// AiMessage AI 算命师与用户的对话消息（同时作为「日常咨询」记录供后台管理）。
type AiMessage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    string    `gorm:"size:64;index" json:"user_id"`
	SessionID string    `gorm:"size:64;index" json:"session_id"`
	Role      string    `gorm:"size:16" json:"role"` // user / assistant / admin
	Content   string    `gorm:"type:text" json:"content"`
	Credits   int       `json:"credits"` // 本条消耗的 AI 次数（assistant 应答）
	CreatedAt time.Time `json:"created_at"`
}

// CommunityPost 社区帖子。
type CommunityPost struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    string    `gorm:"size:64;index" json:"user_id"`
	Title     string    `gorm:"size:128" json:"title"`
	Content   string    `gorm:"type:text" json:"content"`
	Category  string    `gorm:"size:32;index" json:"category"` // 命理趣谈 / 合婚故事 / 求职问卦 / 宝宝取名 …
	Likes     int       `json:"likes"`
	Status    string    `gorm:"size:16;index;default:published" json:"status"` // published / pending / hidden
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CommunityComment 社区评论。
type CommunityComment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	PostID    uint      `gorm:"index" json:"post_id"`
	UserID    string    `gorm:"size:64;index" json:"user_id"`
	Content   string    `gorm:"type:text" json:"content"`
	Status    string    `gorm:"size:16;default:published" json:"status"` // published / hidden
	CreatedAt time.Time `json:"created_at"`
}

// ShareRecord 内容分享记录（转发到微信等）。
type ShareRecord struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    string    `gorm:"size:64;index" json:"user_id"`
	ReadingID string    `gorm:"size:64;index" json:"reading_id"`
	Platform  string    `gorm:"size:32" json:"platform"`
	CreatedAt time.Time `json:"created_at"`
}

// AdminUser 管理后台账号（密码存 bcrypt 哈希）。
type AdminUser struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"uniqueIndex;size:64" json:"username"`
	PasswordHash string    `gorm:"size:255" json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

// AdminSession 管理后台会话（存令牌哈希）。
type AdminSession struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	AdminID   uint      `gorm:"index" json:"admin_id"`
	TokenHash string    `gorm:"uniqueIndex;size:64" json:"-"` // SHA-256(token)
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}