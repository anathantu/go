package main

import "fmt"

type mystruct struct {
	name string
	age  int
}

func test(fn func(string) string, ms mystruct) {

	return
}

func ff(s string) string {
	return s
}

func main() {
	my := mystruct{
		name: "tutu",
		age:  10,
	}

	defer func() {
		fmt.Print(recover())
	}()

	test(ff, my)
	panic("error")
}
