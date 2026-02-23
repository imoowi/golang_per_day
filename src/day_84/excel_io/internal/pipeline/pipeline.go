package pipeline

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/xuri/excelize/v2"
)

type Pipeline interface {
	Run() error
}

type ExcelPipeline struct {
}

func NewExcelPipeline() *ExcelPipeline {
	return &ExcelPipeline{}
}

/*
	func (e *ExcelPipeline) Run(in, out string) error {
		start := time.Now().Unix()
		fmt.Println(`start from`, start)
		r, err := excelize.OpenFile(in)
		if err != nil {
			return err
		}

		// 获取excel的所有sheet
		sheets := r.GetSheetList()
		if len(sheets) == 0 {
			return fmt.Errorf("no sheets found in the Excel file")
		}

		jobs := make(chan []string, 1000)
		results := make(chan []any, 1000)

		for _, sheetName := range sheets {
			rows, err := r.Rows(sheetName)
			if err != nil {
				return err
			}
			defer rows.Close()

			// 读取 Excel → jobs
			go func() {
				for rows.Next() {
					row, _ := rows.Columns()
					jobs <- row
				}
				close(jobs)
			}()
		}

		workerNum := runtime.NumCPU() * 2
		var wg sync.WaitGroup

		// 启动 worker 池
		for i := 0; i < workerNum; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				e.Worker(jobs, results)
			}()
		}

		// 等所有 worker 结束 → 关闭 results
		go func() {
			wg.Wait()
			close(results)
		}()
		w := excelize.NewFile()
		sw, err := w.NewStreamWriter("Sheet1")
		if err != nil {
			return err
		}

		idx := 1
		for res := range results {
			axis, _ := excelize.CoordinatesToCellName(1, idx)
			sw.SetRow(axis, res)
			idx++
		}

		sw.Flush()
		w.SaveAs(out)
		end := time.Now().Unix()
		fmt.Println(`end at`, end)
		fmt.Println(`cost`, end-start, `seconds`)
		return nil
	}

//
*/
func (e *ExcelPipeline) Run(in, out string) error {
	start := time.Now()
	fmt.Println("start:", start)

	r, err := excelize.OpenFile(in)
	if err != nil {
		return err
	}

	sheets := r.GetSheetList()
	if len(sheets) == 0 {
		return fmt.Errorf("no sheets found")
	}

	w := excelize.NewFile()
	w.DeleteSheet("Sheet1") // 删除默认 sheet

	for _, sheetName := range sheets {
		// fmt.Println("processing:", sheetName)

		if err := e.processSheet(r, w, sheetName); err != nil {
			return err
		}
	}

	if err := w.SaveAs(out); err != nil {
		return err
	}

	fmt.Println("done, cost:", time.Since(start))
	return nil
}
func (e *ExcelPipeline) processSheet(r *excelize.File, w *excelize.File, sheet string) error {

	rows, err := r.Rows(sheet)
	if err != nil {
		return err
	}
	defer rows.Close()

	jobs := make(chan []string, 1000)
	results := make(chan []any, 1000)

	workerNum := runtime.NumCPU() * 2
	var wg sync.WaitGroup

	// worker pool
	for i := 0; i < workerNum; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			e.Worker(jobs, results)
		}()
	}

	// 生产者
	go func() {
		for rows.Next() {
			row, _ := rows.Columns()
			jobs <- row
		}
		close(jobs)
	}()

	// 回收者
	go func() {
		wg.Wait()
		close(results)
	}()

	// 创建输出 sheet
	outSheet := sheet
	w.NewSheet(outSheet)
	sw, err := w.NewStreamWriter(outSheet)
	if err != nil {
		return err
	}

	idx := 1
	for res := range results {
		axis, _ := excelize.CoordinatesToCellName(1, idx)
		sw.SetRow(axis, res)
		idx++
	}

	return sw.Flush()
}

// Worker
func (e *ExcelPipeline) Worker(jobs <-chan []string, results chan<- []any) {
	for job := range jobs {
		row := make([]any, 0, len(job))
		for _, cell := range job {
			row = append(row, cell)
		}
		results <- row
	}
}
