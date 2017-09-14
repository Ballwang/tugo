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

var stepNumber = 2000



func main() {
	n := 0
	params := config.NewConfig()
	//mysqlHost := params.GetConfig("mysql", "mysqlHostDev")
	//mysqlPort := params.GetConfig("mysql", "mysqlPortDev")
	//mysqlUser := params.GetConfig("mysql", "mysqlUserDev")
	//mysqlPassword := params.GetConfig("mysql", "mysqlPasswordDev")
	//mysqlCharset := params.GetConfig("mysql", "mysqlCharsetDev")
	//mysqlDatabase := params.GetConfig("mysql", "mysqlDatabaseDev")
	mysqlHost := params.GetConfig("mysql", "mysqlHost")
	mysqlPort := params.GetConfig("mysql", "mysqlPort")
	mysqlUser := params.GetConfig("mysql", "mysqlUser")
	mysqlPassword := params.GetConfig("mysql", "mysqlPassword")
	mysqlCharset := params.GetConfig("mysql", "mysqlCharset")
	mysqlDatabase := params.GetConfig("mysql", "mysqlDatabase")


	db, err := sql.Open("mysql", mysqlUser+":"+mysqlPassword+"@tcp("+mysqlHost+":"+mysqlPort+")/"+mysqlDatabase+"?charset="+mysqlCharset)
	if err != nil {
		fmt.Println(err)
	}
	c,_:=tool.NewRedis()
	defer c.Close()
	r,_:=redis.String(c.Do("GET","artical1-500"))
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
		rows, err := db.Query("SELECT t1.id,t1.title,t1.keywords,t1.description,t1.url,t1.author,t1.updatetime,t2.content,t2.copyfrom,t2.fbtitle FROM js_news as t1, js_news_data as t2 WHERE t1.id=t2.id Limit " + strconv.Itoa(n) + "," + strconv.Itoa(stepNumber))
		if err != nil {
			fmt.Println(err)
		}


		for rows.Next() {
			var id, title, keywords, description, url, author, content, copyfrom, fbtitles, updatetime string
			if err := rows.Scan(&id, &title, &keywords, &description, &url, &author, &updatetime, &content, &copyfrom, &fbtitles); err != nil {
				fmt.Println(err)
			}
			c := &tool.Content{
				Title:       title,
				TitleSub:    fbtitles,
				Keywords:    keywords,
				Description: description,
				Author:      author,
				Time:        updatetime,
				CopyFrom:    copyfrom,
				Content:     content,
				Url:         url,
				ContentID:   id,
			}


			md5string := tool.Md5String(jiangsu + id)
			_,err := es.AddDataWithClient(cxt, client,c, md5string, "artical", "detial")
			if err != nil {
				fmt.Println(err)
			}

			isbreak = false
		}
		n = n + stepNumber
		c.Do("SET","artical1-500",n)
		//c.Do("SET","artical500-1000",n)
		fmt.Println("===========================================")
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
