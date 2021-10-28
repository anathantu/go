package main

import (
	"context"
	"fmt"
	"github/anathantu/go/"
	"time"
)

func main() {

}

func contextFunc() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	go HelloHandle(ctx, 2*time.Second)
	select {
	case <-ctx.Done():
		fmt.Println("Hello Handle ", ctx.Err())
	}
}

func HelloHandle(ctx context.Context, duration time.Duration) {

	select {
	case <-ctx.Done():
		fmt.Println(ctx.Err())
	case <-time.After(duration):
		fmt.Println("process request with", duration)
	}

}
