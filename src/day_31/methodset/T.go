package methodset
import "errors"
type ICodeeJun interface{ B() }
type CodeeJun struct {
	Name string
}

func (t CodeeJun) A() {
}

func (t *CodeeJun) B() {
}

type CjWriter struct{}

func (CjWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

type MyError struct{}

func (MyError) Error() string {
	return "MyError"
}
func MyErrorFunc() error {
	var p *MyError = nil
	return p
}
//定义一个存储接口
type Storage interface {
    Set(key string, data []byte) error
    Get(key string) ([]byte, error)
}
//定义一个文件存储结构体
type FileStorage struct {	}

func (f *FileStorage) Set(key string, data []byte) error {
	return nil
}

func (f *FileStorage) Get(key string) ([]byte, error) {
	return nil, nil
}
//定义一个内存存储结构体
type MemStorage struct {	}

func (m *MemStorage) Set(key string, data []byte) error {
	return nil
}

func (m *MemStorage) Get(key string) ([]byte, error) {
	return nil, nil
}
//定义一个Redis存储结构体
type RedisStorage struct {	}

func (r *RedisStorage) Set(key string, data []byte) error {
	return nil
}

func (r *RedisStorage) Get(key string) ([]byte, error) {
	return nil, nil
}
//定义一个存储工厂结构体
type StorageFactory struct {	}

func (s *StorageFactory) Create(name string) (Storage, error) {
	switch name {
	case "file":
		return &FileStorage{}, nil	
	case "mem":
		return &MemStorage{}, nil
	case "redis":
		return &RedisStorage{}, nil
	default:
		return nil, errors.New("invalid storage name")
	}
}
func DoUpload(storage Storage, key string, data []byte) error {
	return storage.Set(key, data)
}