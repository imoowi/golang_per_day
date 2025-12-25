package errgroupdemo

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func compute(ctx context.Context, name string, delay time.Duration) error {
	select {
	case <-time.After(delay):
		if rand.Intn(10) < 2 {
			fmt.Println(name, "失败")
			return fmt.Errorf("计算%s失败", name)
		}
		fmt.Println("计算", name, "成功")
		return nil
	case <-ctx.Done():
		fmt.Println("计算", name, "取消")
		return ctx.Err()
	}
}
