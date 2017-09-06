package tool

import (
	"time"
	"fmt"
)

func CurrentTimeMillis() int64 {
	return time.Now().UnixNano() / 1000000
}

func ShowTime(s int64,e int64)  {
	fmt.Printf("本次调用用时:%d-%d=%d毫秒\n", e, s, (e - s))
}
