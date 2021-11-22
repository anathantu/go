package main

import (
	"fmt"
	"net/http"
)

func sayHello(w http.ResponseWriter, r *http.Request) {
	fmt.Println("【say hello】访问成功")
}

func sayHello2(w http.ResponseWriter, r *http.Request) {
	fmt.Println("【say hello2】访问成功")
}

func main() {
	http.HandleFunc("/metrics", sayHello)
	http.HandleFunc("/find", sayHello2)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("http server failed, err:%v\n", err)
		return
	}
}
