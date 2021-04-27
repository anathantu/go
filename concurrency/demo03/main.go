package main

import (
	"fmt"
	"runtime"
)

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
}
