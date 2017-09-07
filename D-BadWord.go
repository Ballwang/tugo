package main

import (
	"github.com/Ballwang/tugo/tool"
	"github.com/garyburd/redigo/redis"
	"github.com/Ballwang/tugo/config"
	"time"
	"fmt"
	"encoding/json"
)

type dataContent struct {
	Nodeid  string
	Url     string
	DataID  string
	Badword []string
}

func main() {
	fmt.Println("关键词过滤系统开始运行...")

	//无限循环
	for {

		//获取所有关键词 支持热更新关键词
		badWord := tool.RedisSMEMBERS(config.BadWordSet)

		c, _ := tool.NewRedis()

		//遍历采集总数
		for {
			//从集合中取出新采集内容进行匹配
			reply, _ := c.Do("LPOP", config.DataFilterList)
			if reply != nil {
				isBadword := false
				//json 转换到 结构体中，自动匹配对应字段，不区分大小写匹配，但是json中字段开头如果小写则不匹配
				var d = &dataContent{}
				content, _ := redis.String(reply, nil)
				json.Unmarshal(reply.([]byte), &d)

				//过滤关键词
				for _, v := range badWord {
					if tool.FindBadWord(content, v) {
						isBadword = true
						d.Badword = append(d.Badword, v)
					}
				}
				//判断是否有敏感词存在
				if !isBadword {
					//未查到有敏感词转存采集内容
					if d.Nodeid != "" && d.DataID != "" {
						c.Do("HSET", config.DataPrefix+d.Nodeid, d.DataID, content)
					}
				} else {
					b, _ := json.Marshal(*d)
					c.Do("RPUSH", config.DataBadWordList, b)
				}

			} else {
				break
			}

		}
		c.Close()

		time.Sleep(1 * time.Second)
	}

}
