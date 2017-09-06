package main

import (
	"github.com/Ballwang/tugo/tool"

	"github.com/Ballwang/tugo/config"
	"time"

)

//定时刷新监控队列 平均6秒推送一次,
func main()  {
	params:=config.NewMainParams()
	for {
		c,_:=tool.NewRedis()
		mapString:=tool.RedisHGETALL(params.MonitorSiteHash)
		for _,v:=range mapString{
			c.Do("RPUSH",params.MonitorList,v)
		}
		c.Close()
		time.Sleep(6*time.Second)
	}
}
