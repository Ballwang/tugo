package main

import (
	"github.com/Ballwang/tugo/tool"
	"strconv"

	"strings"

	"fmt"

	"github.com/Ballwang/tugo/config"
)

func ExampleScrape() {
	//var r  []int
	s := tool.CurrentTimeMillis()
	content := make(map[int]map[int]string)
	c := make(map[int]map[int]string)
	url := []string{
		"http://hn.rednet.cn/c/2017/08/24/4402951.htm",
		"http://hn.rednet.cn/c/2017/08/24/4403054.htm",
		"http://hn.rednet.cn/c/2017/08/24/4403056.htm",
		"http://hn.rednet.cn/c/2017/08/24/4402805.htm",
		"http://hn.rednet.cn/c/2017/08/24/4402660.htm",
	}
	for k, v := range url {
		content[k], c[k] = tool.GetHtmlDom("div", v, "gbk")
	}
	AnalysisContent(content, c)
	e := tool.CurrentTimeMillis()
	tool.ShowTime(s, e)
}

func main() {
	ExampleScrape()
}

//分析算法
func AnalysisContent(content map[int]map[int]string, count map[int]map[int]string) {
	c, _ := tool.NewRedis()
	for k, v := range content {
		if v != nil {
			key := "content--" + strconv.Itoa(k)
			for _, v1 := range v {
				v1 = strings.Replace(v1, " ", "", -1)
				v1 = strings.Replace(v1, "\n", "", -1)
				c.Do("SADD", key, v1)
			}
		}
	}

	//利用redis 判断重复   content1 为所有页面不同元素位置
	content1 := make(map[int]map[int]int)
	conut := len(content)
	for k, v := range content {
		if v != nil {
			contentSub := make(map[int]int)
			for k1, v1 := range v {
				for i := 0; i < conut; i++ {
					key := "content--" + strconv.Itoa(i)
					v1 = strings.Replace(v1, " ", "", -1)
					v1 = strings.Replace(v1, "\n", "", -1)
					if v1 != "" {
						re, _ := c.Do("SISMEMBER", key, v1)
						var j int64 = 1
						if re != j {
							contentSub[k1]++
							//if contentSub[k1]>=4{
							//	//fmt.Println(k1)
							//	//fmt.Println("=======================\n")
							//	//fmt.Println(v1)
							//}
						}
					}
				}
			}
			content1[k] = contentSub
		}
	}

	//fmt.Println(content1)
	//判断字符串大小挑选最小满足项目

	minContent := make(map[int]map[int]int)

	for k1, v1 := range content1 {
		subContent := make(map[int]int)
		for k2, _ := range v1 {
			len := tool.LenOfString(content[k1][k2])
			if len >= config.Minstringlen {
				subContent[k2] = len
			}
			fmt.Print()
		}
		minContent[k1] = subContent
	}
	//查找中文字符多的片段

	fmt.Print(minContent[0])

	fmt.Print("\n")

	fmt.Print(content1[0])
	fmt.Print("\n")

	//fmt.Println(len(content[0][3]))
	//写出从属关系
	//第一层
	analysisCount := make(map[int]map[int]int)
	analysisString := make(map[int]map[int]string)
	countString := 0

	for k3, v3 := range content1 {
		//第二层
		count1 := make(map[int]int)
		stringMap := make(map[int]string)
		for k4, _ := range v3 {
			//变量第二层
			for k5, _ := range v3 {
				//排除自身
				if k5 != k4 {
					//判断重复
					if strings.Contains(content[k3][k5], content[k3][k4]) {
						count1[k5]++
						stringMap[countString] = strconv.Itoa(k4) + "->" + strconv.Itoa(k5)
						countString++
					}
				}
			}
		}
		analysisCount[k3] = count1
		analysisString[k3] = stringMap
	}


	//判断非相同提取段落
	countCollection:=make(map[int]map[int]int)
	for k,v:=range minContent {
		countCollectionSub:=make(map[int]int)
		for k1,_:=range v {
			for k2,_:=range content1[k]{
				if strings.Contains(content[k][k1],content[k][k2]){
					countCollectionSub[k1]++
				}
			}
		}
		countCollection[k]=countCollectionSub
	}
	fmt.Print(countCollection[0])
}
