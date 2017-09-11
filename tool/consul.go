package tool

import (
	"strconv"
	"github.com/Ballwang/tugo/soft/softServer"
	consulapi "github.com/hashicorp/consul/api"
	"fmt"
)

type ConsulRegister struct {
	Id string
	Name string
	Port int
	Tags []string
}

//注册服务到consul集群中
func (c *ConsulRegister)RegisterConsulService()  {
	ip := GetIP()
	registration := new(consulapi.AgentServiceRegistration)
	registration.ID = c.Id
	registration.Name = c.Name
	registration.Address = ip
	registration.Port = c.Port
	registration.Tags = c.Tags
	registration.Check = &consulapi.AgentServiceCheck{
		//TCP:                          fmt.Sprintf("http://%s:%d", registration.Address, checkPort),
		TCP:                            ip + ":" + strconv.Itoa(registration.Port),
		Timeout:                        "3s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "5s", //check失败后5秒删除本服务
	}
	softServer.RegisterMcService(registration)
	fmt.Println("Server starting at :" + ip + ":" + strconv.Itoa(registration.Port))
}
