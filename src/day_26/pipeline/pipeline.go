package pipeline

import "fmt"

type Play struct {
	// 玩家
	Player string
	// 牌
	Tile string
}

// 生成模拟出牌数据
func genPlays() <-chan Play {
	out := make(chan Play)
	go func() {
		out <- Play{Player: "张三", Tile: "1万"}
		out <- Play{Player: "李四", Tile: "2万"}
		out <- Play{Player: "王五", Tile: "3万"}
		close(out)
	}()
	return out
}

// 验证阶段
func stageValidate(in <-chan Play) <-chan Play {
	out := make(chan Play)
	go func() {
		for p := range in {
			// 简单验证逻辑,加入1万不合法
			if p.Tile == "1万" {
				fmt.Printf("%s的牌%s非法，已被丢弃\n", p.Player, p.Tile)
				continue
			}
			out <- p
		}
		close(out)
	}()
	return out
}

// 合法化阶段
func stageLegalize(in <-chan Play) <-chan Play {
	out := make(chan Play)
	go func() {
		for p := range in {
			// 假设条不合法
			if p.Tile == `条` {
				continue
			}
			// 简单合法化逻辑
			fmt.Printf("%s的牌%s已被合法化检查\n", p.Player, p.Tile)
			out <- p
		}
		close(out)
	}()
	return out
}

// 广播阶段
func stageBroadcast(in <-chan Play) {
	for p := range in {
		fmt.Printf("广播阶段:通知所有玩家 %s 出了 %s\n", p.Player, p.Tile)
	}
}
