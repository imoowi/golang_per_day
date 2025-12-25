package main

import (
	"fmt"
	"unsafe"
)

func main() {
	fmt.Println("Hello,Codee君")
	// *T <--> unsafe.Pointer
	var x int = 10
	p := unsafe.Pointer(&x)
	fmt.Println("x:", x)
	fmt.Println("p:", p)
	px := (*int)(p)
	fmt.Println("*px:", *px)

	// unsafe.Pointer <--> uintptr
	addr := uintptr(p) + unsafe.Sizeof(x)
	fmt.Println("addr:", addr)
	p2 := unsafe.Pointer(addr)
	fmt.Println("p2:", p2)
}
