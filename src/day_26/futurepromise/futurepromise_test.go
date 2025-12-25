package futurepromise

import (
	"fmt"
	"testing"
	"time"
)

func TestFuturePromise(t *testing.T) {
	hand := Hand{
		Tiles: []Tile{
			{Name: "万", No: 1},
			{Name: "万", No: 2},
			{Name: "万", No: 3},
			{Name: "万", No: 4},
			{Name: "万", No: 5},
			{Name: "万", No: 6},
			{Name: "万", No: 7},
			{Name: "万", No: 7},
			{Name: "万", No: 8},
			{Name: "万", No: 8},
			{Name: "万", No: 8},
			{Name: "万", No: 9},
			{Name: "万", No: 9},
			{Name: "万", No: 9},
		},
		IsSelfDraw: false,
	}
	fmt.Println("---end---")
	fmt.Println("玩家胡牌了！")
	fanFuture := CalcFanScore(hand)
	fmt.Println("启动异步计算倍数")
	time.Sleep(time.Millisecond * 500)
	fmt.Println("播放玩家胡牌音乐，通知其他玩家XXX胡牌了")
	time.Sleep(time.Millisecond * 500)
	//最后，调用fanFuture.Get来获取最终结果，知道结果计算出来为止
	finalRes, err := fanFuture.Get()
	if err != nil {
		fmt.Printf("倍数计算出错:%v\n", err)
		return
	}
	fmt.Println("倍数计算结果出来了")
	fmt.Println("总倍数是", finalRes.TotalScore)
	fmt.Println("详情：", finalRes.Details)
	fmt.Println("---end---")

}
