package main

import (
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	consulapi "github.com/hashicorp/consul/api"
	"os"
	"github.com/Ballwang/tugo/soft/softServer"
	"github.com/Ballwang/mcserver/gen-go/UserAgent"
	"github.com/Ballwang/tugo/service/userAgentService"
	"strconv"
	"time"
	"math/rand"
	"github.com/Ballwang/tugo/tool"
)

var serverPort=8091

func main() {
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	serverTransport, err := thrift.NewTServerSocket(":"+strconv.Itoa(serverPort))
	if err != nil {
		fmt.Print("Error", err)
		os.Exit(1)
	}
	hander := &userAgentService.WebAgent{}
	processor := UserAgent.NewUserAgentProcessor(hander)
	server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)


	//配置注册服务器信息
	ip:=tool.GetIP()
	rand.Seed(time.Now().Unix())
	registration :=new(consulapi.AgentServiceRegistration)
	registration.ID="user_agent"
	registration.Name="user_agent"
	registration.Address=ip
	registration.Port=serverPort
	registration.Tags=[]string{"浏览器代理控制服务"}
	registration.Check=&consulapi.AgentServiceCheck{
		//TCP:                          fmt.Sprintf("http://%s:%d", registration.Address, checkPort),
		TCP:ip+":"+strconv.Itoa(serverPort),
		Timeout:                        "3s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "5s", //check失败后10秒删除本服务
	}
	softServer.RegisterMcService(registration)
	println("Server starting at "+ip+":"+strconv.Itoa(serverPort))
	server.Serve()
}
