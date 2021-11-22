package main

import (
	"fmt"
	"github/anathantu/go/monitor_cache/cache"
	"sync"
	"time"
)

func main() {
	c := cache.New(1048576)
	now := time.Now() //获取当前时间
	fmt.Printf("current time:%v\n", now)

	year := now.Year()     //年
	month := now.Month()   //月
	day := now.Day()       //日
	hour := now.Hour()     //小时
	minute := now.Minute() //分钟
	second := now.Second() //秒
	fmt.Printf("%d-%02d-%02d %02d:%02d:%02d\n", year, month, day, hour, minute, second)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		key := "__name__=\"requests_total\", path=\"/status\", method=\"GET\", instance=”10.0.0.1:80”"
		i := 1

		for {
			v := &cache.Sample{
				T: FromTime(time.Now()),
				V: 132145465,
			}
			tmp := fmt.Sprintf("%sid%d", key, i)
			if ok := c.Check(tmp, v); !ok {
				break
			}
			c.Add(tmp, v)
			i++
			i = i % 7000
		}
	}()
	wg.Wait()
}

// FromTime returns a new millisecond timestamp from a time.
func FromTime(t time.Time) int64 {
	return t.Unix()*1000 + int64(t.Nanosecond())/int64(time.Millisecond)
}
