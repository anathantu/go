package context_with_value

import (
	"context"
	"fmt"
)

type key string

func contextWithValue() {
	ctx := context.WithValue(context.Background(), key("asong"), "Golang梦工厂")
	Get(ctx, key("asong"))
	Get(ctx, key("song"))
}

func Get(ctx context.Context, k key) {
	if v, ok := ctx.Value(k).(string); ok {
		fmt.Println(v)
	}
}
