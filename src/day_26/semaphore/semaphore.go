package semaphore

import (
	"fmt"
	"sync"
	"time"
)

const MAX_ROOMS = 10

type RoomManager struct {
	roomSlots   chan struct{}
	activeRooms map[int]bool
	mu          sync.Mutex
	nextRoomID  int
}

func NewRoomManager() *RoomManager {
	return &RoomManager{
		roomSlots:   make(chan struct{}, MAX_ROOMS),
		activeRooms: make(map[int]bool),
		nextRoomID:  100,
	}
}

func (rm *RoomManager) OpenRoom() int {
	fmt.Println("[请求]房间创建请求到达，等待空闲槽位...")
	// 如果roomSlots已满，就在这儿等待
	rm.roomSlots <- struct{}{}
	rm.mu.Lock()
	roomid := rm.nextRoomID
	rm.nextRoomID++
	rm.activeRooms[roomid] = true
	currentLoad := len(rm.roomSlots)
	rm.mu.Unlock()
	fmt.Printf("[创建] 房间%d成功创建。当前房间数是：%d/%d\n", roomid, currentLoad, MAX_ROOMS)
	time.Sleep(time.Duration(roomid%3) * time.Second)
	return roomid
}

func (rm *RoomManager) CloseRoom(roomId int) {
	rm.mu.Lock()
	if !rm.activeRooms[roomId] {
		rm.mu.Unlock()
		return
	}
	delete(rm.activeRooms, roomId)
	rm.mu.Unlock()
	// 释放信号量
	<-rm.roomSlots
	currentLoad := len(rm.roomSlots)
	fmt.Printf("[释放] 房间%d成功释放。当前房间数是：%d/%d\n", roomId, currentLoad, MAX_ROOMS)

}
