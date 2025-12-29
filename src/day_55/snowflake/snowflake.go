package snowflake

import (
	"errors"
	"strconv"
	"sync"
	"time"
)

// 常量定义
const (
	workerBits   = 10                      // 工作节点位数
	seqBits      = 12                      // 序列号位数
	workerMax    = -1 ^ (-1 << workerBits) // 工作节点最大ID
	seqMask      = -1 ^ (-1 << seqBits)    // 序列号掩码
	timeShift    = workerBits + seqBits    // 时间戳左移位数
	workerShift  = seqBits                 // 工作节点左移位数
	defaultEpoch = int64(1672531200000)    // 默认起始时间戳 (2023-01-01)
)

// ID 自定义类型，用于区分雪花ID和普通int64
type ID int64

// Snowflake 雪花算法生成器
type Snowflake struct {
	mu       sync.Mutex
	lastTime int64
	workerID int64
	sequence int64
	epoch    int64
}

// New 创建雪花算法生成器
// workerID: 工作节点ID，范围 0~1023
// 返回错误如果workerID超出范围
func New(workerID int64) (*Snowflake, error) {
	return NewWithEpoch(workerID, defaultEpoch)
}

// NewWithEpoch 创建带自定义起始时间的雪花算法生成器
// workerID: 工作节点ID，范围 0~1023
// epoch: 自定义起始时间戳（毫秒）
// 返回错误如果workerID超出范围
func NewWithEpoch(workerID int64, epoch int64) (*Snowflake, error) {
	if workerID < 0 || workerID > workerMax {
		return nil, errors.New("worker ID out of range [0, 1023]")
	}

	return &Snowflake{
		workerID: workerID,
		epoch:    epoch,
	}, nil
}

// NextID 生成下一个雪花ID
// 返回ID类型的雪花ID和可能的错误
func (s *Snowflake) NextID() (ID, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UnixMilli()

	// 处理时间回拨
	if now < s.lastTime {
		return 0, errors.New("time is back, ID generation failed")
	}

	if now == s.lastTime {
		// 同一毫秒内，递增序列号
		s.sequence = (s.sequence + 1) & seqMask
		// 序列号耗尽，等待下一个毫秒
		if s.sequence == 0 {
			// 使用短暂休眠代替自旋等待，减少CPU占用
			time.Sleep(time.Millisecond)
			now = time.Now().UnixMilli()

			// 处理时间回拨（再次检查）
			if now < s.lastTime {
				return 0, errors.New("time is back, ID generation failed")
			}

			s.lastTime = now
			s.sequence = 0
		}
	} else {
		// 新的毫秒，重置序列号
		s.lastTime = now
		s.sequence = 0
	}

	// 生成ID
	id := ((now - s.epoch) << timeShift) |
		(s.workerID << workerShift) |
		s.sequence

	return ID(id), nil
}

// ParseID 解析雪花ID
// 返回ID的各组成部分：时间戳、工作节点ID、序列号
func ParseID(id ID, epoch int64) (time.Time, int64, int64) {
	idInt := int64(id)
	timestamp := (idInt >> timeShift) + epoch
	workerID := (idInt >> workerShift) & ((1 << workerBits) - 1)
	sequence := idInt & seqMask

	return time.UnixMilli(timestamp), workerID, sequence
}

// String 将ID转换为字符串
func (id ID) String() string {
	return strconv.FormatInt(int64(id), 10)
}

// Int64 将ID转换为int64
func (id ID) Int64() int64 {
	return int64(id)
}
