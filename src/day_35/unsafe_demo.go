package main

import "unsafe"

type CodeeJun struct {
	A int32
	B int64
	C byte
}

type CodeeJunSlice struct {
	Data uintptr
	Len  int
	Cap  int
}

func MakeAnySlice[T any](addr uintptr, len, cap int) []T {
	hdr := CodeeJunSlice{Data: addr, Len: len, Cap: cap}
	return *(*[]T)(unsafe.Pointer(&hdr))
}
func Add(ptr unsafe.Pointer, offset uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(ptr) + offset)
}
