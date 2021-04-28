package main

import (
	"fmt"
	"sync"
)

var x int
var wg sync.WaitGroup

var y int
var lock sync.Mutex

func addWithoutLock() {
	for i := 0; i < 5000; i++ {
		x++

	}
	wg.Done()
}

func addWithLock() {
	for i := 0; i < 5000; i++ {
		lock.Lock() //互斥锁
		y++
		lock.Unlock()
	}
	wg.Done()
}

func main() {
	wg.Add(2)
	go addWithoutLock()
	go addWithoutLock()
	wg.Wait()
	fmt.Println(x)

	wg.Add(2)
	go addWithLock()
	go addWithLock()
	wg.Wait()
	fmt.Print(y)
}
