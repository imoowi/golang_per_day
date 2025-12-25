package contextdemo

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	actionCh := make(chan string, 1)
	go func() {
		if act, err := waitPlayerAction(context.Background(), "张三"); err == nil {
			actionCh <- act
		} else {
			//没有动作
		}
	}()

	select {
	case act := <-actionCh:
		fmt.Println("玩家动作", act)
	case <-ctx.Done():
		fmt.Println("超时，自动打一张牌")
	}
}
