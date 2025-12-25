package workerpool

import (
	"fmt"
	"sync"
	"testing"
)

func TestWorker(t *testing.T) {
	// 创建任务通道和等待组
	taskCh := make(chan Task, 10)
	var wg sync.WaitGroup
	// worker 数量
	numWorkers := 5
	// 启动 worker 池
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, taskCh, &wg)
	}
	// 发送任务到任务通道
	for i := 1; i <= 5; i++ {
		taskCh <- Task{
			TableID: i,
			Op:      "洗牌",
		}
		taskCh <- Task{
			TableID: i,
			Op:      "发牌",
		}
	}
	// 关闭任务通道并等待所有 worker 完成
	close(taskCh)
	wg.Wait()

	fmt.Println("所有任务结束")
}
