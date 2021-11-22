package context_with_cancel

import (
	"context"
	"fmt"
	"time"
)

func ContextWithCancel() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go Speak(ctx)
	time.Sleep(10 * time.Second)
}

func Speak(ctx context.Context) {
	for range time.Tick(time.Second) {
		select {
		case <-ctx.Done():
			return
		default:
			fmt.Println("balabalabalabala")
		}
	}
}

func WithCancel() {
	gen := func(ctx context.Context) <-chan int {
		dst := make(chan int)
		n := 1
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case dst <- n:
					n++
				}
			}
		}()
		return dst
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for n := range gen(ctx) {
		fmt.Println(n)
		if n == 5 {
			break
		}
	}
}
