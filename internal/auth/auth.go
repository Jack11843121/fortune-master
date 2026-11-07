// Package auth 提供密码哈希（bcrypt）与会话令牌（随机 + SHA-256 哈希）的生成与校验。
// 明文密码永不落盘；落库的只有 bcrypt 哈希与令牌哈希。
package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword 用 bcrypt 生成密码哈希。
func HashPassword(pw string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// CheckPassword 校验明文密码与 bcrypt 哈希是否匹配。
func CheckPassword(hash, pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw)) == nil
}

// NewToken 生成高熵随机令牌，返回明文与它的 SHA-256 哈希（哈希落库，明文只回给客户端一次）。
func NewToken() (token, hash string, err error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", err
	}
	token = hex.EncodeToString(b)
	return token, HashToken(token), nil
}

// HashToken 计算令牌的 SHA-256 十六进制哈希。
func HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// CheckToken 以常数时间比较令牌与哈希是否相等，避免时序侧信道。
func CheckToken(hash, token string) bool {
	h := HashToken(token)
	return subtle.ConstantTimeCompare([]byte(hash), []byte(h)) == 1
}