package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/Ballwang/tugo/config"
	"database/sql"
	"fmt"
	"github.com/Ballwang/tugo/tool"

	"net/http"
	"strconv"
	"github.com/Ballwang/tugo/soft/softServer"
	consulapi "github.com/hashicorp/consul/api"
	"encoding/json"
	"strings"
)

var badToken = "72800634e9d0e3ccc0e32aca1154ff130894879283"

//权限验证
func CheckRightBadword(w http.ResponseWriter, req *http.Request) bool {
	is := true
	key, isToken := req.Form["token"]
	if !isToken {
		is = false
	} else {
		if key[0] != badToken {
			is = false
		}
	}

	return is
}

//更新所有关键词
func AddNewBadword(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	isAccess := CheckRightBadword(w, req)
	result := make(map[string]interface{})
	if !isAccess {
		result["success"] = "false"
		result["code"]="401"
		result["message"] = "权限验证失败"
		bytes, _ := json.Marshal(result)
		fmt.Fprint(w, string(bytes))
		return
	}
	badwordResult, isbad := req.Form["badword"]
	if !isbad {
		result["success"] = "false"
		result["code"]="418"
		result["message"] = "badword 必须填写！"
		bytes, _ := json.Marshal(result)
		fmt.Fprint(w, string(bytes))
		return
	}


	c, err := tool.NewRedis()
	if err != nil {
		result["success"] = "false"
		result["code"]="511"
		result["message"] = "Redis 链接失败！"
		bytes, _ := json.Marshal(result)
		fmt.Fprint(w, string(bytes))
		return
	}


	defer c.Close()
	badword:=badwordResult[0]

	if badword != "" {
		badword = strings.Replace(badword, " ", "", -1)
		badword = strings.Replace(badword, "\n", "", -1)
		badword = strings.Replace(badword, "\r\n", "", -1)
		badword = strings.Replace(badword, "\r", "", -1)
		if badword!=""{
			r,e:=c.Do("SADD", config.BadWordSet, badword)
			if e!=nil{
				result["success"] = "false"
				result["code"]="510"
				result["message"] = "敏感词添加失败！"
				bytes, _ := json.Marshal(result)
				fmt.Fprint(w, string(bytes))
			}
			if r!=nil{
				result["success"] = "true"
				result["success"] = "200"
				result["message"] = "敏感词添加成功！"
				bytes, _ := json.Marshal(result)
				fmt.Fprint(w, string(bytes))
			}
		}else {
			result["success"] = "false"
			result["success"] = "422"
			result["message"] = "请勿添加空值！"
			bytes, _ := json.Marshal(result)
			fmt.Fprint(w, string(bytes))
		}


	}else {
		result["success"] = "false"
		result["code"]="510"
		result["message"] = "敏感词添加失败！"
		bytes, _ := json.Marshal(result)
		fmt.Fprint(w, string(bytes))
	}


}

