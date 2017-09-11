package main

import (
	"github.com/Ballwang/tugo/tool"
	"github.com/garyburd/redigo/redis"
	"time"
	"github.com/Ballwang/tugo/config"
	"net/http"
	"strconv"
	"fmt"
)


//迁移有变动的网站链接 6秒一次
func StartUpdateList(w http.ResponseWriter, req *http.Request)  {

	//for {
	//	c,_:=tool.NewRedis()
	//	mapString:=tool.RedisHGETALL(params.MonitorSiteHash)
	//	for _,v:=range mapString{
	//		c.Do("RPUSH",params.MonitorList,v)
	//	}
	//	c.Close()
	//	time.Sleep(6*time.Second)
	//	tool.SetServerState("P5-UpdateList","7")
	//}

	for{
		c,_:=tool.NewRedis()
		string:=tool.RedisSMEMBERS(config.UpdateListSet)
		if len(string)>0{
			for _,v:=range string{
				c.Do("RPUSH",params.MonitorList,v)
			}
		}
		c.Close()
		time.Sleep(6*time.Second)
		tool.SetServerState("P5-UpdateList","8")
	}
}

//服务运行状态监控
func UpdateListState(w http.ResponseWriter, req *http.Request)  {
	fmt.Fprint(w,tool.GetServerState("P5-UpdateList"))
}
//更新队列
func main()  {

	var serverID = "P5-UpdateList"
	var serverPort = 8092
	ip := tool.GetIP()
	http.HandleFunc("/StartUpdateList", StartUpdateList)
	http.HandleFunc("/State", UpdateListState)
	register := &tool.ConsulRegister{Id: serverID, Name: "P5-迁移变动列表服务", Port: serverPort, Tags: []string{"P5 能够迁移有变动的列表链接到带采集队列中！"}}
	register.RegisterConsulService()
	err := http.ListenAndServe(ip+":"+strconv.Itoa(serverPort), nil)

	if err != nil {
		fmt.Println("Listen And Serve error: ", err.Error())
	}


}
