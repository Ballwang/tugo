package tool

import (
	"github.com/Ballwang/tugo/config"

)


//设置服务器状态并且设定过期时间
func SetServerState(serverName string,expTime string)  {
	c,_:=NewRedis()
	defer c.Close()
	c.Do("SET",config.ServerState+serverName,"1")
	c.Do("EXPIRE",config.ServerState+serverName,expTime)

}

//获取服务器状态
func GetServerState(serverName string) bool {
	c,_:=NewRedis()
	defer c.Close()
	reply,_:=c.Do("GET",config.ServerState+serverName)
	if reply!=nil{
		return true
	}
	return false

}
