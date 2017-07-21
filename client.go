package main

import (
	"fmt"
	"time"
	"github.com/Ballwang/tugo/McService"

)

const (
	HOST = "192.168.4.50"
	PORT = "9090"
)

func main() {
	startTime := currentTimeMillis()
	client,transport,err:=McService.NewMcClient("service_url3")
	defer 	transport.Close()
	for i := 0; i < 10000; i++ {
		if err !=nil{
			println("service_url12222 client creat failed!",err)
		}
		client.Add(10,30)
	}
	endTime := currentTimeMillis()
	fmt.Printf("本次调用用时:%d-%d=%d毫秒\n", endTime, startTime, (endTime - startTime))
}


func currentTimeMillis() int64 {
	return time.Now().UnixNano() / 1000000
}





