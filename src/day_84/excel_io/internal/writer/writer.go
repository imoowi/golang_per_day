package writer

import (
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
)

type Writer interface {
	Write(data [][]string) error
}

type ExcelWriter struct {
}

func (e *ExcelWriter) Write(filename string) error {
	start := time.Now().Unix()
	fmt.Println(`start from`, start)
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	totalSize := 10000000 //一千万
	// 将一千万拆成组，每组 100 万行
	groupSize := 1000000 // 每组 100 万行
	totalGroups := totalSize / groupSize

	for g := 0; g < totalGroups; g++ {
		sheet := "Sheet" + fmt.Sprintf("%d", g+1)
		f.NewSheet(sheet)
		sw, err := f.NewStreamWriter(sheet)
		if err != nil {
			return err
		}

		// 批量写入 100 万行
		for row := 0; row <= groupSize; row++ {
			axis, _ := excelize.CoordinatesToCellName(1, row)
			sw.SetRow(axis, []interface{}{row, row + 1, row + 2, row + 3, row + 4, row + 5, row + 6, row + 7, row + 8, `Codee君`})
		}

		// 每写完一组就 Flush 一次，防止内存暴涨
		sw.Flush()
	}

	err := f.SaveAs(filename)
	if err != nil {
		fmt.Println(err)
		return err
	}
	end := time.Now().Unix()
	fmt.Println(`end at`, end)
	fmt.Println(`cost`, end-start, `seconds`)
	return nil
}

func NewExcelWriter() *ExcelWriter {
	return &ExcelWriter{}
}
