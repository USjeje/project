package utils

import (
	"math/rand"
	"sync"
	"time"
)

var (
	//RandomIntMutex 用于得到随机Int的互斥锁
	RandomIntMutex sync.Mutex
)

//GetRandomInt 返回介于[start, end)随机Int数,要注意并发，要加锁,保证多线程安全
func GetRandomInt(start, end int) int {
	RandomIntMutex.Lock()
	time.Sleep(time.Microsecond)
	rand.Seed(time.Now().UnixNano())

	//这才是返回(start, end]之间的随机数
	res := start + rand.Intn(end - start)
	RandomIntMutex.Unlock()
	return res
}