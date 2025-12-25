package services

import "sync/atomic"

var ready int32

func SetReady() {
	atomic.StoreInt32(&ready, 1)
}
func SetNotReady() {
	atomic.StoreInt32(&ready, 0)
}
func IsReady() bool {
	return atomic.LoadInt32(&ready) == 1
}
