package userAgentService

import (
	"math/rand"
	"time"
	"github.com/Ballwang/mcserver/gen-go/UserAgent"
)

type WebAgent struct {
}

var WebUserAgent = []string{
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/536.11 (KHTML, like Gecko) Chrome/20.0.1132.11 TaoBrowser/2.0 Safari/536.11",
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.1 (KHTML, like Gecko) Chrome/21.0.1180.71 Safari/537.1 LBBROWSER",
	"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; WOW64; Trident/5.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; .NET4.0C; .NET4.0E; LBBROWSER)",
	"Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1; QQDownload 732; .NET4.0C; .NET4.0E; LBBROWSER)",
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.84 Safari/535.11 LBBROWSER",
	"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.1; WOW64; Trident/5.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; .NET4.0C; .NET4.0E)",
	"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; WOW64; Trident/5.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; .NET4.0C; .NET4.0E; QQBrowser/7.0.3698.400)",
	"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.1; WOW64; Trident/5.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; .NET4.0C; .NET4.0E)",
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.1 (KHTML, like Gecko) Chrome/21.0.1180.89 Safari/537.1",
	"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.1; WOW64; Trident/5.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; Media Center PC 6.0; .NET4.0C; .NET4.0E)",
	"Mozilla/5.0 (Windows NT 5.1) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.84 Safari/535.11 SE 2.X MetaSr 1.0",
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:2.0b13pre) Gecko/20110307 Firefox/4.0b13pre",
	"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11",
	"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Win64; x64; Trident/5.0)",
	"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; WOW64; Trident/5.0)",
	"Mozilla/5.0 (X11; U; Linux x86_64; zh-CN; rv:1.9.2.10) Gecko/20100922 Ubuntu/10.10 (maverick) Firefox/3.6.10",
	"Mozilla/5.0(Macintosh;U;IntelMacOSX10_6_8;en-us)AppleWebKit/534.50(KHTML,likeGecko)Version/5.1Safari/534.50",
	"Mozilla/5.0(Windows;U;WindowsNT6.1;en-us)AppleWebKit/534.50(KHTML,likeGecko)Version/5.1Safari/534.50",
	"Mozilla/5.0(compatible;MSIE9.0;WindowsNT6.1;Trident/5.0;",
	"Mozilla/4.0(compatible;MSIE8.0;WindowsNT6.0;Trident/4.0)",
	"Mozilla/4.0(compatible;MSIE7.0;WindowsNT6.0)",
	"Mozilla/4.0(compatible;MSIE6.0;WindowsNT5.1)",
	"Mozilla/5.0(Macintosh;IntelMacOSX10.6;rv:2.0.1)Gecko/20100101Firefox/4.0.1",
	"Mozilla/5.0(WindowsNT6.1;rv:2.0.1)Gecko/20100101Firefox/4.0.1",
	"Opera/9.80(Macintosh;IntelMacOSX10.6.8;U;en)Presto/2.8.131Version/11.11",
	"Opera/9.80(WindowsNT6.1;U;en)Presto/2.8.131Version/11.11",
	"Mozilla/5.0(Macintosh;IntelMacOSX10_7_0)AppleWebKit/535.11(KHTML,likeGecko)Chrome/17.0.963.56Safari/535.11",
	"Mozilla/4.0(compatible;MSIE7.0;WindowsNT5.1)",
	"Mozilla/4.0(compatible;MSIE7.0;WindowsNT5.1;360SE)",
	"Mozilla/5.0 (compatible; Baiduspider/2.0; +http://www.baidu.com/search/spider.html)",
	"Mozilla/5.0 (compatible; Baiduspider/2.0; +http://www.baidu.com/search/spider.html)",
	"Mozilla/5.0 (compatible; Baiduspider/2.0; +http://www.baidu.com/search/spider.html)",
	"Mozilla/5.0 (compatible; Baiduspider/2.0; +http://www.baidu.com/search/spider.html)",
	"Mozilla/5.0 (compatible; Baiduspider/2.0; +http://www.baidu.com/search/spider.html)",
	"Mozilla/5.0 (compatible; Baiduspider/2.0; +http://www.baidu.com/search/spider.html)",
	"Sogou web spider/4.0(+http://www.sogou.com/docs/help/webmasters.htm#07)",
	"Mozilla/5.0 (compatible; Yahoo! Slurp/3.0; http://help.yahoo.com/help/us/ysearch/slurp)",
	"Sosospider+(+http://help.soso.com/webspider.htm)",
	"Mozilla/5.0 (compatible; YoudaoBot/1.0; http://www.youdao.com/help/webmaster/spider/; )",
	"Mozilla/5.0 (compatible; Baiduspider/2.0; +http://www.baidu.com/search/spider.html)",
	"Sogou web spider/4.0(+http://www.sogou.com/docs/help/webmasters.htm#07)",
	"Mozilla/5.0 (compatible; Yahoo! Slurp/3.0; http://help.yahoo.com/help/us/ysearch/slurp)",
	"Sosospider+(+http://help.soso.com/webspider.htm)",
	"Mozilla/5.0 (compatible; YoudaoBot/1.0; http://www.youdao.com/help/webmaster/spider/; )",
	"Mozilla/5.0 (compatible; Baiduspider/2.0; +http://www.baidu.com/search/spider.html) ",
	"Sogou web spider/4.0(+http://www.sogou.com/docs/help/webmasters.htm#07)",
	"Mozilla/5.0 (compatible; Yahoo! Slurp/3.0; http://help.yahoo.com/help/us/ysearch/slurp)",
	"Sosospider+(+http://help.soso.com/webspider.htm)",
	"Mozilla/5.0 (compatible; YoudaoBot/1.0; http://www.youdao.com/help/webmaster/spider/; )",
	"Mozilla/5.0 (compatible; Baiduspider/2.0; +http://www.baidu.com/search/spider.html)",
	"Sogou web spider/4.0(+http://www.sogou.com/docs/help/webmasters.htm#07)",
	"Mozilla/5.0 (compatible; Yahoo! Slurp/3.0; http://help.yahoo.com/help/us/ysearch/slurp)",
	"Sosospider+(+http://help.soso.com/webspider.htm)",
	"Mozilla/5.0 (compatible; YoudaoBot/1.0; http://www.youdao.com/help/webmaster/spider/; )",
}

