package dynamicworkerpool

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestDynamicWorkerPool(t *testing.T) {
	jobs := make(chan int, 1000)
	var workerQuits []chan struct{}
	stop := make(chan struct{})
	maxWorkers := 10
	var wg sync.WaitGroup

	// 调度器：根据 jobs 队列长度调整 worker 数量（非常简单的策略）
	go func() {
		workers := 0
		for {
			select {
			case <-stop:
				// 停止调度器时，关闭剩余的 worker quit 通道以确保 worker 退出
				for _, q := range workerQuits {
					close(q)
				}
				return
			default:
			}

			log.Println("检测jobs长度")
			l := len(jobs)
			target := 1

			switch {
			case l > 20:
				target = 10
			case l > 10:
				target = 3
			case l > 5:
				target = 2
			}
			target = min(target, maxWorkers)
			if target > workers {
				// 启动更多 worker（为每个 worker 创建独立 quit 通道）
				add := target - workers
				for i := 0; i < add; i++ {
					wg.Add(1)
					q := make(chan struct{})
					workerQuits = append(workerQuits, q)
					go worker(workers+i+1, jobs, q, &wg)
				}
				workers = target
				fmt.Println("扩容 ->", workers, "个工人")
			} else if target < workers {
				// 精确缩容：关闭最近创建的 worker 的 quit 通道，并从列表移除
				for i := 0; i < workers-target; i++ {
					last := len(workerQuits) - 1
					if last >= 0 {
						close(workerQuits[last])
						workerQuits = workerQuits[:last]
					}
				}
				workers = target
				fmt.Println("缩容 ->", workers, "个工人")
			}
			time.Sleep(200 * time.Millisecond)
		}
	}()

	// 产生任务（模拟高峰 / 低谷）
	go func() {
		for i := 1; i <= 600; i++ {
			jobs <- i
			time.Sleep(time.Duration(rand.Intn(80)) * time.Millisecond)
		}
		close(jobs)
	}()

	// 等待所有 workers 结束（注意：简化版，这里睡一会儿以等待）
	time.Sleep(6 * time.Second)
	// 停止调度器（它会关闭剩余的 worker quit 通道）
	close(stop)
	wg.Wait()
	fmt.Println("结束")
}
