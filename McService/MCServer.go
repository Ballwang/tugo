package McService

import (
	consulapi "github.com/hashicorp/consul/api"
)

func RegisterMcService() (err error){
	//注册服务
	//配置consul
	//这里可以写成全局变量,配置信息里面写consul
	//精简客户端调用参数变量编写
	consulConfig:=&consulapi.Config{Address:"192.168.3.36:8500",Scheme:"http"}
	//创建consul 客户端
	client, err:=consulapi.NewClient(consulConfig)
	if err !=nil {
		//打印错误
		println("Consul 服务器连接失败！",err)
	}

	//配置注册服务器信息
	registration :=new(consulapi.AgentServiceRegistration)
	registration.ID="service_url3"
	registration.Name="service_url_name2"
	registration.Address="192.168.4.50"
	registration.Port=9090
	registration.Tags=[]string{"URL服务第三版本"}
	registration.Check=&consulapi.AgentServiceCheck{
		//TCP:                           fmt.Sprintf("http://%s:%d", registration.Address, checkPort),
		TCP:                           "192.168.3.36:22",
		Timeout:                        "3s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s", //check失败后10秒删除本服务
	}

	err:=client.Agent().ServiceRegister(registration)
	if err !=nil{
		println(registration.ID+"服务注册失败")

	}
	return
	
}

