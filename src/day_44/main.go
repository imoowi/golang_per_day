package main

import (
	"fmt"
	"sync"
)

func main() {
	fmt.Println("hello Codee君")
	var bufPool = sync.Pool{
		New: func() any {
			return make([]byte, 32*1024)
		},
	}
	b := bufPool.Get().([]byte)
	fmt.Println(len(b))
	bufPool.Put(b)
	b = bufPool.Get().([]byte)
	fmt.Println(len(b))
}
