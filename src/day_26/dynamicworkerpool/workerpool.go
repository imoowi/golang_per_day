package dynamicworkerpool

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func worker(id int, jobs <-chan int, quit <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case j, ok := <-jobs:
			if !ok {
				return
			}
			fmt.Println("工人", id, "处理任务", j)
			time.Sleep(time.Duration(rand.Intn(200)+50) * time.Millisecond)
		case <-quit:
			fmt.Println("工人", id, "收到结束信号")
			return
		}
	}
}
