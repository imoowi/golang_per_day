package interfaces

// 模型接口
type IModel interface {
	GetId() uint
	SetId(uint)
	TableName() string
}
