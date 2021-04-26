package main

import "fmt"

func main() {
	var str string
	str = "aaaa"
	for i, n := 0, len(str); i < n; i++ {
		fmt.Printf("%c\n", str[i])
	}

	arr := [6]int{1, 2, 3, 4, 5, 6}
	l := len(arr) - 1
	for l >= 0 {
		fmt.Println(arr[l])
		l--
	}

	for i, x := range arr {
		fmt.Printf("%d  %d\n", i, x)
		if i > 2 {
			break
		}
	}
}
