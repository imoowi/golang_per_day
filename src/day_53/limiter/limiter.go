package limiter

import (
	"sync"
	"time"
)

// TokenBucket 实现了令牌桶限流算法
type TokenBucket struct {
	capacity int        // 令牌桶容量
	tokens   int        // 当前令牌数量
	rate     int        // 每秒生成令牌速率
	lock     sync.Mutex // 互斥锁，保证线程安全
}

// NewTokenBucket 创建一个新的令牌桶限流器
// capacity: 令牌桶容量，表示最多可以存储多少个令牌
// rate: 每秒生成的令牌数量
func NewTokenBucket(capacity, rate int) *TokenBucket {
	tb := &TokenBucket{
		capacity: capacity, // 设置令牌桶容量
		tokens:   capacity, // 初始时令牌桶是满的
		rate:     rate,     // 设置每秒生成令牌的速率
	}
	// 启动令牌生成协程
	go tb.refill()
	return tb
}

// refill 定时生成令牌的方法
func (tb *TokenBucket) refill() {
	// 创建一个每秒触发一次的定时器
	ticker := time.NewTicker(time.Second)
	// 无限循环，每秒向令牌桶中添加令牌
	for range ticker.C {
		tb.lock.Lock()
		// 向令牌桶中添加rate个令牌
		tb.tokens += tb.rate
		// 确保令牌数量不超过桶容量
		if tb.tokens > tb.capacity {
			tb.tokens = tb.capacity
		}
		tb.lock.Unlock()
	}
}

// Allow 判断是否允许请求通过
// 返回true表示允许，false表示拒绝
func (tb *TokenBucket) Allow() bool {
	tb.lock.Lock()
	defer tb.lock.Unlock()

	// 如果有可用令牌，消耗一个令牌并返回true
	if tb.tokens > 0 {
		tb.tokens--
		return true
	}
	// 没有可用令牌，返回false
	return false
}
