package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/Ballwang/tugo/tool"

	"github.com/Ballwang/tugo/soft/softServer"
	consulapi "github.com/hashicorp/consul/api"
	"strconv"
	"github.com/Ballwang/tugo/config"
)

var tokenRedis = "eo99s001144999999381111"

//推送请求结果
func ShowRequestRedis(w http.ResponseWriter, r interface{}) {
	bytes, _ := json.Marshal(r)
	fmt.Fprint(w, string(bytes))
}

//权限验证
func CheckRightRedis(req *http.Request) bool {
	is := true
	key, isToken := req.Form["token"]
	if !isToken {
		is = false
	} else {
		if key[0] != tokenRedis {
			is = false
		}
	}
	return is
}


//删除Hash列表中元素并且更新Redis中统计数据
func DelHashData(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	ok := CheckRightRedis(req)
	result := make(map[string]interface{})
	if !ok {
		result["success"] = "false"
		result["code"]="401"
		result["message"] = "权限验证失败！"
		ShowRequestRedis(w, result)
		return
	}
	nodeid, okNodeid := req.Form["nodeid"]
	if !okNodeid {
		result["success"] = "false"
		result["code"]="418"
		result["message"] = "nodeid 不能为空"
		ShowRequestRedis(w, result)
		return
	}
	key, okkey := req.Form["key"]
	if !okkey {
		result["success"] = "false"
		result["code"]="418"
		result["message"] = "Key 不能为空"
		ShowRequestRedis(w, result)
		return
	}

	c, _ := tool.NewRedis()
	r, _ := c.Do("HDEL", config.DataPrefix+nodeid[0], key[0])
	data := r.(int64)
	if data > 0 {
		_, err := c.Do("HINCRBY", config.NodeCount, nodeid[0], -1)
		if err != nil {
			result["success"] = "false"
			result["code"]="511"
			result["message"] = nodeid[0] + "：文章删除成功但是统计减一失败！"
			ShowRequestRedis(w, result)
			return
		} else {
			result["success"] = "true"
			result["code"] = "200"
			result["message"] = nodeid[0] + "：文章删除成功！"
			ShowRequestRedis(w, result)
			return
		}
	}else {
		result["success"] = "false"
		result["code"] = "421"
		result["message"] = "该条数据不存在！nodeid:"+nodeid[0] + " ,key:"+key[0]
		ShowRequestRedis(w, result)
		return
	}
}

func main() {
	//第一个参数为客户端发起http请求时的接口名，第二个参数是一个func，负责处理这个请求。
	http.HandleFunc("/DelRedisDataAndCount", DelHashData)

	//服务器要监听的主机地址和端口号
	//配置注册服务器信息

	config:=config.NewConfig()
	serverPort,_:=strconv.Atoi(config.GetConfig("D-Redis","port"))

	ip := tool.GetIP()
	registration := new(consulapi.AgentServiceRegistration)
	registration.ID = "D-Redis:"+ip
	registration.Name = "D-Redis"
	registration.Address = ip
	registration.Port = serverPort
	registration.Tags = []string{"Redis 数据服务接口!"}
	registration.Check = &consulapi.AgentServiceCheck{
		//TCP:                          fmt.Sprintf("http://%s:%d", registration.Address, checkPort),
		TCP:                            ip + ":" + strconv.Itoa(registration.Port),
		Timeout:                        "3s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "5s", //check失败后5秒删除本服务
	}
	softServer.RegisterMcService(registration)
	fmt.Println("Server starting at :" + ip + ":" + strconv.Itoa(registration.Port))
	err := http.ListenAndServe(ip+":"+strconv.Itoa(registration.Port), nil)

	if err != nil {
		fmt.Println("ListenAndServe error: ", err.Error())
	}
}
