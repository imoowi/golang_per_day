package models

import (
	"log"
	"testing"
	"time"
	"unsafe"
)

func TestStruct(t *testing.T) {
	log.Println("BadStruct.length=", unsafe.Sizeof(BadStruct{}))
	log.Println("GoodStruct.length=", unsafe.Sizeof(GoodStruct{}))

	var c Counters
	go func() {
		for i := 0; i < 1000000; i++ {
			c.a++
			log.Println("c.a=", c.a, ";c.b=", c.b)
		}
	}()
	go func() {
		for i := 0; i < 1000000; i++ {
			c.b++
			log.Println("c.a=", c.a, ";c.b=", c.b)
		}
	}()

	time.Sleep(5 * time.Second)
	log.Println("test end")
}
