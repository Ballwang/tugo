package main

import (
	"fmt"

	"database/sql"
	"github.com/Ballwang/tugo/config"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"github.com/Ballwang/tugo/tool"
	"context"
)

var stepNumber = 100



func main() {

	params := config.NewConfig()
	mysqlHost := params.GetConfig("mysql", "mysqlHostDev")
	mysqlPort := params.GetConfig("mysql", "mysqlPortDev")
	mysqlUser := params.GetConfig("mysql", "mysqlUserDev")
	mysqlPassword := params.GetConfig("mysql", "mysqlPasswordDev")
	mysqlCharset := params.GetConfig("mysql", "mysqlCharsetDev")
	mysqlDatabase := params.GetConfig("mysql", "mysqlDatabaseDev")
	fmt.Println(mysqlUser + ":" + mysqlPassword + "@tcp(" + mysqlHost + ":" + mysqlPort + ")/" + mysqlDatabase + "?charset=" + mysqlCharset)
	db, err := sql.Open("mysql", mysqlUser+":"+mysqlPassword+"@tcp("+mysqlHost+":"+mysqlPort+")/"+mysqlDatabase+"?charset="+mysqlCharset)
	if err != nil {
		fmt.Println(err)
	}
	n := 0
	isbreak := true
	for {
		rows, err := db.Query("SELECT t1.id,t1.title,t1.keywords,t1.description,t1.url,t1.author,t1.updatetime,t2.content,t2.copyfrom,t2.fbtitle FROM js_news as t1, js_news_data as t2 WHERE t1.id=t2.id Limit " + strconv.Itoa(n) + "," + strconv.Itoa(stepNumber))
		if err != nil {
			fmt.Println(err)

		}
		n = n + stepNumber
		cxt := context.Background()
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
			es := tool.NewESFromConfig()
			jiangsu := params.GetConfig("siteName", "jiangsu")
			md5string := tool.Md5String(jiangsu + id)
			_,err := es.AddData(cxt, c, md5string, "artical", "detial")
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(id)
			isbreak = false
		}
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
	}
	db.Close()

}
