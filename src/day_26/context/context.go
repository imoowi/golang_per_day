package contextdemo

import (
	"context"
	"fmt"
	"time"
)

func waitPlayerAction(ctx context.Context, player string) (string, error) {
	select {
	case <-time.After(5 * time.Second): //模拟玩家5秒都没有出牌
		return "", fmt.Errorf("玩家%s5秒内没有任何动作", player)
	case <-ctx.Done():
		return "", ctx.Err()
	}
}
