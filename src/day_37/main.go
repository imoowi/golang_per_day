package main

import (
	"fmt"
)

var a, b int

func main() {
	fmt.Println("hello,Codee君")
	for i := 0; i < 5; i++ {

		go aa()
		go bb()
		fmt.Println(a, b)
	}
}

func aa() {
	a = 1
	b = 2
}

func bb() {
	print(b)
	print(a)
}
