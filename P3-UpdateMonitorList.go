package main

import (
	"github.com/Ballwang/tugo/tool"

	"github.com/Ballwang/tugo/config"
	"time"

	"net/http"
	"strconv"
	"fmt"
)


//链接采集频率控制
func StartUpdateMonitorList(w http.ResponseWriter, req *http.Request)  {
	params:=config.NewMainParams()
	for {
		c,_:=tool.NewRedisCluster()
		mapString:=tool.RedisClusterHGETALL(params.MonitorSiteHash)
		for _,v:=range mapString{
			c.Do("RPUSH",params.MonitorList,v)
		}
		c.Close()
		time.Sleep(6*time.Second)
		tool.SetServerState("P3-UpdateMonitorList","7")
	}
}

//链接采集服务 状态监控
func UpdateMonitorState(w http.ResponseWriter, req *http.Request)  {
	fmt.Fprint(w,tool.GetServerState("P3-UpdateMonitorList"))
}




//定时刷新监控队列 平均6秒推送一次,
func main()  {

	var serverID = "P3-UpdateMonitorList"
	var serverPort = 8089
	ip := tool.GetIP()
	http.HandleFunc("/StartUpdateMonitorList", StartUpdateMonitorList)
	http.HandleFunc("/State", UpdateMonitorState)
	register := &tool.ConsulRegister{Id: serverID, Name: "P3-采集频率控制服务", Port: serverPort, Tags: []string{"P3 能够控制特定链接的采集频率，能够根据网站更新频率判断伸缩节点被采集频率！"}}
	register.RegisterConsulService()
	err := http.ListenAndServe(ip+":"+strconv.Itoa(serverPort), nil)

	if err != nil {
		fmt.Println("Listen And Serve error: ", err.Error())
	}

}
