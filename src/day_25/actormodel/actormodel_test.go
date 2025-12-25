package actormodel

import (
	"testing"
	"time"
)

func TestActor(t *testing.T) {

	in := make(chan RoomMsg, 10)
	go roomActor(1, in)
	in <- Enter{"张三"}
	in <- Enter{"李四"}
	in <- Enter{"王五"}
	in <- Play{"张三", "1万"}
	in <- Exit{"王五"}
	time.Sleep(100 * time.Millisecond)
	close(in)
	time.Sleep(time.Second)

}
