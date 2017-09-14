package main

import (
	"github.com/Ballwang/tugo/soft/softClient"
	"fmt"
	"time"
	"net/http"
	"os"
	"github.com/Ballwang/tugo/tool"
	"io/ioutil"
	"compress/gzip"
	"io"
	"crypto/md5"
	"encoding/hex"
	"github.com/chasex/redis-go-cluster"
	"github.com/Ballwang/tugo/config"
	"strconv"

)

var complete chan int = make(chan int)

//定期获取监控队列进行扫描监控
func MonitorCategory(w http.ResponseWriter, req *http.Request) {

	for {
		c, _ := tool.NewRedisCluster()
		IsEnd := 0
		client, transport, _ := softClient.NewMcUserAgentClient("user_agent")
		for {
			SiteList := []string{}
			for i := 0; i < config.MaxListProcess; i++ {
				reply, err := c.Do("LPOP", config.MonitorList)
				tool.Error(err)
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
					//获取不同的User-agent 列表
					ua, _ := client.GetAgentBySiteID("baidu.com")
					//进入并发程序
					go HtmlEye(c,v, ua.UserAgent, ua.AgentIp)
				}
				transport.Close()
				for i := 0; i < process; i++ {
					<-complete
				}
			}

			if IsEnd == 1 {
				break
			}
			tool.SetServerState("S4-MonitorCategory","5")
		}
		c.Close()
		tool.SetServerState("S4-MonitorCategory","5")
		time.Sleep(1 * time.Second)
	}

}

func MonitorCategoryState(w http.ResponseWriter, req *http.Request)  {
	fmt.Fprint(w,tool.GetServerState("S4-MonitorCategory"))
}

func main() {

	var serverID = "S4-MonitorCategory"
	var serverPort = 8090
	ip := tool.GetIP()

	http.HandleFunc("/MonitorCategory", MonitorCategory)
	http.HandleFunc("/State", MonitorCategoryState)
	register := &tool.ConsulRegister{Id: serverID, Name: "列表更新监控服务", Port: serverPort, Tags: []string{"列表更新监控服务，监控目标采集点是否更新内容！"}}
	register.RegisterConsulService()
	err := http.ListenAndServe(ip+":"+strconv.Itoa(serverPort), nil)

	if err != nil {
		fmt.Println("Listen And Serve error: ", err.Error())
	}


}

//监控网站是否有更新，使用缓冲区进行Md5加密，存入redis 进行新旧值进行比较
func HtmlEye(c *redis.Cluster,url string, ua string, ip string) {

	client := &http.Client{}
	params := config.NewMainParams()
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Fatal error ", err.Error())

	}
	request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	request.Header.Add("Accept-Encoding", "gzip, deflate")
	request.Header.Add("Accept-Language", "zh-cn,zh;q=0.8,en-us;q=0.5,en;q=0.3")
	request.Header.Add("Connection", "keep-alive")
	request.Header.Add("Referer", url)
	request.Header.Add("User-Agent", ua)
	response, _ := client.Do(request)

	//创建 redis，在这里创建是因为go 分布式运行避免创建错误

	//判断网页是否更新
	if response != nil {
		//请求成功
		if response.StatusCode == 200 {

			if err != nil {
				fmt.Println("Redis Fatal error ", err.Error())
			}
			var body string
			switch response.Header.Get("Content-Encoding") {
			case "gzip":
				reader, _ := gzip.NewReader(response.Body)
				for {
					buf := make([]byte, 1024)
					n, err := reader.Read(buf)
					if err != nil && err != io.EOF {
						panic(err)
					}
					if n == 0 {
						break
					}
					body += string(buf)
				}
			default:
				bodyByte, _ := ioutil.ReadAll(response.Body)
				body = string(bodyByte)
			}

			//对返回值进行 md5 加密判断是否更新
			h := md5.New()
			h.Write([]byte(body))
			md5String := hex.EncodeToString(h.Sum(nil))

			//获取历史记录判断是否有存在历史
			reply, err := c.Do("GET", url)
			//判断是否查找到历史记录
			if reply != nil {
				//获取历史md5记录
				replyMd5String, _ := redis.String(reply, err)
				//判断新记录是否与历史记录相同
				if md5String != replyMd5String {
					//新值和旧值不相同，写入更新队列，排除并发带来的重复值
					c.Do("SADD", params.UpdateListSet, url)
					c.Do("SET", config.PrefixCategory+url, []byte(body))
					//c.Do("SADD",config.ValueSet,url)
					//设置新值
					c.Do("SET", url, md5String)
				}
			} else {
				//未查找到历史记录，则记录首次历史
				c.Do("SET", url, md5String)
				c.Do("SET", config.PrefixCategory+url, []byte(body))
				//c.Do("SADD",config.ValueSet,url)
				//加入更新队列，排除并发带来的重复值
				c.Do("SADD", params.UpdateListSet, url)
			}
		} else if response.StatusCode >= 400 {
			_, err = c.Do("SADD", params.BadSite, url)
			if err != nil {
				fmt.Println("Redis Fatal error ", err.Error())
			}
		}
		defer response.Body.Close()
	} else {
		_, err = c.Do("SADD", params.BadSite, url)
		if err != nil {
			fmt.Println("Redis Fatal error ", err.Error())
		}
	}

	complete <- 1

	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(0)
	}
}
