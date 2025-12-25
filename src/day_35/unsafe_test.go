package main

import (
	"log"
	"reflect"
	"testing"
	"unsafe"
)

func TestUnsafeOffsetof(t *testing.T) {
	offsetA := unsafe.Offsetof((*CodeeJun)(nil).A)
	offsetB := unsafe.Offsetof((*CodeeJun)(nil).B)
	offsetC := unsafe.Offsetof((*CodeeJun)(nil).C)
	log.Println(offsetA, offsetB, offsetC)
	if offsetA == 0 && offsetB == 8 && offsetC == 16 {
		t.Log("offset of A, B and C are correct")
	} else {
		t.Error("offset of B and C are not correct")
	}
}

// Output:
// 0 8 16

func TestUnsafeOffsetofRW(t *testing.T) {
	cj := CodeeJun{A: 1, B: 2, C: 3}
	base := unsafe.Pointer(&cj)
	offsetB := unsafe.Offsetof((*CodeeJun)(nil).B)
	pb := (*int64)(unsafe.Pointer(uintptr(base) + offsetB))
	*pb = 99
	if *pb == 99 {
		t.Log("B is written correctly")
	} else {
		t.Error("B is not written correctly")
	}
}

func TestUnsafeZeroCopy(t *testing.T) {
	s := "Codee君"
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	res := *(*[]byte)(unsafe.Pointer(&bh))
	// res[0] = 67 //不能修改[]byte的值，否则会panic
	log.Println(res[0])
	log.Println(res)
	if string(res) == "Codee君" {
		t.Log("zero copy is correct")
	} else {
		t.Error("zero copy is not correct")
	}

	sb := []byte(s)
	log.Println(sb)
	// 真正零拷贝
	sbStr := *(*string)(unsafe.Pointer(&sb))
	log.Println(sbStr)
	if sbStr == "Codee君" {
		t.Log("zero copy is correct")
	} else {
		t.Error("zero copy is not correct")
	}
}

func TestUnsafeSlice(t *testing.T) {
	array := [5]int{1, 2, 3, 4, 5}
	addr := uintptr(unsafe.Pointer(&array[0]))
	newSlice := MakeAnySlice[int](addr, 3, 5)
	log.Println(newSlice)
	if newSlice[0] == 1 && newSlice[1] == 2 && newSlice[2] == 3 {
		t.Log("zero copy is correct")
	} else {
		t.Error("zero copy is not correct")
	}
	newSlice[0] = 9
	log.Println(newSlice)
	if array[0] == 9 {
		t.Log("zero copy is correct")
	} else {
		t.Error("zero copy is not correct")
	}
}

func TestUnsafePointerArithmetics(t *testing.T) {
	array := [5]int{1, 2, 3, 4, 5}
	addr := uintptr(unsafe.Pointer(&array[0]))
	secondAddr := Add(unsafe.Pointer(addr), unsafe.Sizeof(int(0)))
	log.Println(secondAddr)
	if *(*int)(secondAddr) == 2 {
		t.Log("pointer arithmetics is correct")
	} else {
		t.Error("pointer arithmetics is not correct")
	}
	thirdAddr := Add(secondAddr, unsafe.Sizeof(int(0)))
	log.Println(thirdAddr)
	if *(*int)(thirdAddr) == 3 {
		t.Log("pointer arithmetics is correct")
	} else {
		t.Error("pointer arithmetics is not correct")
	}
}
