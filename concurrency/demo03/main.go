package main

import (
	"fmt"
	"runtime"
)

func counter(out chan<- int) {
	for i := 0; i < 100; i++ {
		out <- i
	}
	close(out)
}

func squarer(out chan<- int, in <-chan int) {
	for i := range in {
		out <- i * i
	}
	close(out)
}

func printer(in <-chan int) {
	for i := range in {
		fmt.Println(i)
	}
}

func main() {

	go func(s string) {
		for i := 0; i < 2; i++ {
			fmt.Println(s)
		}
	}("hello world")
	for i := 0; i < 2; i++ {
		//让cpu让出时间片
		runtime.Gosched()
		fmt.Println("hello")
	}

	ch1 := make(chan int)
	ch2 := make(chan int)
	go counter(ch1)
	go squarer(ch2, ch1)
	printer(ch2)
}
