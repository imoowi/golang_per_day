package repos

import (
	"codee_jun/internal/components"
	"codee_jun/internal/interfaces"
	"codee_jun/internal/models"
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

// Goods 商品仓储实例，提供商品数据的CRUD操作
var (
	Goods            *GoodsRepo
	GoodsBloomFilter *components.RedisBloomFilter
)

// GoodsRepo 商品仓储结构体，实现了Repo接口
// 封装了商品数据的数据库操作、缓存操作和布隆过滤器
type GoodsRepo struct {
	interfaces.Repo[*models.Goods]
}

// NewGoodsRepo 创建商品仓储实例
// 初始化Goods实例，并初始化布隆过滤器
func NewGoodsRepo() {
	Goods = &GoodsRepo{
		Repo: *interfaces.NewRepo[*models.Goods](components.DB),
	}
	initGoodsBloomFilter()
}

// init 包初始化函数，将商品仓储注册到仓储注册表中
func init() {
	RegisterRepos(NewGoodsRepo)
}

// initGoodsBloomFilter 初始化商品布隆过滤器
// 使用Redis实现布隆过滤器，用于快速判断商品ID是否存在
func initGoodsBloomFilter() {
	GoodsBloomFilter = components.NewRedisBloomFilter(components.Redis, "goods_bloom_filter")
	Goods.InitAllIds()
}

// InitAllIds 初始化所有商品ID到布隆过滤器
// 从数据库中查询所有商品ID，并添加到布隆过滤器中
// 用于防止缓存穿透，快速判断商品ID是否存在
func (r *GoodsRepo) InitAllIds() {
	var ids []uint
	r.Repo.DB.Model(&models.Goods{}).Pluck("id", &ids)
	for _, id := range ids {
		GoodsBloomFilter.Add(context.Background(), id)
	}
}

// AddOne 添加一个新商品
// 在添加商品后，将新商品的ID添加到布隆过滤器中
// 参数：
//   - c: Gin上下文
//   - model: 商品模型
//
// 返回：
//   - newId: 新插入的商品ID
//   - err: 错误信息
func (r *GoodsRepo) AddOne(c *gin.Context, model *models.Goods) (newId uint, err error) {
	newId, err = r.Add(c, model)
	GoodsBloomFilter.Add(context.Background(), newId)
	return
}

// GetById 根据ID获取商品信息
// 实现了多级缓存策略：
// 1. 布隆过滤器快速拦截不存在的ID（防止缓存穿透）
// 2. Redis缓存读取
// 3. SingleFlight击穿保护（防止大量并发请求打到数据库）
// 4. 数据库查询
// 5. 缓存回填（带随机过期时间防止缓存雪崩）
// 参数：
//   - c: Gin上下文
//   - id: 商品ID
//
// 返回：
//   - res: 商品模型
//   - err: 错误信息
func (r *GoodsRepo) GetById(c *gin.Context, id uint) (res *models.Goods, err error) {
	// --- 0. 第一步：Redis 布隆过滤器拦截 ---
	// 只有返回 false 时才确定不存在，返回 true 可能是误报，需进一步查缓存/DB
	// 布隆过滤器的作用：快速判断商品ID是否可能存在，避免无效的缓存和数据库查询
	exists, err := GoodsBloomFilter.Exists(c.Request.Context(), id)
	if err == nil && !exists {
		return nil, gorm.ErrRecordNotFound
	}

	// 1. 尝试读缓存
	// 构建缓存key，格式为 "goods:{id}"
	key := "goods:" + cast.ToString(id)
	val, err := components.Redis.Get(c.Request.Context(), key).Result()
	if err == nil {
		// 命中缓存穿透占位符，说明商品不存在
		if val == cacheNilValue {
			return nil, gorm.ErrRecordNotFound
		}
		// 反序列化缓存数据
		err = json.Unmarshal([]byte(val), &res)
		if err != nil {
			return
		}
		return
	}
	// 2. 击穿保护：使用singleflight合并请求
	// 当缓存失效时，大量并发请求会同时打到数据库，造成击穿
	// 使用singleflight确保同一时间只有一个请求去查数据库，其他请求等待结果
	rawRes, err, _ := sfGroup.Do(key, func() (interface{}, error) {
		// 双重检查：防止在等待锁的时候，缓存已被其他Goroutine填好
		val, err := components.Redis.Get(c.Request.Context(), key).Result()
		if err == nil {
			if val == cacheNilValue {
				return nil, gorm.ErrRecordNotFound
			}
			err = json.Unmarshal([]byte(val), &res)
			return res, err
		}

		// 3. 查DB并处理穿透
		var goods *models.Goods
		err = r.Repo.DB.WithContext(c.Request.Context()).Where(`id=?`, id).First(&goods).Error
		if err != nil {
			// 如果商品不存在，使用缓存穿透防护策略
			// 缓存空对象（占位符），设置较短过期时间（5分钟）
			// 这样后续请求会直接从缓存中获取到占位符，避免频繁查询数据库
			if errors.Is(err, gorm.ErrRecordNotFound) {
				components.Redis.Set(c.Request.Context(), key, cacheNilValue, 5*time.Minute)
				return nil, err
			}
			return nil, err
		}

		// 4. 缓存数据回填+雪崩保护（随机过期时间）
		// 将查询到的商品数据序列化后存入Redis
		data, _ := json.Marshal(goods)
		// 随机过期时间，防止缓存雪崩
		// 所有缓存同时过期会导致大量请求同时打到数据库
		// 设置300-600秒之间的随机过期时间，分散缓存失效时间
		ttl := time.Duration(rand.Intn(300)+300) * time.Second
		components.Redis.Set(c.Request.Context(), key, data, ttl)
		return goods, nil
	})
	if err != nil {
		return
	}
	res = rawRes.(*models.Goods)
	return
}

// UpdateById 根据ID更新商品信息
// 实现了数据库事务更新和缓存删除策略：
// 1. 开启事务更新数据库
// 2. 更新成功后删除缓存
// 3. 使用延迟双删策略，防止并发场景下的脏数据
// 参数：
//   - c: Gin上下文
//   - id: 商品ID
//   - data: 需要更新的字段和值
//
// 返回：
//   - err: 错误信息
func (r *GoodsRepo) UpdateById(c *gin.Context, id uint, data map[string]interface{}) (err error) {
	key := "goods:" + cast.ToString(id)
	// 1. 开启事务更新DB
	err = r.Repo.DB.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		// 2. 更新DB
		// Omit(`created_at`) 排除创建时间字段，不更新
		var goods *models.Goods
		err = tx.Model(&goods).Omit(`created_at`).Where(`id=?`, id).Updates(data).Error
		if err != nil {
			return err
		}
		// 3. 更新完DB后直接删除缓存
		// 采用Cache Aside模式：先更新数据库，再删除缓存
		// 删除缓存而不是更新缓存，因为删除操作更简单，且避免并发更新的问题
		components.Redis.Del(c.Request.Context(), key)
		return nil
	})

	// 延迟双删
	// 在删除缓存后，延迟500ms再次删除缓存
	// 目的：防止在更新数据库和删除缓存之间，有其他请求读取了旧数据并写入了缓存
	// 延迟时间需要大于主从复制的延迟时间（如果有主从复制）
	time.AfterFunc(500*time.Millisecond, func() {
		components.Redis.Del(c.Request.Context(), key)
	})
	return
}
