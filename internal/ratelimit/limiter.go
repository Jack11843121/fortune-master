// Package ratelimit 提供进程内固定窗口限流器，用于登录、AI 对话等敏感接口的滥用防护。
// 单实例部署下即可用；若多实例，需替换为 Redis 等共享实现。
package ratelimit

import (
	"sync"
	"time"
)

type bucket struct {
	count int
	reset time.Time
}

// Limiter 按 key 限制窗口内请求次数。
type Limiter struct {
	mu     sync.Mutex
	items  map[string]*bucket
	limit  int
	window time.Duration
}

// New 创建限流器：window 内最多 limit 次。
func New(limit int, window time.Duration) *Limiter {
	return &Limiter{items: map[string]*bucket{}, limit: limit, window: window}
}

// Allow 检查 key 是否在限流内；允许则计数并返回 true。
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	l.prune(now)

	b, ok := l.items[key]
	if !ok || !now.Before(b.reset) {
		l.items[key] = &bucket{count: 1, reset: now.Add(l.window)}
		return true
	}
	if b.count >= l.limit {
		return false
	}
	b.count++
	return true
}

// prune 惰性清理过期桶，避免 map 无限增长。
func (l *Limiter) prune(now time.Time) {
	if len(l.items) < 10000 {
		return
	}
	for k, b := range l.items {
		if !now.Before(b.reset) {
			delete(l.items, k)
		}
	}
}