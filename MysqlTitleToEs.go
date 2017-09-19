package main

import (
	"fmt"

	"database/sql"
	"github.com/Ballwang/tugo/config"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"github.com/Ballwang/tugo/tool"
	"context"

	"github.com/garyburd/redigo/redis"

	"gopkg.in/olivere/elastic.v5"
)

var stepNumberTitle = 5000

//数据迁移具有毁灭性质

func main() {
	n := 0
	params := config.NewConfig()
	mysqlHost := params.GetConfig("mysql", "mysqlHostDev")
	mysqlPort := params.GetConfig("mysql", "mysqlPortDev")
	mysqlUser := params.GetConfig("mysql", "mysqlUserDev")
	mysqlPassword := params.GetConfig("mysql", "mysqlPasswordDev")
	mysqlCharset := params.GetConfig("mysql", "mysqlCharsetDev")
	mysqlDatabase := params.GetConfig("mysql", "mysqlDatabaseDev")
	//mysqlHost := params.GetConfig("mysql", "mysqlHost")
	//mysqlPort := params.GetConfig("mysql", "mysqlPort")
	//mysqlUser := params.GetConfig("mysql", "mysqlUser")
	//mysqlPassword := params.GetConfig("mysql", "mysqlPassword")
	//mysqlCharset := params.GetConfig("mysql", "mysqlCharset")
	//mysqlDatabase := params.GetConfig("mysql", "mysqlDatabase")


	db, err := sql.Open("mysql", mysqlUser+":"+mysqlPassword+"@tcp("+mysqlHost+":"+mysqlPort+")/"+mysqlDatabase+"?charset="+mysqlCharset)
	if err != nil {
		fmt.Println(err)
	}
	c,_:=tool.NewRedis()
	defer c.Close()
	r,_:=redis.String(c.Do("GET","title"))
	//r,_:=redis.String(c.Do("GET","artical500-1000"))
	fmt.Println(r)
	if r!=""{
		n,_=strconv.Atoi(r)
	}
	isbreak := true
	es := tool.NewESFromConfig()
	client, err := elastic.NewClient(elastic.SetURL("http://" + es.Host + ":" + es.Port))
	cxt := context.Background()
	jiangsu := params.GetConfig("siteName", "jiangsu")
	for {
		s:=tool.CurrentTimeMillis()
		rows, err := db.Query("SELECT id,title,url FROM js_news where id >="+strconv.Itoa(n)+" AND id<="+strconv.Itoa(n+stepNumberTitle))
		if err != nil {
			fmt.Println(err)
		}


		for rows.Next() {
			var id, title, url  string
			if err := rows.Scan(&id, &title,  &url); err != nil {
				fmt.Println(err)
			}
			c := &tool.Title{
				Title:       title,
				Url:         url,
				ContentID:   id,
			}


			md5string := tool.Md5String(jiangsu + id)
			_,err := es.AddDataWithTitleClient(cxt, client,c, md5string, "news", "title")
			if err != nil {
				fmt.Println(err)
			}

			isbreak = false
		}
		n = n + stepNumberTitle
		c.Do("SET","title",n)
		//c.Do("SET","artical500-1000",n)

		if db != nil {

		}
		if err != nil {

		}
		rows.Close()
		if isbreak {
			break
		}
		isbreak = true
		e:=tool.CurrentTimeMillis()
		tool.ShowTime(s,e)

	}
	db.Close()

}
