package main

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

func main() {
	x := Min(1, 3)              // T = int
	y := Min(1.0, 2.0)          // T = float64
	z := Min("hello", "Codee君") // T = string
	// xx := Min(1, "2")
	fmt.Println(x, y, z)

	var dog Dog
	dog.name = "丹尼"
	var cat Cat
	cat.name = "坎迪"
	var pig Pig
	pig.name = "佩奇"
	var animals Animal[any]
	animals.Append(dog)
	animals.Append(cat)
	animals.Append(pig)
	fmt.Println(animals)

	arr := []int{1, 2, 3, 4, 5}
	fmt.Println(IndexOf(arr, 3))
	fmt.Println(IndexOf(arr, 6))
	c := Add[MyInt](1, 2)
	fmt.Println(c)
	fmt.Println(Sum(arr))
	fmt.Println(Sum([]float64{1.1, 2.2, 3.3}))
	fmt.Println(Sum([]MyInt{1, 2, 3}))
}
func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

type Animal[T any] struct {
	data []T
}

func (a *Animal[T]) Append(b T) {
	a.data = append(a.data, b)
}

type Dog struct {
	name string
}
type Cat struct {
	name string
}
type Pig struct {
	name string
}

// func Get(x any) // x 是 interface{}类型
// func Get[T any](x T) // x 是具体类型
func Print[T any](x T) {
	fmt.Println(x)
}

func IndexOf[T comparable](arr []T, x T) int {
	for i, v := range arr {
		if v == x {
			return i
		}
	}
	return -1
}

type MyInt int

// ~int 表示类型约束：T 可以是 int 或任何以 int 为底层类型的自定义类型（如 MyInt）
func Add[T ~int](a, b T) T {
	return a + b
}

type Number interface {
	~int | ~float64
}

func Sum[T Number](arr []T) T {
	var sum T
	for _, v := range arr {
		sum += v
	}
	return sum
}

// func (a *Animal[T]) AppendMore[U any](x U) {} // 错误！
