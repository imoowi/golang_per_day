package main

type Day46 struct {
	a uint64 // 8字节
	b uint64 // 8字节
}
type Day46WithPadding struct {
	a uint64
	_ [56]byte // 填充56字节,保证b在64字节边界上
	b uint64
	_ [64]byte // 填充64字节,足以撑开任何缓存行
}
