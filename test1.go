package main

import (


	"github.com/Ballwang/tugo/tool"

	"github.com/garyburd/redigo/redis"
	"fmt"
)

var completeCat1 chan int = make(chan int)
func main() {
	s:=tool.CurrentTimeMillis()
	//for i:=0;i<=1000;i++ {
	//	go test(i)
	//}
	//
	//for i:=0;i<=1000;i++ {
	//	<-completeCat1
	//}
	c, _ := tool.NewRedis()
	defer c.Close()

	r,_:=redis.String(c.Do("HGET","node:10461","dda7a6812a49ea24f6783b7777416195"))


	if r==""{
		fmt.Println(r)
		fmt.Println("==============")
	}else {
		fmt.Println("--------------")
	}


	e:=tool.CurrentTimeMillis()
	tool.ShowTime(s,e)
}

func test(i int)  {

}


