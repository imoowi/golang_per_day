package main

import "C"
import (
	"fmt"
)

func main() {
	fmt.Println("Hello, Codee君!")
}

//export GoPrint
func GoPrint(i C.int) {
	fmt.Println("GoPrint:", int(i))
}
