package fanoutfanin

import "sync"

type Tile struct {
	// 牌名
	Name string
	// 牌值
	No int
}

// 倍数结构体
type Fan struct {
	// 倍数名
	Name string
	// 倍数值
	Score int
}

// 胡牌结构体
type Hand struct {
	// 牌
	Tiles []Tile
	// 是否自摸
	IsSelfDraw bool
}

// --- Fan-out Worker 阶段---
// 清一色检查
func CheckQingYiSe(hand Hand) <-chan Fan {
	out := make(chan Fan)
	go func() {
		defer close(out)
		if len(hand.Tiles) > 0 {
			count := 0
			for _, tile := range hand.Tiles {
				if tile.Name == "万" {
					count++
				}
			}
			if count == len(hand.Tiles) {
				//假定清一色是6倍
				out <- Fan{Name: "清一色", Score: 6}
			}

		}
	}()
	return out
}

// 7对检查
func Check7Dui(hand Hand) <-chan Fan {
	out := make(chan Fan)
	go func() {
		defer close(out)
		if true {
			out <- Fan{Name: "7对", Score: 5}
		}
	}()
	return out
}

// --- Fan-in 聚合阶段---
func mergeFans(in ...<-chan Fan) <-chan Fan {
	var wg sync.WaitGroup
	out := make(chan Fan)
	output := func(c <-chan Fan) {
		defer wg.Done()
		for f := range c {
			out <- f
		}
	}
	wg.Add(len(in))
	for _, c := range in {
		go output(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
