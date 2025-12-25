package pubsub

import "sync"

type Broker struct {
	subs []chan string
	mu   sync.Mutex
}

// 订阅消息
func (b *Broker) Sub() chan string {
	b.mu.Lock()
	defer b.mu.Unlock()
	ch := make(chan string, 10)
	b.subs = append(b.subs, ch)
	return ch
}

// 发布消息
func (b *Broker) Pub(msg string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, s := range b.subs {
		select {
		case s <- msg:
		default:
		}
	}
}
