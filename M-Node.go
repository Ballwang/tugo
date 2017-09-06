package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/Ballwang/tugo/config"
	"database/sql"
	"fmt"
	"github.com/Ballwang/tugo/tool"

)


//节点迁移工具把采集节点配置信息定期迁移到redis中
func main()  {
	params := config.NewConfig()
	mysqlHost:=params.GetConfig("mysql","mysqlHost")
	mysqlPort:=params.GetConfig("mysql","mysqlPort")
	mysqlUser:=params.GetConfig("mysql","mysqlUser")
	mysqlPassword:=params.GetConfig("mysql","mysqlPassword")
	mysqlCharset:=params.GetConfig("mysql","mysqlCharset")
	mysqlDatabase:=params.GetConfig("mysql","mysqlDatabase")
	fmt.Println(mysqlUser+":"+mysqlPassword+"@tcp("+mysqlHost+":"+mysqlPort+")/"+mysqlDatabase+"?charset="+mysqlCharset)
	db,err:=sql.Open("mysql",mysqlUser+":"+mysqlPassword+"@tcp("+mysqlHost+":"+mysqlPort+")/"+mysqlDatabase+"?charset="+mysqlCharset)
	rows,err:=db.Query("SELECT nodeid, name,sourcecharset,urlpage,url_start,url_end FROM js_collection_node")
	if err!=nil{

		fmt.Println(err)
	}

	c,err:=tool.NewRedis()
	tool.ErrorPrint(err)
	defer c.Close()

	for rows.Next(){
		var nodeid, name,sourcecharset,urlpage,url_start,url_end  string

		if err := rows.Scan(&nodeid, &name,&sourcecharset,&urlpage,&url_start,&url_end); err != nil {
			fmt.Println(err)
		}
		fmt.Println(nodeid)
		if urlpage !="" {
			c.Do("HSET",config.MonitorHash,"Host:-"+urlpage,urlpage)
			c.Do("HSET",config.MonitorHash,"Time:-"+urlpage,10)
			c.Do("HSET",config.MonitorHash,"nodeid:-"+urlpage,nodeid)
			c.Do("HSET",config.MonitorHash,"sourcecharset:-"+urlpage,sourcecharset)
			c.Do("HSET",config.MonitorHash,"url_start:-"+urlpage,url_start)
			c.Do("HSET",config.MonitorHash,"url_end:-"+urlpage,url_end)
			c.Do("HSET",config.MonitorMiddleSiteHash,"M:-"+urlpage,urlpage)
		}
	}


	if db !=nil{

	}
	if err!=nil{

	}
}
