package main

import (
	"github.com/Ballwang/tugo/tool"
	"github.com/garyburd/redigo/redis"
	"time"
	"github.com/Ballwang/tugo/config"
)

//更新队列定时更新
func main()  {
	//处理集合中URL，更新列表
	params:=config.NewMainParams()
	for {
		c,_:=tool.NewRedis()
		rep,_:=c.Do("SPOP",params.UpdateListSet)
		if rep !=nil {
			url,_:=redis.String(rep,nil)
			c.Do("RPUSH",params.UpdateList,url)
			c.Close()
		}else {
			c.Close()
			time.Sleep(3*time.Second)
		}
	}
}
