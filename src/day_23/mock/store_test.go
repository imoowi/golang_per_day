package mock

import "testing"

func TestDoSave(t *testing.T) {
	ms := &mockStore{}
	DoSave(ms)
	if !ms.called {
		t.Fatal("方法Save()需要被调用")
	}
}
