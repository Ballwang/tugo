package softServer

import (
	consulapi "github.com/hashicorp/consul/api"
	"github.com/Ballwang/tugo/config"
)

//注册新服务到Consul 中
func RegisterMcService(registration *consulapi.AgentServiceRegistration) (err error) {
	//注册服务
	//配置consul
	//这里可以写成全局变量,配置信息里面写consul
	//精简客户端调用参数变量编写
	consulConfig := config.NewConfig()
	consulHost := consulConfig.GetConfig("consul", "consulHost")
	consulPort := consulConfig.GetConfig("consul", "consulPort")
	consulScheme := consulConfig.GetConfig("consul", "Scheme")
	Address := consulHost + ":" + consulPort
	consulConfigResult :=&consulapi.Config{Address: Address, Scheme: consulScheme}
	//创建 consul 客户端
	client, err := consulapi.NewClient(consulConfigResult)
	if err != nil {
		//打印错误
		println("Consul 服务器连接失败！", err)
		return
	}
	//注册详细信息到 Consul 集群中
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		println(registration.ID + "服务注册失败")
	}

	return
}
