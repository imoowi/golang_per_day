package main

import (
	"log"
	"sync"
	"testing"
)

func TestModel(t *testing.T) {
	d := Day46{
		a: 1,
		b: 2,
	}
	// 并发写入 a 和 b，制造伪共享
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			d.a++
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			d.b++
		}
	}()

	wg.Wait()
	log.Printf("a=%d, b=%d", d.a, d.b)
	log.Println("TestModel passed")
}

// go test -bench=. ./day_46/
func BenchmarkModel(b *testing.B) {
	d := Day46{}
	b.Run("NoPadding", func(b *testing.B) {
		var wg sync.WaitGroup
		wg.Add(2)
		for i := 0; i < 2; i++ {
			go func(idx int) {
				defer wg.Done()
				for j := 0; j < b.N; j++ {
					if idx == 0 {
						d.a++
					} else {
						d.b++
					}
				}
			}(i)
		}
		wg.Wait()
	})
	dp := Day46WithPadding{}
	b.Run("WithPadding", func(b *testing.B) {
		var wg sync.WaitGroup
		wg.Add(2)
		for i := 0; i < 2; i++ {
			go func(idx int) {
				defer wg.Done()
				for j := 0; j < b.N; j++ {
					if idx == 0 {
						dp.a++
					} else {
						dp.b++
					}
				}
			}(i)
		}
		wg.Wait()
	})
}
