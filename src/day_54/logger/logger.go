package logger

import "go.uber.org/zap"

// Log 是全局的zap日志实例
var Log *zap.Logger

// Init 初始化全局日志实例
// 配置使用生产级别日志格式，并将输出重定向到标准输出
func Init() {
	// 创建生产级别日志配置
	cfg := zap.NewProductionConfig()
	// 设置日志输出路径为标准输出
	cfg.OutputPaths = []string{"stdout"}

	// 构建日志实例
	var err error
	Log, err = cfg.Build()
	if err != nil {
		// 如果初始化失败，直接panic
		panic(err)
	}
}
