package methodset

import (
	"io"
	"reflect"
	"testing"
)

//*
	func TestMethodSet(t *testing.T) {
		var cj ICodeeJun = CodeeJun{Name: "A"}
		var cj2 ICodeeJun = &CodeeJun{Name: "B"}
		cj.B()
		cj2.B()
	}

//*/
func TestMethodSet2(t *testing.T) {
	codeeJun := &CodeeJun{Name: "Codee君"}
	var cj ICodeeJun = codeeJun
	cj.B()
}

func TestMethodSet3(t *testing.T) {
	var w io.Writer = CjWriter{}
	w.Write([]byte("Hello, CodeeJun!"))
}

func TestMethodSet4(t *testing.T) {
	err := MyErrorFunc()
	if err != nil {
		t.Error("MyErrorFunc() should return nil")
	}
}

func TestMethodSet5(t *testing.T) {
	rv := reflect.ValueOf(CodeeJun{Name: "Codee君"})
	t.Log(rv.Type().Name())
	t.Log(rv.NumMethod())

}

func TestStorage(t *testing.T) {
	factory := &StorageFactory{}
	storage, err := factory.Create("file")
	err = DoUpload(storage, "test.txt", []byte("Hello, CodeeJun!"))
	if err!= nil {
		t.Error(err)
	}
	storage, err = factory.Create("mem")
	err = DoUpload(storage, "test.txt", []byte("Hello, CodeeJun!"))
	if err!= nil {
		t.Error(err)
	}
	storage, err = factory.Create("redis")
	err = DoUpload(storage, "test.txt", []byte("Hello, CodeeJun!"))
	if err!= nil {
		t.Error(err)
	}
}