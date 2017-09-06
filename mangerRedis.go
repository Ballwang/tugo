package main

import (
	"github.com/Ballwang/tugo/config"
	"github.com/Ballwang/tugo/tool"

	"fmt"
)

var params = config.NewMainParams()

func getMonitorList() {
	mapString:=tool.RedisHGETALL(params.MonitorSiteHash)
	fmt.Println(mapString)
}

//初始化采集队列
func setMonitorList() {
	c, _ := tool.NewRedis()
	//添加MonitorSiteHash 表的同时需要添加MonitorHash，HOST项目 初始化只要初始化这个两个 redis 队列就可以
	c.Do("HSET", params.MonitorSiteHash, params.MonitorSiteHash+"-5", "http://www.test.com/1.js")
	c.Do("HSET", params.MonitorHash, "Host:"+"-"+"http://www.test.com/1111.php", "http://www.test.com/1111.php")
	c.Do("HSET", params.MonitorHash, "Time:"+"-"+"http://www.test.com/1111.php", params.MonitorTime)
	c.Do("HSET", params.MonitorHash, "Time:"+"-"+"http://www.test.com/1111.php", params.MonitorTime)
}

func setBadWord(badWord string)  {
	c, _ := tool.NewRedis()
	c.Do("SADD",config.BadWordSet,badWord)
	fmt.Println(tool.RedisSMEMBERS(config.BadWordSet))
}



func main() {
	setBadWord("womeng")
}
