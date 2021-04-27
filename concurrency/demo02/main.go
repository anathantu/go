package main

import (
	"fmt"
	"time"
)

func main() {
	go func() {
		var times int
		for times < 100 {
			times++
			fmt.Println("tick", times)

			time.Sleep(time.Second)
		}
	}()

	var input string
	fmt.Scanln(&input)
}
