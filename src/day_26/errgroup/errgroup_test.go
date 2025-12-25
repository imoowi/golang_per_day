package errgroupdemo

import (
	"context"
	"fmt"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"
)

func TestErrGroup(t *testing.T) {
	ctx := context.Background()
	g, ctx := errgroup.WithContext(ctx)
	methods := []struct {
		name  string
		delay time.Duration
	}{
		{"平胡", 100 * time.Millisecond},
		{"七对", 200 * time.Millisecond},
		{"清一色", 150 * time.Millisecond},
	}
	for _, m := range methods {
		m := m //核心：创建循环变量 m 的局部副本
		g.Go(func() error {
			return compute(ctx, m.name, m.delay)
		})
	}
	if err := g.Wait(); err != nil {
		fmt.Println("有计算失败了：", err)
		return
	}
	fmt.Println("所有计算都成功了")
}
