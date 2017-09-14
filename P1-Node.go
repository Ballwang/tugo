package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/Ballwang/tugo/config"
	"fmt"
	"github.com/Ballwang/tugo/tool"
	"net/http"
	"strconv"

)

//节点迁移工具把采集节点配置信息定期迁移到redis中
func AddNodeToMonitor(w http.ResponseWriter, req *http.Request) {
	db, err := tool.NewMysql()
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	rows, err := db.Query("SELECT nodeid, name,sourcecharset,urlpage,url_start,url_end FROM js_collection_node")
	defer db.Close()
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	c,err:=tool.NewRedisCluster()
	tool.ErrorPrint(err)
	defer c.Close()
	for rows != nil && rows.Next() {
		var nodeid, name, sourcecharset, urlpage, url_start, url_end string
		if err := rows.Scan(&nodeid, &name, &sourcecharset, &urlpage, &url_start, &url_end);err != nil {
			fmt.Fprint(w, err)
			return
		}
		urlpage=tool.TrimReplace(urlpage,config.Separate)
		if urlpage != ""&&nodeid!="" {
			c.Do("HSET",config.MonitorHash,"Host:-"+urlpage,urlpage)
			c.Do("HSET",config.MonitorHash,"Time:-"+urlpage,10)
			c.Do("HSET",config.MonitorHash,"nodeid:-"+urlpage,nodeid)
			c.Do("HSET",config.MonitorHash,"sourcecharset:-"+urlpage,sourcecharset)
			c.Do("HSET",config.MonitorHash,"url_start:-"+urlpage,url_start)
			c.Do("HSET",config.MonitorHash,"url_end:-"+urlpage,url_end)
			c.Do("HSET", config.MonitorMiddleSiteHash, "M:-"+urlpage, urlpage)
		}
	}
	fmt.Fprint(w, true)
	return
}

//服务状态监控接口
func NodeToMonitorState(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, true)
}

func main() {
	var serverID = "P1-MysqlNodeToRedis"
	var serverPort = 8087
	ip := tool.GetIP()

	http.HandleFunc("/AddNodeToMonitor", AddNodeToMonitor)
	http.HandleFunc("/State", NodeToMonitorState)
	register := &tool.ConsulRegister{Id: serverID, Name: "P1-数据库节点同步监控服务", Port: serverPort, Tags: []string{"数据库节点同步监控服务"}}
	register.RegisterConsulService()
	err := http.ListenAndServe(ip+":"+strconv.Itoa(serverPort), nil)

	if err != nil {
		fmt.Println("Listen And Serve error: ", err.Error())
	}

}
