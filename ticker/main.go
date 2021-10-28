package main

import (
	"fmt"
	"time"
)

func main() {
	fun1()
}

func fun1() {
	// 1.获取ticker对象
	ticker := time.NewTicker(5 * time.Second)
	// 子协程
	go func() {
		for {
			fmt.Println("打印")
			select {
			case <-ticker.C:
			}
		}
	}()
	for {

	}
}