//更新所有关键词
func UpdateAllBadword(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	isAccess := CheckRightBadword(w, req)
	result := make(map[string]interface{})
	if !isAccess {
		result["success"] = "false"
		result["code"]="401"
		result["message"] = "权限验证失败"
		bytes, _ := json.Marshal(result)
		fmt.Fprint(w, string(bytes))

		return
	}
	params := config.NewConfig()
	mysqlHost := params.GetConfig("mysql", "mysqlHost")
	mysqlPort := params.GetConfig("mysql", "mysqlPort")
	mysqlUser := params.GetConfig("mysql", "mysqlUser")
	mysqlPassword := params.GetConfig("mysql", "mysqlPassword")
	mysqlCharset := params.GetConfig("mysql", "mysqlCharset")
	mysqlDatabase := params.GetConfig("mysql", "mysqlDatabase")
	fmt.Println(mysqlUser + ":" + mysqlPassword + "@tcp(" + mysqlHost + ":" + mysqlPort + ")/" + mysqlDatabase + "?charset=" + mysqlCharset)
	db, err := sql.Open("mysql", mysqlUser+":"+mysqlPassword+"@tcp("+mysqlHost+":"+mysqlPort+")/"+mysqlDatabase+"?charset="+mysqlCharset)
	rows, err := db.Query("SELECT badword FROM js_badword")
	if err != nil {
		result["success"] = "false"
		result["code"]="512"
		result["message"] = "数据库连接失败！"
		bytes, _ := json.Marshal(result)
		fmt.Fprint(w, string(bytes))
		return
	}

	c, err := tool.NewRedis()
	if err != nil {
		result["success"] = "false"
		result["code"]="511"
		result["message"] = "Redis 链接失败！"
		bytes, _ := json.Marshal(result)
		fmt.Fprint(w, string(bytes))
		return
	}

	defer c.Close()

	for rows.Next() {
		var badword string
		if err := rows.Scan(&badword); err != nil {
			result["success"] = "false"
			result["code"]="513"
			result["message"] = "敏感词查询错误！"
			bytes, _ := json.Marshal(result)
			fmt.Fprint(w, string(bytes))
			return
		}

		//网页提交和入库字符串比较坑爹，会带有隐藏字符 空格，换行等
		if badword != "" {
			badword = strings.Replace(badword, " ", "", -1)
			badword = strings.Replace(badword, "\n", "", -1)
			badword = strings.Replace(badword, "\r\n", "", -1)
			badword = strings.Replace(badword, "\r", "", -1)
			c.Do("SADD", config.BadWordSet, badword)
		}
	}
	result["success"] = "true"
	result["code"]="200"
	result["message"] = "敏感词更新成功！"
	bytes, _ := json.Marshal(result)
	fmt.Fprint(w, string(bytes))
}


//删除关键词
func DelBadeword(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	isAccess := CheckRightBadword(w, req)
	result := make(map[string]interface{})
	if !isAccess {
		result["success"] = "false"
		result["code"]="401"
		result["message"] = "权限验证失败"
		bytes, _ := json.Marshal(result)
		fmt.Fprint(w, string(bytes))
		return
	}


	badword, isbad := req.Form["badword"]
	if !isbad {
		result["success"] = "false"
		result["code"]="418"
		result["message"] = "badword 必须填写！"
		bytes, _ := json.Marshal(result)
		fmt.Fprint(w, string(bytes))
		return
	}

	c, err := tool.NewRedis()
	if err != nil {
		result["success"] = "false"
		result["code"]="511"
		result["message"] = "Redis 创建失败！"
		bytes, _ := json.Marshal(result)
		fmt.Fprint(w, string(bytes))
		fmt.Fprint(w, err)
		return
	}

	reply, err := c.Do("SREM", config.BadWordSet, badword[0])
	if err != nil {
		result["success"] = "false"
		result["code"] = "511"
		result["message"] = "Redis 命令执行失败！"
		bytes, _ := json.Marshal(result)
		fmt.Fprint(w, string(bytes))
		return
	}

	r := reply.(int64)
	if r > 0 {
		result["success"] = "true"
		result["code"] = "200"
		result["message"] = badword[0] + "敏感删除成功！"
		bytes, _ := json.Marshal(result)
		fmt.Fprint(w, string(bytes))
		return
	} else {
		result["success"] = "false"
		result["code"] = "421"
		result["message"] = "该敏感不存在！"
		bytes, _ := json.Marshal(result)
		fmt.Fprint(w, string(bytes))
		return
	}
}

//节点迁移工具把采集节点配置信息定期迁移到redis中
func main() {
	//第一个参数为客户端发起http请求时的接口名，第二个参数是一个func，负责处理这个请求。
	http.HandleFunc("/AddBadword", AddNewBadword)
	http.HandleFunc("/DelBadword", DelBadeword)
	http.HandleFunc("/UpdateAllBadword", UpdateAllBadword)

	//服务器要监听的主机地址和端口号
	//配置注册服务器信息
	id:=tool.RandNum(100)
	ip := tool.GetIP()
	registration := new(consulapi.AgentServiceRegistration)
	registration.ID = "M-Badword:"+strconv.Itoa(id)
	registration.Name = "M-Badword"
	registration.Address = ip
	registration.Port = 8082
	registration.Tags = []string{"Badword 敏感词管理服务！"}
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
