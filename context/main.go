package main

import (
	"context"
	"fmt"
	"github/anathantu/go/context/context_with_cancel"
	"time"
)

func main() {
	//context_with_value.ContextWithValue()
	//context_with_cancel.ContextWithCancel()
	context_with_cancel.WithCancel()
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
