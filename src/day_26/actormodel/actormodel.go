package actormodel

import "fmt"

type RoomMsg any

type Enter struct {
	Player string
}
type Exit struct {
	Player string
}
type Play struct {
	// 玩家
	Player string
	// 牌
	Tile string
}

func roomActor(roomId int, in <-chan RoomMsg) {
	fmt.Printf("房间%d开始游戏\n", roomId)
	viewers := 0
	for msg := range in {
		switch m := msg.(type) {
		case Enter:
			viewers++
			fmt.Printf("%s进入了房间\n", m.Player)
		case Exit:
			viewers--
			fmt.Printf("%s离开了房间\n", m.Player)
		case Play:
			fmt.Printf("%s打了一张牌：%s\n", m.Player, m.Tile)
		}
	}
	fmt.Println("房间", roomId, "结束游戏")
}
