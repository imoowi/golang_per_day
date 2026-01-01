package components

import (
	"time"

	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// NatsBroker NATS消息队列代理
// 封装了NATS连接和JetStream上下文，提供消息发布订阅功能
type NatsBroker struct {
	NC *nats.Conn            // NATS连接对象
	JS nats.JetStreamContext // JetStream上下文，支持持久化消息
}

// NewNatsBroker 创建并初始化NATS消息队列代理
// 从配置文件读取NATS连接参数，建立连接并初始化JetStream
// 返回初始化完成的NatsBroker实例和可能的错误
func NewNatsBroker() (*NatsBroker, error) {
	url := viper.GetString("nats.url")
	// 1. 建立连接
	nc, err := nats.Connect(url,
		nats.MaxReconnects(10),            // 最大重连次数
		nats.ReconnectWait(2*time.Second), // 重连等待时间
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			zap.L().Panic("NATS Disconnected:", zap.String("url", url))
		}),
	)
	if err != nil {
		return nil, err
	}

	// 2. 启用 JetStream
	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	// 3. 预定义 Stream (如果不存在则创建)
	// 就像 Kafka 的 Topic 配置
	_name := viper.GetString("nats.stream.name")
	_subjects := viper.GetStringSlice("nats.stream.subjects")
	_max_age := viper.GetDuration("nats.stream.max_age")
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     _name,
		Subjects: _subjects,
		Storage:  nats.FileStorage, // 消息持久化到磁盘
		MaxAge:   _max_age,         // 消息保留24小时
	})
	if err != nil {
		zap.L().Panic("Stream info/create error:", zap.String("err", err.Error()))
	}

	return &NatsBroker{NC: nc, JS: js}, nil
}

// Close 关闭NATS连接
// 使用Drain方法优雅关闭，处理完本地缓存的消息后再断开连接
func (b *NatsBroker) Close() {
	if b.NC != nil {
		b.NC.Drain() // 优雅关闭：处理完本地缓存再断开
	}
}
