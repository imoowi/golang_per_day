package components

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisBloomFilter struct {
	client *redis.Client
	key    string
}

func NewRedisBloomFilter(client *redis.Client, key string) *RedisBloomFilter {
	return &RedisBloomFilter{client: client, key: key}
}

// 初始化布隆过滤器：设置错误率 0.01 和 预期容量 1000000
func (b *RedisBloomFilter) Reserve(ctx context.Context, errorRate float64, capacity int64) error {
	return b.client.Do(ctx, "BF.RESERVE", b.key, errorRate, capacity).Err()
}

// 添加元素
func (b *RedisBloomFilter) Add(ctx context.Context, id interface{}) error {
	return b.client.Do(ctx, "BF.ADD", b.key, id).Err()
}

// 检查元素是否存在
func (b *RedisBloomFilter) Exists(ctx context.Context, id interface{}) (bool, error) {
	exists, err := b.client.Do(ctx, "BF.EXISTS", b.key, id).Bool()
	if err != nil {
		return false, err
	}
	return exists, nil
}
