package main

import (
	"fmt"
	"github.com/Ballwang/tugo/tool"
	"github.com/Ballwang/tugo/config"
	"strings"
	"strconv"
	"net/http"
	"net"
	"time"
	"runtime"
)

var MaxPage int = 2

var MaxProcess  = 1000

var completeCategory chan int = make(chan int)

//需要定期执行

func main() {
	s := tool.CurrentTimeMillis()
	fmt.Println("筛选分类链接是否可以访问...")
	mapString := tool.RedisHGETALL(config.MonitorMiddleSiteHash)
	urlArray := []string{}
	c,_:=tool.NewRedis()
	defer c.Close()
	for _, v := range mapString {
		if strings.Contains(v, "(*)") {
			for i := 1; i <= MaxPage; i++ {
				UrlString:=string(strings.Replace(v, "(*)", strconv.Itoa(i), -1))
				urlArray = append(urlArray,UrlString )
				c.Do("HSET",config.MonitorShowHash,UrlString,v)
			}
		} else {
			//fmt.Println(strings.Split(v, "\n"))
			url := strings.Fields(v)
			if len(url) > 1 {
				i := 0
				for _, v1 := range url {
					if i <=2 { 
						urlArray = append(urlArray, v1)
						c.Do("HSET",config.MonitorShowHash,v1,v)
					} else {
						break
					}
					i++
				}
			} else {
				urlArray = append(urlArray, v)
			}
		}
	}
	process := len(urlArray)
	runtime.GOMAXPROCS(6)
	j:=1
	if process > 0 {
		for _, v := range urlArray {
			go GetCatgoryListUrl(v)
			if j%MaxProcess==0{
				for i := 0; i < MaxProcess; i++ {
					<-completeCategory
				}
			}
			j++
		}
	}
	fmt.Println(len(urlArray))

	e := tool.CurrentTimeMillis()
	fmt.Printf("本次调用用时:%d-%d=%d毫秒\n", e, s, (e - s))
}

//判断列表页链接是否可以正常获取
func GetCatgoryListUrl(siteUrl string) {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(15 * time.Second)
				c, err := net.DialTimeout(netw, addr, time.Second*20)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
		}, }
	reqest, err := http.NewRequest("GET", siteUrl, nil)
	tool.ErrorPrint(err)
	reqest.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	reqest.Header.Add("Accept-Encoding", "gzip, deflate")
	reqest.Header.Add("Accept-Language", "zh-cn,zh;q=0.8,en-us;q=0.5,en;q=0.3")
	reqest.Header.Add("Connection", "keep-alive")
	reqest.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11")
	response, _ := client.Do(reqest)
	if response != nil {
		c,_:=tool.NewRedis()
		if response.StatusCode == 200 {
			if c!=nil{
				c.Do("HSET",config.MonitorSiteHash,"S-:"+siteUrl,siteUrl)
			}
		} else {
			c.Do("SADD",config.BadSite,siteUrl)
		}
	}
	completeCategory <- 1
}
