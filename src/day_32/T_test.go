package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestT(t *testing.T) {
	u := User{
		ID:   1,
		Name: "张三",
		Age:  10,
	}
	tp := reflect.TypeOf(u)
	for i := 0; i < tp.NumField(); i++ {
		f := tp.Field(i)
		fmt.Println(f.Name, f.Tag.Get("json"))
	}
	//修改name
	v := reflect.ValueOf(&u).Elem()
	v.FieldByName("Name").SetString("李四")
	fmt.Println(u)
	//调用Call
	m := reflect.ValueOf(&u).MethodByName("SetName")
	m.Call([]reflect.Value{reflect.ValueOf("王五")})
	fmt.Println(u)
	u2 := User{
		ID:   1,
		Name: "张三",
		Age:  10,
	}
	fmt.Println(Parse(u2))
}
