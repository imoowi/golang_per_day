package models

import "gorm.io/gorm"

type Goods struct {
	gorm.Model `swaggerignore:"true" `
	Name       string  `json:"name" gorm:"column:name;type:varchar(100);not null;comment:'商品名称'"`
	Price      float64 `json:"price" gorm:"column:price;type:decimal(10,2);not null;comment:'商品价格'"`
	Stock      int     `json:"stock" gorm:"column:stock;type:int;not null;comment:'库存数量'"`
}

func (m *Goods) GetId() uint {
	return m.ID
}
func (m *Goods) SetId(id uint) {
	m.ID = id
}
func (m *Goods) TableName() string {
	return "goods"
}
