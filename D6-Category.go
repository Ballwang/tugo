package main

import (
	"github.com/Ballwang/tugo/tool"
	"github.com/Ballwang/tugo/config"
	"time"
	"github.com/axgle/mahonia"
	"github.com/chasex/redis-go-cluster"
	"strings"
	"net/http"
	"strconv"
	"fmt"
)

var completeCat chan int = make(chan int)
var textLen = 24 //8个汉字
var serverIDCategory = "D6-Category"
var serverCategoryPort = 8085

//过滤列表页链接
func FilterCategoryList(w http.ResponseWriter, req *http.Request) {

	for {
		c, _ := tool.NewRedisCluster()
		IsEnd := 0
		for {
			SiteList := []string{}
			for i := 0; i < config.MaxProcess; i++ {
				//for i := 0; i < 1; i++ {
				reply, _ := c.Do("LPOP", config.UpdateList)
				if reply != nil {
					SiteList = append(SiteList, string(reply.([]byte)))
				} else {
					IsEnd = 1
					//退出当前循环
					break
				}
			}

			process := len(SiteList)
			if process > 0 {
				for _, v := range SiteList {
					//进入并发程序
					go GetUrlFromCategory(c, v)
				}
				for i := 0; i < process; i++ {
					<-completeCat
				}
			}

			if IsEnd == 1 {
				break
			}
			tool.SetServerState(serverIDCategory, "5")
			//break
		}
		//break
		tool.SetServerState(serverIDCategory, "5")

		c.Close()
		time.Sleep(1 * time.Second)
	}
}

//返回服务状态
func ServerCategoryState(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, tool.GetServerState(serverIDCategory))
}

//列表页详细链接提取
func main() {

	ip := tool.GetIP()
	var serverIDCategory = "D6-Category:"+ip
	config:=config.NewConfig()
	serverCategoryPort,_:=strconv.Atoi(config.GetConfig("D6-Category","port"))
	http.HandleFunc("/D6-Category", FilterCategoryList)
	http.HandleFunc("/State", ServerCategoryState)
	register := &tool.ConsulRegister{Id: serverIDCategory, Name: "D6-Category", Port: serverCategoryPort, Tags: []string{"列表页链接提取服务！自动识别提取列表页链接！"}}
	register.RegisterConsulService()
	err := http.ListenAndServe(ip+":"+strconv.Itoa(serverCategoryPort), nil)

	if err != nil {
		fmt.Println("Listen And Serve error: ", err.Error())
	}
}


//从栏目页提取链接
func GetUrlFromCategory(c *redis.Cluster, url string) {

	content, _ := c.Do("GET", config.PrefixCategory+url)

	//fmt.Println(config.PrefixCategory+url)
	if content != nil {
		//获取链接提取规则
		urlString, _ := c.Do("HGET", config.MonitorShowHash, url)

		//转化为小写

		var contentString string
		var urlParent string

		if urlString != nil {
			urlParent, _ = redis.String(urlString, nil)
		} else {
			urlParent = url
		}

		//urlStart, _ := redis.String(c.Do("HGET", config.MonitorHash, config.UrlStart+url))
		//urlEnd, _ := redis.String(c.Do("HGET", config.MonitorHash, config.UrlEnd+url))
		//获取链接配置信息
		Sourcecharset, _ := redis.String(c.Do("HGET", config.MonitorHash, config.Sourcecharset+urlParent))

		if !strings.Contains(Sourcecharset, "utf") {
			dec := mahonia.NewDecoder("gbk")
			contentString = dec.ConvertString(string(content.([]byte)))
		} else {
			contentString, _ = redis.String(content, nil)
		}
		//开始链接提取
		contentString = strings.ToLower(contentString)

		resultMap := map[string]string{}

		result := tool.ListHref(contentString)

		//fmt.Print(contentString)

		//分析网站链接是否为文章链接
		if len(result) > 0 {
			for _, v := range result {
				if v != "" {
					aText := tool.GetALabelText(v)
					if len(aText) > textLen {
						resultMap[v] = aText
					}
				} else {

				}
			}
		}

		//os.Exit(-1)

		//分析链接特征
		count := map[int]int{}
		countMap := map[string]int{}
		countSubfix := map[string]int{}
		if len(resultMap) > 0 {
			for k, _ := range resultMap {
				url := tool.GetALabelUrl(k)
				subfix := tool.GetHostUrlSuffix(url)
				urlNoHost := tool.RemovHostNameByUrl(url)
				urlArray := strings.Split(urlNoHost, "/")
				if len(urlArray) > 0 {
					count[len(urlArray)]++
					countMap[url] = len(urlArray)
				}
				if subfix != "" {
					subArray := strings.Split(url, subfix)
					if len(subArray) > 0 {
						countSubfix[subfix]++
					}
				} else {
					countSubfix["nilSubfix"]++
				}
			}

			//筛选权重最高的链接形式
			maxLen := 0
			secondLen := 0
			maxKey := 0
			secondKey := 0
			if len(count) > 0 {
				for k, v := range count {
					if v > maxLen {
						secondLen = maxLen
						secondKey = maxKey
						maxLen = v
						maxKey = k
					}
					if v > secondLen && v < maxLen {
						secondLen = v
						secondKey = k
					}
				}
			}
			//筛选使用最多的后缀
			maxLenSubfix := 0
			maxkeySubfix := ""
			if len(countSubfix) > 0 {
				for key, v := range countSubfix {
					if v > maxLenSubfix {
						maxLenSubfix = v
						maxkeySubfix = key
					}

				}
			}

			result := []string{}
			if len(countMap) > 0 && maxLen >= config.MinNumOfUrl {
				for k, v := range countMap {
					//兼容网页侧边栏相似链接现象

					isDo := false
					if (maxLen - secondLen) <= config.MaxMinusValue {
						if v == maxKey || v == secondKey {
							isDo = true
						}
					} else {
						if v == maxKey {
							isDo = true
						}
					}
					if isDo {
						if maxkeySubfix != "" {
							if strings.Contains(k, maxkeySubfix) {
								//result = append(result, k)
								AbsoluteUrl := tool.GetAbsoluteUrl(url, k)
								result = append(result, AbsoluteUrl)
								historyKey := tool.GetHostUri(url)
								md5String := tool.Md5String(AbsoluteUrl)
								key := config.HistoryPrefix + historyKey
								r, _ := c.Do("SISMEMBER", key, md5String)
								var j int64 = 1
								if r != j {
									c.Do("SADD", key, md5String)
									c.Do("HSET", config.ContentParentHash, AbsoluteUrl, urlParent)
									c.Do("SADD", config.ContentUrlSet, AbsoluteUrl)
									SetUpdateTime(c, urlParent)
								}
								if historyKey == "" {
									c.Do("SADD", config.NullKey, url)
								}
							}
						} else {
							c.Do("SADD", config.NullUrlInCategory, urlParent)
						}
					}

				}
			}
			//fmt.Println(result)

		}
	}
	completeCat <- 1
}

//记录采集条数历史
func SetUpdateTime(c *redis.Cluster, url string) {
	t := time.Now()
	y := t.Year()
	m := t.Month().String()
	d := t.Day()
	h := t.Hour()
	s := "count:" + strconv.Itoa(y) + "-" + m + "-" + strconv.Itoa(d) + "-" + strconv.Itoa(h)
	c.Do("HINCRBY", s, url, 1)
}
