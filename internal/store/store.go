package store

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"fortune/internal/auth"
)

// Open 根据 driver 打开数据库连接。sqlite 使用纯 Go 驱动（免 CGO，Windows/Linux 均无需额外工具链）。
func Open(driver, dsn string) (*gorm.DB, error) {
	var dialector gorm.Dialector
	switch driver {
	case "mysql":
		dialector = mysql.Open(dsn)
	case "sqlite":
		if dir := filepath.Dir(dsn); dir != "." && dir != "" {
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return nil, fmt.Errorf("创建 sqlite 目录失败: %w", err)
			}
		}
		dialector = sqlite.Open(dsn)
	default:
		return nil, fmt.Errorf("不支持的数据库驱动: %q（支持 sqlite / mysql）", driver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}
	return db, nil
}

// Migrate 自动建表 / 迁移。
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Reading{}, &Order{},
		&User{}, &MembershipLog{}, &AiMessage{},
		&CommunityPost{}, &CommunityComment{}, &ShareRecord{},
		&AdminUser{}, &AdminSession{}, &UserSession{},
	)
}

// SeedAdmin 首次启动时创建默认管理员（密码用 bcrypt 哈希落库）。
// 已存在同名账号时不做任何事，以库中哈希为准（配置中的密码仅作初始种子）。
func SeedAdmin(db *gorm.DB, username, password string) error {
	if username == "" || password == "" {
		return nil
	}
	var existing AdminUser
	if err := db.First(&existing, "username = ?", username).Error; err == nil {
		return nil
	}
	hash, err := auth.HashPassword(password)
	if err != nil {
		return fmt.Errorf("管理员密码哈希失败: %w", err)
	}
	return db.Create(&AdminUser{Username: username, PasswordHash: hash}).Error
}