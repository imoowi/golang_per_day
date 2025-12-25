package main

import (
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/exp/constraints"
)

type Maybe[T any] struct {
	v  T
	ok bool
}

func Some[T any](v T) Maybe[T] { return Maybe[T]{v, true} }
func None[T any]() Maybe[T]    { return Maybe[T]{ok: false} }

func (m Maybe[T]) Get() (T, bool) { return m.v, m.ok }

func (m Maybe[T]) Or(defaultV T) T {
	if m.ok {
		return m.v
	}
	return defaultV
}

// 查找map中是否存在key
func lookupMaybe[T comparable, U any](m map[T]U, k T) Maybe[U] {
	if v, ok := m[k]; ok {
		return Some(v)
	}
	return None[U]()
}

// 切片A转切片B
func SliceA2B[A any, B any](a []A, f func(A) B) []B {
	b := make([]B, len(a))
	for i, v := range a {
		b[i] = f(v)
	}
	return b
}

// 切片过滤
func Filter[T any](in []T, f func(T) bool) []T {
	out := make([]T, 0)
	for _, v := range in {
		if f(v) {
			out = append(out, v)
		}
	}
	return out
}

type Stack[T any] struct {
	data []*T
	mu   sync.Mutex
	len  atomic.Int64
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		data: make([]*T, 0),
	}
}
func (s *Stack[T]) Push(v *T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// 避免元素复制：直接存储指针*T
	s.data = append(s.data, v)
	s.len.Add(1)
}
func (s *Stack[T]) Pop() (*T, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.len.Load() == 0 {
		return nil, false
	}
	v := s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]
	s.len.Add(-1)
	return v, true
}

// 辅助API
// 返回栈当前元素数量
func (s *Stack[T]) Len() int64 {
	return s.len.Load()
}

type CodeeMap[K comparable, V any] struct {
	data map[K]V
	mu   sync.RWMutex
}

func NewCodeeMap[K comparable, V any]() *CodeeMap[K, V] {
	return &CodeeMap[K, V]{
		data: make(map[K]V),
	}
}

func (m *CodeeMap[K, V]) Set(k K, v V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[k] = v
}

func (m *CodeeMap[K, V]) Get(k K) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.data[k]
	return v, ok
}

// 定制可排序约束
func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// Ordered的意义：支持数字、string、rune等，只要是可用比大小的类型都可以

// 泛型二分查找
func BinarySearch[T constraints.Ordered](arr []T, target T) int {
	left, right := 0, len(arr)-1
	for left <= right {
		mid := left + (right-left)/2
		if arr[mid] == target {
			return mid
		} else if arr[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return -1
}

// 泛型归约算法
func Reduce[T any, R any](arr []T, init R, f func(R, T) R) R {
	acc := init
	for _, v := range arr {
		acc = f(acc, v)
	}
	return acc
}

// 泛型优先队列
type Item[T any] struct {
	value    T
	Priority int
}
type PriorityQueue[T any] []*Item[T]

func (pq PriorityQueue[T]) Len() int {
	return len(pq)
}
func (pq PriorityQueue[T]) Less(i, j int) bool {
	return pq[i].Priority < pq[j].Priority
}
func (pq PriorityQueue[T]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}
func (pq *PriorityQueue[T]) Push(x interface{}) {
	item := x.(*Item[T])
	*pq = append(*pq, item)
}
func (pq *PriorityQueue[T]) Pop() interface{} {
	old := *pq
	n := len(old)
	if n == 0 {
		var zero T
		return zero
	}
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}
func NewPriorityQueue[T any]() *PriorityQueue[T] {
	return &PriorityQueue[T]{}
}

type Task struct {
	Name    string
	DueDate time.Time
}
