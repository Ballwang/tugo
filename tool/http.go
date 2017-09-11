package tool

import (
	"net/http"
	"net"
	"time"
	"github.com/Ballwang/tugo/soft/softClient"

	"compress/gzip"
	"io"
	"io/ioutil"
)

//获取网站内容
func HttpRequestContent(url string) (int,string) {

	userClient, transport, _ := softClient.NewMcUserAgentClient("user_agent")
	ua, _ := userClient.GetAgentBySiteID(url)

	defer transport.Close()

	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(10 * time.Second)
				c, err := net.DialTimeout(netw, addr, time.Second*10)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
		}, }


	reqest, _ := http.NewRequest("GET", url, nil)
	//tool.ErrorPrint(err)
	reqest.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	reqest.Header.Add("Accept-Encoding", "gzip, deflate")
	reqest.Header.Add("Accept-Language", "zh-cn,zh;q=0.8,en-us;q=0.5,en;q=0.3")
	reqest.Header.Add("Connection", "keep-alive")
	reqest.Header.Add("Referer",url)
	reqest.Header.Add("User-Agent",ua.UserAgent)

	response, _ := client.Do(reqest)
	if response != nil {
		if response.StatusCode==200{
			var body string
			switch response.Header.Get("Content-Encoding") {
			case "gzip":
				reader, _ := gzip.NewReader(response.Body)
				for {
					buf := make([]byte, 1024)
					n, err := reader.Read(buf)
					if err != nil && err != io.EOF {
						break
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
			return response.StatusCode,body
		}else {
			return response.StatusCode,""
		}
	}
	return 404,""
}
