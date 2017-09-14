package main

import (

	"fmt"

	"github.com/Ballwang/tugo/tool"
)

func main()  {

	//c,e:=tool.NewRedisCluster()
	//if e!=nil{
	//	fmt.Println(e)
	//	fmt.Println("=====")
	//}
	//
	//
	//reply, err := redis.StringMap(c.Do("HGETALL", "MonitorMiddleSiteHash"))
	//
	//
	//if err != nil {
	//
	//} else {
	//	for k, v := range reply {
	//		fmt.Println(k)
	//		fmt.Println("===========")
	//		fmt.Println(v)
	//		break
	//	}
	//}
	v:="http://news.enorth.com.cn/tj/wytd/|||http://news.enorth.com.cn/system/more/17001009000000000/0028/17001009000000000_00002844.shtml|||http://news.enorth.com.cn/system/more/17001009000000000/0028/17001009000000000_00002843.shtml"

	fmt.Println(tool.GetSplitOne(v,"---"))

	//for i:=0;i<8000 ;i++  {
	//	k:="M:-http://finance.jxgdw.com/news/index.html"+
	//"http://finance.jxgdw.com/news/index1.html"+
	//"http://finance.jxgdw.com/news/index2.html"+
	//"http://finance.jxgdw.com/news/index3.html"+
	//"http://finance.jxgdw.com/news/index4.html"
	//	r,e:=c.Do("HSET","MonitorMiddleSiteHash11111",k+strconv.Itoa(i),k)
	//	if e!=nil{
	//		fmt.Println(e)
	//		fmt.Println("-----")
	//	}
	//	fmt.Println(r)
	//}








	
}
