package main

import (
	"fmt"
	"time"
)

/**
为一个普通函数创建goroutine的写法
go 函数名（参数列表）
函数名：要调用的函数
参数列表：调用函数需要传入的参数

这种调用方式下，被调用函数的返回值会被忽略
*/
func running() {
	var times int
	for times < 10 {
		times++
		fmt.Println("tick", times)

		time.Sleep(time.Second)
	}
}

func main() {
	go running()

	var input string
	fmt.Scanln(&input)
	fmt.Print(input)
}
