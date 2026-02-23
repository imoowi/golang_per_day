package reader

import (
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
)

type Reader interface {
// Reader 定义了读取器的接口，用于从文件中读取数据
	Read(filename string) ([][]string, error)
}

type ExcelReader struct {
// ExcelReader 是 Reader 接口的具体实现，用于读取 Excel 文件
}

func (e *ExcelReader) Read(filename string) ([][]string, error) {
// Read 从指定的 Excel 文件中读取数据并返回二维字符串数组
// filename: 要读取的 Excel 文件路径
// 返回: 二维字符串数组表示的 Excel 数据，以及可能的错误
	start := time.Now().Unix()
	fmt.Println(`start from`, start)
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	rows, err := f.Rows("Sheet1")
	if err != nil {
		return nil, err
	}
	var data [][]string
	for rows.Next() {
		row, err := rows.Columns()
		if err != nil {
			return nil, err
		}
		data = append(data, row)
	}
	end := time.Now().Unix()
	fmt.Println(`end at`, end)
	fmt.Println(`cost`, end-start, `seconds`)
	return data, nil
}

func NewExcelReader() *ExcelReader {
// NewExcelReader 创建并返回一个新的 ExcelReader 实例
// 返回: 指向新创建的 ExcelReader 的指针
	return &ExcelReader{}
}
