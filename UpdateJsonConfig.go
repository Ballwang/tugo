package main

import (
	"github.com/Ballwang/tugo/tool"
	"github.com/Ballwang/tugo/config"

	"strconv"
)

func main()  {
	config:=config.NewConfig()
	tool.MakeMarathonJson("P1-Node",0.2,123,2,8087)
	tool.MakeMarathonJson("P2-FindBadMonitorList",0.2,123,2,8088)
	tool.MakeMarathonJson("P3-UpdateMonitorList",0.2,123,2,8088)
	tool.MakeMarathonJson("S4-MonitorCategory",0.2,123,2,8089)
	tool.MakeMarathonJson("P5-UpdateList",0.2,123,2,8088)
	tool.MakeMarathonJson("D6-Category",0.2,123,2,8088)
	tool.MakeMarathonJson("D7-Content",0.2,123,2,8088)
	tool.MakeMarathonJson("D8-BadWord",0.2,123,2,8088)
	tool.MakeMarathonJson("D-Redis",0.2,123,2,8088)
	tool.MakeMarathonJson("E-EsService",0.2,123,2,8088)
	tool.MakeMarathonJson("M-Badword",0.2,123,2,8088)
}

func doJson(c *config.Config,name string)  {
	cpus,_:=strconv.ParseFloat(c.GetConfig(name,"cpus"),32)
	mem,_:=strconv.ParseFloat(c.GetConfig(name,"mem"),32)
	instances,_:=strconv.Atoi(c.GetConfig(name,"instances"))
	port,_:=strconv.Atoi(c.GetConfig(name,"port"))


	tool.MakeMarathonJson(name,cpus,mem,instances,port)
}
