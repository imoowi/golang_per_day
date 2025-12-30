package migrates

import (
	"codee_jun/internal/components"

	"gorm.io/gorm"
)

// 定义迁移函数
type migrate func(*gorm.DB)

// 这里放所有的迁移
var routers = []migrate{}

// 调用迁移
func DoMigrate() {
	for _, route := range routers {
		route(components.DB)
	}
}

// 注册迁移函数
func regMigrate(r ...migrate) {
	routers = append(routers, r...)
}
