package main

import (
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/Ballwang/mcserver/gen-go/business"
	"github.com/Ballwang/tugo/url"
	consulapi "github.com/hashicorp/consul/api"
	"os"
	"github.com/Ballwang/tugo/soft/softServer"

)

func main() {
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	serverTransport, err := thrift.NewTServerSocket(":9090")

	if err != nil {
		fmt.Print("Error", err)
		os.Exit(1)
	}
	hander := &url.UrlTestGo{}
	processor := business.NewBusinessProcessor(hander)
	server := thrift.NewTSimpleServer4(processor, serverTransport, transportFactory, protocolFactory)

	//配置注册服务器信息
	registration :=new(consulapi.AgentServiceRegistration)
	registration.ID="service_url3"
	registration.Name="service_url_name2"
	registration.Address="192.168.4.50"
	registration.Port=9090
	registration.Tags=[]string{"URL服务第三版本"}
	registration.Check=&consulapi.AgentServiceCheck{
		//TCP:                          fmt.Sprintf("http://%s:%d", registration.Address, checkPort),
		TCP:                            "192.168.3.36:22",
		Timeout:                        "3s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "5s", //check失败后10秒删除本服务
	}
	softServer.RegisterMcService(registration)
	println("Server starting at 192.168.3.50:9090")
	server.Serve()
}
