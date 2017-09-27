package main

import (
	"github.com/Ballwang/tugo/tool"
	"github.com/Ballwang/tugo/config"

	"strconv"
)

func main()  {
    config:=config.NewConfig()
	servers :=[]string{"P1-Node","P2-FindBadMonitorList","P3-UpdateMonitorList","S4-MonitorCategory","P5-UpdateList","D6-Category","D7-Content","D8-BadWord","D-Redis","E-EsService","M-Badword"}

	for _,v:=range servers{
		doJson(config,v)
	}

}

func doJson(c *config.Config,name string)  {
	cpus,_:=strconv.ParseFloat(c.GetConfig(name,"cpus"),32)
	mem,_:=strconv.ParseFloat(c.GetConfig(name,"mem"),32)
	instances,_:=strconv.Atoi(c.GetConfig(name,"instances"))
	port,_:=strconv.Atoi(c.GetConfig(name,"port"))
	tool.MakeMarathonJson(name,cpus,mem,instances,port)
}
