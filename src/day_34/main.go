package main

import (
	"container/heap"
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("Hello, Codee君!")

	users := map[int]string{
		100: "Codee君",
		200: "詹姆斯邦德",
		300: "乔治华盛顿",
	}
	res1 := lookupMaybe(users, 100)
	fmt.Println(res1)
	if name, ok := res1.Get(); ok {
		fmt.Println(name)
	}
	res2 := lookupMaybe(users, 011)
	nameOrDefault := res2.Or("未知")
	fmt.Println(nameOrDefault)
	nameOrDefault2 := res1.Or("未知")
	fmt.Println(nameOrDefault2)

	numbers := []int{1, 2, 3, 4, 5}
	// 转换成字符串
	strings := SliceA2B(numbers, func(i int) string {
		return fmt.Sprintf("ID=%d", i)
	})
	fmt.Println(strings)
	// 过滤出偶数
	even := Filter(numbers, func(i int) bool {
		return i%2 == 0
	})
	fmt.Println(even)
	// 栈
	stack := NewStack[string]()
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			value := fmt.Sprintf("value-%d", i)
			stack.Push(&value)
		}(i)
	}
	wg.Wait()
	fmt.Println("栈元素数量:", stack.Len())
	if itm, ok := stack.Pop(); ok {
		fmt.Println(*itm)
		fmt.Println("栈元素数量:", stack.Len())
	}

	fmt.Println(stack.Pop())
	fmt.Println(stack.Pop())
	fmt.Println(stack.Pop())
	fmt.Println(stack.Pop()) // 栈空，返回 nil
	// 映射
	m := NewCodeeMap[int, string]()
	m.Set(1, "ID=1")
	m.Set(2, "ID=2")
	m.Set(3, "ID=3")
	fmt.Println(m.Get(1))
	fmt.Println(m.Get(2))
	fmt.Println(m.Get(3))
	fmt.Println(m.Get(4)) // 不存在，返回零值
	sum := Reduce([]int{1, 2, 3, 4}, 0, func(a, b int) int { return a + b })
	fmt.Println(sum)
	mult := Reduce([]int{1, 2, 3, 4}, 1, func(a, b int) int { return a * b })
	fmt.Println(mult)
	// 优先队列
	pq := make(PriorityQueue[Task], 0)
	tasks := []Task{
		{
			Name:    "任务1",
			DueDate: time.Now(),
		},
		{
			Name:    "任务2",
			DueDate: time.Now().AddDate(0, 0, 2),
		},
		{
			Name:    "任务3",
			DueDate: time.Now().AddDate(0, 0, 1),
		},
	}
	for _, task := range tasks {
		// 将Task包装在Item中，并根据DueDate设置优先级（时间越早优先级越高）
		priority := int(time.Since(task.DueDate).Milliseconds())
		if priority < 0 {
			priority = -priority // 确保过期时间越早，优先级越高
		}
		heap.Push(&pq, &Item[Task]{
			value:    task,
			Priority: priority,
		})
	}
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item[Task])
		fmt.Println(item.value)
	}

}
