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
	"github.com/chasex/redis-go-cluster"
	"encoding/json"

)

var MaxPage int = 1

var MaxProcess = 1000

var completeCategory chan int = make(chan int)

//检查监控链接是否能被正常访问
func FindBadMonitorList(w http.ResponseWriter, req *http.Request) {

	s := tool.CurrentTimeMillis()
	fmt.Println("筛选分类链接是否可以访问...")
	mapString := tool.RedisClusterHGETALL(config.MonitorMiddleSiteHash)

	urlArray := []string{}
	c, _ := tool.NewRedisCluster()
	defer c.Close()
	//分析链接
	for _, v := range mapString {
		if strings.Contains(v, "(*)") {
			for i := 1; i <= MaxPage; i++ {
				UrlString := string(strings.Replace(v, "(*)", strconv.Itoa(i), -1))
				urlArray = append(urlArray, UrlString)
				c.Do("HSET", config.MonitorShowHash, UrlString, v)
			}
		} else {
			//fmt.Println(strings.Split(v, "\n"))
			url := tool.GetSplitOne(v,config.Separate)
			if len(url) > 1 {
				i := 1
				for _, v1 := range url {
					if i <= MaxPage {
						urlArray = append(urlArray, v1)
						c.Do("HSET", config.MonitorShowHash, v1, v)
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

	//开始请求
	process := len(urlArray)
	runtime.GOMAXPROCS(config.MaxCpu)
	j := 0
	if process > 0 {
		for _, v := range urlArray {
			go GetCatgoryListUrl(v,c)
			j++
			if j%MaxProcess == 0 {
				for i := 0; i < MaxProcess; i++ {
					<-completeCategory
				}
				j=0
			}
		}
		for i := 0; i < j; i++ {
			<-completeCategory
		}

	}
	e := tool.CurrentTimeMillis()
	fmt.Printf("本次调用用时:%d-%d=%d毫秒\n", e, s, (e - s))

}

//获取不能访问监控地址
func FindBadMonitorSite(w http.ResponseWriter, req *http.Request) {
	c, _ := tool.NewRedisCluster()
	r, _ := c.Do("SMEMBERS", config.BadSite)
	list := []string{}
	if r != nil {
		result, _ := redis.Values(r, nil)
		for _, v := range result {
			list = append(list, string(v.([]byte)))
		}
	}

	b, _ := json.Marshal(list)
	fmt.Fprint(w, string(b))
	return
}

//状态反馈
func BadWordMonitorState(w http.ResponseWriter, req *http.Request)  {
	fmt.Fprint(w, true)
}

func main() {

	ip := tool.GetIP()
	var serverID = "P2-FindBadMonitorList:"+ip
	config:=config.NewConfig()
	serverPort,_:=strconv.Atoi(config.GetConfig("P2-FindBadMonitorList","port"))


	http.HandleFunc("/StartMonitorList", FindBadMonitorList)
	http.HandleFunc("/FindBdaWordList", FindBadMonitorSite)
	http.HandleFunc("/State", BadWordMonitorState)
	register := &tool.ConsulRegister{Id: serverID, Name: "P2-FindBadMonitorList", Port: serverPort, Tags: []string{"P2-链接失效监控服务"}}
	register.RegisterConsulService()
	err := http.ListenAndServe(ip+":"+strconv.Itoa(serverPort), nil)

	if err != nil {
		fmt.Println("Listen And Serve error: ", err.Error())
	}
}

//判断列表页链接是否可以正常获取,加入redis 客户端防止redis重复创建
func GetCatgoryListUrl(siteUrl string,c *redis.Cluster) {
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
	s:=tool.CurrentTimeMillis()
	reqest, _ := http.NewRequest("GET", siteUrl, nil)
	//tool.ErrorPrint(err)
	reqest.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	reqest.Header.Add("Accept-Encoding", "gzip, deflate")
	reqest.Header.Add("Accept-Language", "zh-cn,zh;q=0.8,en-us;q=0.5,en;q=0.3")
	reqest.Header.Add("Connection", "keep-alive")
	reqest.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11")
	response, _ := client.Do(reqest)
	e:=tool.CurrentTimeMillis()
	fmt.Println(siteUrl)
	tool.ShowTime(s,e)
	fmt.Println("--------------------------\n")
	if response != nil {

		if response.StatusCode == 200 {
			if c != nil {
				c.Do("HSET", config.MonitorSiteHash, "S-:"+siteUrl, siteUrl)
			}
		} else {
			c.Do("SADD", config.BadSite, siteUrl)
		}
	}

	completeCategory <- 1
}
