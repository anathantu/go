package main

import "fmt"

func getCircleArea(r float32) (ans float32) {
	if r < 0 {
		panic("半径小于0")
	}
	return 3.14 * r * r
}

func test01() {
	a := [5]int{0, 1, 2, 3, 4}
	a[1] = 123
	fmt.Println(a)
	//a[10] = 11
	index := 10
	a[index] = 10
	fmt.Println(a)
}

func test03() {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	//getCircleArea(-1.0)
	test01()
}

func test04() {
	test03()
	fmt.Println("test04")
}

func main() {
	test04()
}
