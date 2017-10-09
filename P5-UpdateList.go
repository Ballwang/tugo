package main

import (
	"github.com/Ballwang/tugo/tool"
	"time"
	"github.com/Ballwang/tugo/config"
	"net/http"
	"strconv"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

//迁移有变动的网站链接 6秒一次
func StartUpdateList(w http.ResponseWriter, req *http.Request) {

	for {
		c, _ := tool.NewRedisCluster()
		for {
			string, err := redis.String(c.Do("SPOP", config.UpdateListSet))

			if err != nil {

			}

			if string != "" {
				c.Do("RPUSH", config.UpdateList, string)
			} else {
				break
			}
		}

		c.Close()
		time.Sleep(6 * time.Second)
		tool.SetServerState("P5-UpdateList", "8")
	}
}

//服务运行状态监控
func UpdateListState(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, tool.GetServerState("P5-UpdateList"))
}

//更新队列
func main() {

	ip := tool.GetIP()
	var serverID = "P5-UpdateList:"+ip
	config:=config.NewConfig()
	serverPort,_:=strconv.Atoi(config.GetConfig("P5-UpdateList","port"))
	http.HandleFunc("/P5-UpdateList", StartUpdateList)
	http.HandleFunc("/State", UpdateListState)
	register := &tool.ConsulRegister{Id: serverID, Name: "P5-UpdateList", Port: serverPort, Tags: []string{"P5 能够迁移有变动的列表链接到带采集队列中！"}}
	register.RegisterConsulService()
	err := http.ListenAndServe(ip+":"+strconv.Itoa(serverPort), nil)

	if err != nil {
		fmt.Println("Listen And Serve error: ", err.Error())
	}

}