//这里的IP需要独立出去
var IP  =[]string{
	"192.168.3.50",
	"192.168.3.51",
	"192.168.3.52",
	"192.168.3.53",
	"192.168.3.54",
}



//随机返回UserAgent 头部信息
func (A *WebAgent) GetUserAgent() (string, error) {
	i:=rand.Int63n(1000)
	rand.Seed(time.Now().Unix()+i)
	return WebUserAgent[rand.Intn(len(WebUserAgent))], nil
}

//根据网站名称返回网站头部
func (A *WebAgent) GetUserAgentById(site string) (string, error) {
	i:=rand.Int63n(1000)
	rand.Seed(time.Now().Unix()+i)
	return WebUserAgent[rand.Intn(len(WebUserAgent))], nil
}

//根据网站名称返回网站带来IP
func (A *WebAgent)GetAgentIp(site string)(string,error)  {
	i:=rand.Int63n(1000)
	rand.Seed(time.Now().Unix()+i)
	return IP[rand.Intn(len(IP))], nil
}

//根据网站名称返回网站信息
func (A *WebAgent)GetAgentBySiteID(site string)(*UserAgent.Agent,error) {
	userString,_:=A.GetUserAgentById(site)
	ipString,_:=A.GetAgentIp(site)
	siteAgent:=&UserAgent.Agent{UserAgent:userString,AgentIp:ipString,SiteId:""}
	return siteAgent,nil
}
