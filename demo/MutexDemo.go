package main

import (
	"fmt"
	"runtime"
	"sync"
)

/**
sync.Mutex为互斥锁（也叫全局锁），Lock()加锁，Unlock()解锁。


*/
var wg = &sync.WaitGroup{}
var lock sync.Mutex
var s = 1000

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	for i := 0; i <= 10; i++ {
		wg.Add(1)
		go add(1)
	}
	wg.Wait()
	//wg.Wait()
	print(s)
	print("====================================")

}
func add(count int) {
	lock.Lock()
	fmt.Printf("加锁----第%d个携程\n", count)
	for i := 0; i < 1; i++ {
		s++
		fmt.Printf("j %d gorount %d \n", s, count)
	}
	fmt.Printf("解锁----第%d个携程\n", count)
	wg.Done()
	defer lock.Unlock()
}
