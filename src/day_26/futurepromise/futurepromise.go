package futurepromise

import (
	"fmt"
	"time"
)

// 异步计算的最终结果
type FanResult struct {
	// 总倍数
	TotalScore int
	// 详情
	Details string
	Err     error
}

// 结果占位符
type Future chan FanResult

// Get方法：阻塞等待并获取最终结果
func (f Future) Get() (FanResult, error) {
	result := <-f
	return result, result.Err
}

type Tile struct {
	// 牌名
	Name string
	// 牌值
	No int
}

// 胡牌结构体
type Hand struct {
	// 牌
	Tiles []Tile
	// 是否自摸
	IsSelfDraw bool
}

// CalcFanScore 是Promise函数，它接受胡牌信息，启动一个Goroutine异步任务，并立即返回Future占位符
func CalcFanScore(hand Hand) Future {
	resultCh := make(chan FanResult, 1)
	go func() {
		fmt.Println("[Promise] 计算倍数开始...")
		time.Sleep(2 * time.Second)
		totalScore := 11
		details := "清一色（6）倍+7对（5）倍=11倍"
		resultCh <- FanResult{
			TotalScore: totalScore,
			Details:    details,
			Err:        nil,
		}
		fmt.Println("[Promise] 计算倍数结束，并发送结果")
		close(resultCh)
	}()
	// 立即返回Future占位符
	return resultCh
}
