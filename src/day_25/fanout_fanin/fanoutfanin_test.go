package fanoutfanin

import (
	"fmt"
	"testing"
)

func TestFanoutFanin(t *testing.T) {
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
	fmt.Printf("1.并行检查所有倍数\n")
	c1 := CheckQingYiSe(hand)
	c2 := Check7Dui(hand)

	fmt.Printf("2.聚合结果\n")
	finalCh := mergeFans(c1, c2)
	totalScore := 0
	foundFans := 0
	for fan := range finalCh {
		fmt.Printf("发现种类:%s(倍数:%d)\n", fan.Name, fan.Score)
		totalScore += fan.Score
		foundFans++
	}
	fmt.Printf("总共发现%d种倍数\n", foundFans)
	fmt.Printf("最终的倍数是：%d\n", totalScore)
}
