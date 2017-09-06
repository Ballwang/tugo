package main

import (
	"github.com/Ballwang/tugo/tool"
	"github.com/garyburd/redigo/redis"
	"github.com/Ballwang/tugo/config"
	"time"
	"fmt"
)

func main() {
	fmt.Println("关键词过滤系统开始运行...")

	//无限循环
	for {
		//startTime := tool.CurrentTimeMillis()

		//获取所有关键词
		badWord := tool.RedisSMEMBERS(config.BadWordSet)

		c, _ := tool.NewRedis()

		//遍历采集总数
		for {
			//从集合中取出新采集内容进行匹配
			url,_:=c.Do("SPOP",config.ValueSet)
			if url!=nil {
				urlString,_:=redis.String(url,nil)
				html, _ := c.Do("GET", "Value:-"+urlString)
				if html != nil && len(badWord) >= 0 {
					htmlContent, _ := redis.String(html, nil)
					i:=0
					//过滤关键词
					for _, v := range badWord {
						if tool.FindBadWord(htmlContent, v) {
							c.Do("SADD",config.BadWordStoreSet,urlString)
							i++
						}
					}
					//判断是否有敏感词存在
					if i==0{
						c.Do("SADD",config.NullBadWordSet,urlString)
					}else {
						//删除采集内容
						c.Do("del","Value:-"+urlString)
					}
					//统一写入历史
					md5String:=tool.Md5String(urlString)
					c.Do("SADD",config.HistoryUrlSet,md5String)
				}
			}else {
				break
			}
		}
		c.Close()
		//endTime := tool.CurrentTimeMillis()
		//fmt.Printf("本次调用用时:%d-%d=%d毫秒\n", endTime, startTime, (endTime - startTime))
		time.Sleep(1*time.Second)
	}
}
