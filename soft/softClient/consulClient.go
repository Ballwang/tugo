package softClient

import (
	"strconv"
	"os"
	"net"
	"fmt"
	"github.com/Ballwang/mcserver/gen-go/business"
	 consulapi "github.com/hashicorp/consul/api"
	//"github.com/apache/thrift/lib/go/thrift"
	"github.com/Ballwang/tugo/config"
	"github.com/Ballwang/mcserver/gen-go/UserAgent"
	"github.com/apache/thrift/lib/go/thrift"
)

//根据注册的servicID 查找服务详细地址
func getMcServiceConfigFromConsul(servicID string) (is bool,ip string,port string) {
	//这里可以写成全局变量,配置信息里面写consul
	//精简客户端调用参数变量编写
	consulConfig:=config.NewConfig()
	consulHost:=consulConfig.GetConfig("consul","consulHost")
	consulPort:=consulConfig.GetConfig("consul","consulPort")
	consulScheme:=consulConfig.GetConfig("consul","Scheme")
	Address:=consulHost+":"+consulPort
	consulConfigResult:=&consulapi.Config{Address:Address,Scheme:consulScheme}
	client,err := consulapi.NewClient(consulConfigResult)
	if err != nil{
		println("client connect error",err)
	}
	services, err := client.Agent().Services()
	if err != nil{
		println("consul get error",err)
	}

	if _,found := services[servicID];!found{
		println("Servic :"+servicID+" NOT FOUND!!")
		return false,"",""
	}else {
		return true,services[servicID].Address,strconv.Itoa(services[servicID].Port)
	}
}

//获取注册的微服务，并且创建微服务客户端
func NewMcClient(serverID string) (*business.BusinessClient,*thrift.TSocket,error) {
	mc,mcIP,mcPort:=getMcServiceConfigFromConsul(serverID)
	if !mc{
		println("Mc service is down!!!")
		os.Exit(1)
	}
    //创建 thrift 客户端
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	transport, err := thrift.NewTSocket(net.JoinHostPort(mcIP, mcPort))
	if err != nil {
		fmt.Println(os.Stderr, "errror resolving address", err)
		os.Exit(1)
	}
	userTransport := transportFactory.GetTransport(transport)
	client := business.NewBusinessClientFactory(userTransport, protocolFactory)
	if err := transport.Open(); err != nil {
		fmt.Println(os.Stderr, "Error opening socket to "+mcIP+":"+mcPort, "", err)
		os.Exit(1)
	}
	//defer transport.Close()
	return client,transport,err
}

//获取注册的微服务，并且创建微服务客户端
func NewMcUserAgentClient(serverID string) (*UserAgent.UserAgentClient,*thrift.TSocket,error) {
	mc,mcIP,mcPort:=getMcServiceConfigFromConsul(serverID)
	if !mc{
		println("Mc service is down!!!")

	}
	//创建 thrift 客户端
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	transport, err := thrift.NewTSocket(net.JoinHostPort(mcIP, mcPort))
	if err != nil {
		fmt.Println(os.Stderr, "errror resolving address", err)

	}
	userTransport := transportFactory.GetTransport(transport)
	client := UserAgent.NewUserAgentClientFactory(userTransport, protocolFactory)
	if err := transport.Open(); err != nil {
		fmt.Println(os.Stderr, "Error opening socket to "+mcIP+":"+mcPort, "", err)

	}
	//defer transport.Close()
	return client,transport,err
}