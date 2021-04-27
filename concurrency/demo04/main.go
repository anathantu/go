package main

import (
	"fmt"
)

func recv(c chan int) {
	ret := <-c
	fmt.Println(ret)
}

func main() {
	ch := make(chan int)
	go recv(ch)
	ch <- 10
	fmt.Print("success")
}
