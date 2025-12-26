package models

type UserA struct {
	ID       uint
	Online   bool
	NickName string
	Blocked  bool
	Avatar   *string
	Deleted  bool
	Albumns  []string
}

type UserB struct {
	ID       uint
	NickName string
	Online   bool
	Avatar   *string
	Blocked  bool
	Deleted  bool
	Albumns  []string
}

type UserC struct {
	Albumns  []string
	NickName string
	Avatar   *string
	ID       uint
	Online   bool
	Blocked  bool
	Deleted  bool
}
type UserGood struct {
	Albumns  []string //24 bytes（对齐值8）
	NickName string   //16 bytes（对齐值8）
	Avatar   *string  //8 bytes（对齐值8）
	ID       uint     //8 bytes（对齐值8）
	Online   bool     //1 bytes（对齐值1）
	Blocked  bool     //1 bytes（对齐值1）
	Deleted  bool     //1 bytes（对齐值1）
	// 5 bytes 填充（对齐值8）
	// 总计：64 bytes
}

type BadStruct struct {
	a bool
	b int64
	c bool
}

type GoodStruct struct {
	b int64
	a bool
	c bool
}
type Counters struct {
	a int64
	b int64
}
