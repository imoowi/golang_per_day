package pubsub

import (
	"fmt"
	"testing"
)

func TestPubsub(t *testing.T) {
	b := &Broker{}
	s1 := b.Sub()
	s2 := b.Sub()
	s3 := b.Sub()
	go func() {
		for m := range s1 {
			fmt.Println("第一个订阅者获得消息", m)
		}
	}()
	go func() {
		for m := range s2 {
			fmt.Println("第2个订阅者获得消息", m)
		}
	}()
	go func() {
		for m := range s3 {
			fmt.Println("第3个订阅者获得消息", m)
		}
	}()
	b.Pub("张三打了一张牌：1万")
	b.Pub("李四碰了张三的一万")
	b.Pub("王五托管了游戏")
	close(s1)
	close(s2)
	close(s3)
}
