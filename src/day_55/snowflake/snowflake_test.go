package snowflake

import (
	"testing"
	"time"
)

// TestNew 测试创建雪花算法生成器
func TestNew(t *testing.T) {
	// 测试正常创建
	generator, err := New(1)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}
	if generator == nil {
		t.Fatal("Generator should not be nil")
	}

	// 测试工作节点ID超出范围
	_, err = New(-1)
	if err == nil {
		t.Error("Should return error for worker ID < 0")
	}

	_, err = New(1024) // workerMax is 1023
	if err == nil {
		t.Error("Should return error for worker ID > workerMax")
	}
}

// TestNewWithEpoch 测试创建带自定义epoch的雪花算法生成器
func TestNewWithEpoch(t *testing.T) {
	customEpoch := int64(1609459200000) // 2021-01-01
	generator, err := NewWithEpoch(1, customEpoch)
	if err != nil {
		t.Fatalf("Failed to create generator with custom epoch: %v", err)
	}
	if generator.epoch != customEpoch {
		t.Errorf("Expected epoch %d, got %d", customEpoch, generator.epoch)
	}
}

// TestNextID 测试生成雪花ID
func TestNextID(t *testing.T) {
	generator, err := New(1)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	// 生成多个ID
	for i := 0; i < 10; i++ {
		id, err := generator.NextID()
		if err != nil {
			t.Errorf("Failed to generate ID at index %d: %v", i, err)
		}
		if id == 0 {
			t.Errorf("Generated ID should not be 0 at index %d", i)
		}
	}
}

// TestParseID 测试解析雪花ID
func TestParseID(t *testing.T) {
	generator, err := New(1)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	id, err := generator.NextID()
	if err != nil {
		t.Fatalf("Failed to generate ID: %v", err)
	}

	timestamp, workerID, sequence := ParseID(id, defaultEpoch)
	if timestamp.Before(time.Now().Add(-time.Second)) {
		t.Error("Timestamp should be recent")
	}
	if workerID != 1 {
		t.Errorf("Expected workerID 1, got %d", workerID)
	}
	if sequence < 0 || sequence >= (1<<seqBits) {
		t.Errorf("Sequence out of range")
	}
}

// TestIDMethods 测试ID类型的方法
func TestIDMethods(t *testing.T) {
	id := ID(123456789)

	// 测试String方法
	str := id.String()
	if str != "123456789" {
		t.Errorf("Expected string '123456789', got '%s'", str)
	}

	// 测试Int64方法
	int64Val := id.Int64()
	if int64Val != 123456789 {
		t.Errorf("Expected int64 123456789, got %d", int64Val)
	}
}
