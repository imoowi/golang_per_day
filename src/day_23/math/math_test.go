package math

import (
	"fmt"
	"testing"
)

func TestAdd(t *testing.T) {
	res := Add(3, 7)
	if res != 10 {
		t.Errorf("应该是10，结果是%d", res)
	}
}
func TestTableAdd(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{"正数", 3, 7, 10},
		{"零值", 0, 8, 8},
		{"负数", -6, 6, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Add(tt.a, tt.b)
			if got != tt.want {
				t.Fatalf("Add(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
func BenchmarkAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Add(3, 7)
	}
}

func ExampleAdd() {
	fmt.Println(Add(3, 7))
	// Output: 10
}
