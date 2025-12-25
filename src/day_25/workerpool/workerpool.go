package workerpool

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	// 桌面id
	TableID int
	// 操作类型
	Op string
}

// worker 池处理任务
func worker(id int, tasks <-chan Task, wg *sync.WaitGroup) {
	defer wg.Done()
	for t := range tasks {
		fmt.Printf("工人 %d 正在给 %d 号桌 %s\n", id, t.TableID, t.Op)
		// 模拟处理任务
		time.Sleep(time.Millisecond * 500)
	}
}
