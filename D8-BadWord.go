package main

import (
	"github.com/Ballwang/tugo/tool"
	"github.com/garyburd/redigo/redis"
	"github.com/Ballwang/tugo/config"
	"time"
	"fmt"
	"encoding/json"
	"strings"
	"net/http"
	"strconv"
)

type dataContent struct {
	Nodeid  string
	Url     string
	DataID  string
	Badword []string
}


//过滤关键词
func FilterBadword(w http.ResponseWriter, req *http.Request)  {
	//无限循环
	for {

		//获取所有关键词 支持热更新关键词
		badWord := tool.RedisClusterSMEMBERS(config.BadWordSet)

		c, _ := tool.NewRedisCluster()

		//遍历采集总数
		for {
			//从集合中取出新采集内容进行匹配
			reply, _ := c.Do("LPOP", config.DataFilterList)
			if reply != nil {
				isBadword := false
				//json 转换到 结构体中，自动匹配对应字段，不区分大小写匹配，但是json中字段开头如果小写则不匹配
				var d = &dataContent{}
				content, _ := redis.String(reply, nil)
				//去除空值
				content=strings.Replace(content," ","",-1)
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
						r,_:=redis.String(c.Do("HGET",config.DataPrefix+d.Nodeid,d.DataID))
						if r==""{
							c.Do("HSET", config.DataPrefix+d.Nodeid, d.DataID, content)
							c.Do("HINCRBY",config.NodeCount,d.Nodeid,1)
						}
					}
				} else {
					b, _ := json.Marshal(*d)
					c.Do("RPUSH", config.DataBadWordList, b)
				}
				c.Do("HDEL",config.ContentParentHash,d.Url)
				c.Do("DEL",config.PrefixCategory+d.Url)
			} else {
				break
			}
			tool.SetServerState("D8-BadWord","5")

		}
		c.Close()
		time.Sleep(1 * time.Second)
		tool.SetServerState("D8-BadWord","5")
	}
}

//关键词过滤服务状态监控
func FilterBadwordState(w http.ResponseWriter, req *http.Request)  {
	fmt.Fprint(w,tool.GetServerState("D8-BadWord"))
}

func main() {

	ip := tool.GetIP()
	var serverID = "D8-BadWord:"+ip
	config:=config.NewConfig()
	serverPort,_:=strconv.Atoi(config.GetConfig("D8-BadWord","port"))


	http.HandleFunc("/FilterBadword", FilterBadword)
	http.HandleFunc("/State", FilterBadwordState)
	register := &tool.ConsulRegister{Id: serverID, Name: "D8-BadWord", Port: serverPort, Tags: []string{"D8 关键词过滤服务！"}}
	register.RegisterConsulService()
	err := http.ListenAndServe(ip+":"+strconv.Itoa(serverPort), nil)

	if err != nil {
		fmt.Println("Listen And Serve error: ", err.Error())
	}

}
