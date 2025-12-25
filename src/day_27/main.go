package main

import (
	"log"
	"math"
	"net/http"
	_ "net/http/pprof"
)

// 这个函数消耗大量的CPU计算
func slow() {
	for i := 0; i < 1e7; i++ {
		math.Sqrt(float64(i))
	}
}

var data [][]byte

func leak() {
	b := make([]byte, 10<<20) // 10MB
	data = append(data, b)    // 不断增长
}
func main() {
	//开启pprof的方式
	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		slow()
		leak()
		w.Write([]byte("done"))
	})
	http.ListenAndServe(":8080", nil)
}
