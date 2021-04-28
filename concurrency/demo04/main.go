package main

import (
	"fmt"
	"time"
)

func recv(c chan int) {
	ret := <-c
	fmt.Println(ret)
}

func main() {
	done := make(chan bool)
	go func() {
		for i := 0; i < 3; i++ {
			time.Sleep(100 * time.Millisecond)
			fmt.Println("hello world")
		}
		done <- true
	}()

	<-done //通道会阻塞
	//通过这样的方式，可以实现，主协程等待前一个协程执行结束的效果
	fmt.Println("over!")

	ch := make(chan int)
	go recv(ch)
	var a int
	fmt.Scanln(&a)
	ch <- a
	fmt.Print("success")
}
