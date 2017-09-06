package main

import (
	"github.com/garyburd/redigo/redis"
	"github.com/Ballwang/tugo/config"
	"fmt"
)

func main()  {
	cf:=config.NewConfig()
	c,err:=redis.Dial("tcp",cf.GetConfig("redis","redisHost")+":"+cf.GetConfig("redis","redisPort"))
	if err !=nil{
		println("Redis Connect Failed!")
		return
	}
	_ , err =c.Do("HSET", "website","google", "3")
	if err !=nil{
		fmt.Println("pub err:",err)
	}

	defer c.Close()
	
}
