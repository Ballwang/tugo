package main

import (
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/Ballwang/mcserver/gen-go/business"
	"github.com/Ballwang/tugo/url"
	consulapi "github.com/hashicorp/consul/api"
	"os"
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



	println("Server starting at 192.168.3.50:9090")
	server.Serve()
}
