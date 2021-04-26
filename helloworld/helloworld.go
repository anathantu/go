package main // 声明 main 包，表明当前是一个可执行程序

import "fmt" // 导入内置 fmt

func main() { // main函数，是程序执行的入口

	m := make(map[string]string, 4)
	m["11"] = "1111"
	m["22"] = "2222"
	m["33"] = "3333"
	m["44"] = "4444"
	for k, v := range m {
		fmt.Println(k + " " + v)
	}

	fmt.Println("Hello World!") // 在终端打印 Hello World!
}
