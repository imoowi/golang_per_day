package main

import (
	"fmt"
	"reflect"
)

type User struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Age  int    `json:"age" db:"age"`
}

func (u *User) SetName(name string) {
	u.Name = name
	fmt.Println("Name updated to", u.Name)
}
func Parse(v interface{}) map[string]interface{} {
	t := reflect.TypeOf(v)
	val := reflect.ValueOf(v)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		val = val.Elem()
	}

	res := make(map[string]interface{})

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get("db")
		if tag == "" {
			continue
		}
		res[tag] = val.Field(i).Interface()
	}

	return res
}
