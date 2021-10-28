package main

import (
	"fmt"
	"time"
)

func main() {
	fun2()
}

func fun2() {
	ch := make(chan int)
	go func() {
		for {
			fmt.Println("func start")
			select {
			case num := <-ch:
				fmt.Println(num)
			}
		}
	}()

	go func() {
		fmt.Println("sleep")
		time.Sleep(5 * time.Second)
		ch <- 1
		ch <- 1
	}()

	for {

	}
}
