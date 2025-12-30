package models

import "gorm.io/gorm"

// 用户表
type User struct {
	gorm.Model
	// gorm:"uniqueIndex" 是唯一索引，not null表示不为空
	Username string `json:"username" gorm:"type:varchar(30);uniqueIndex;not null"`
	// json:"-", - 表示不输出到json
	PassWd string `json:"-" gorm:"not null"`
}

func (m *User) GetId() uint {
	return m.ID
}
func (m *User) SetId(id uint) {
	m.ID = id
}
func (m *User) TableName() string {
	return "users"
}
