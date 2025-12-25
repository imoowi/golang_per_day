package semaphore

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestSemaphore(t *testing.T) {
	manager := NewRoomManager()
	var wg sync.WaitGroup
	const totalRequests = 12
	fmt.Printf("服务器最大容量：%d个房间。\n", MAX_ROOMS)
	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		go func(requestNum int) {
			defer wg.Done()
			roomId := manager.OpenRoom()
			//模拟房间正在打牌耗时1秒
			time.Sleep(time.Second)
			manager.CloseRoom(roomId)
		}(i + 1)
	}
	wg.Wait()
	fmt.Println("所有房间任务已经处理完毕。")
}
