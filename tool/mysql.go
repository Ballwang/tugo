package tool

import (
	"database/sql"
	"github.com/Ballwang/tugo/config"
)

//使用默认配置创建mysql客户端
func NewMysql() (*sql.DB,error) {
	params := config.NewConfig()
	mysqlHost:=params.GetConfig("mysql","mysqlHost")
	mysqlPort:=params.GetConfig("mysql","mysqlPort")
	mysqlUser:=params.GetConfig("mysql","mysqlUser")
	mysqlPassword:=params.GetConfig("mysql","mysqlPassword")
	mysqlCharset:=params.GetConfig("mysql","mysqlCharset")
	mysqlDatabase:=params.GetConfig("mysql","mysqlDatabase")
	db,err:=sql.Open("mysql",mysqlUser+":"+mysqlPassword+"@tcp("+mysqlHost+":"+mysqlPort+")/"+mysqlDatabase+"?charset="+mysqlCharset)
	return db,err
}
