package tool

import (
	"regexp"
	"fmt"
	"strings"
	"github.com/PuerkitoBio/goquery"
)

////获取网页上所div结构
//func divGet(html string)[]string  {
//
//	return []string
//}

//提取A链接
func ListHref(html string) []string {
	var hrefRegexp = regexp.MustCompile("<[a|A].*?[href|HREF]=.*?>.*?[^<.*?>]?</[a|A]>")
	match := hrefRegexp.FindAllString(html, -1)
	result := []string{}
	if match != nil {
		for _, v := range match {
			result = append(result, v)
		}
	}
	return result
}

//获取该标签所有内容
func GetLabel(html string, label string) []string {
	var hrefRegexp = regexp.MustCompile("<\\s*?" + label + ".*?>.*?[^<*?.*?>*?]/" + label + ">")

	match := hrefRegexp.FindAllString(html, -1)
	result := []string{}
	if match != nil {
		for i, v := range match {
			fmt.Println("[", i, "]-", v)
			//result=append(result,v)
		}
	}
	return result

}

//获取链接文本
func GetALabelText(string string) string {
	reg := regexp.MustCompile(`<a.*\n*href=.*?>`)
	aheader := reg.FindAllString(string, -1)
	if aheader != nil {
		stringArray := strings.Split(string, aheader[0])
		if len(stringArray) > 0 {
			for _, v := range stringArray {
				if v != "" {
					aText := strings.Split(v, "</a>")
					return aText[0]
				}
			}
		}
	} else {

	}

	return ""
}

//获取链接Url
func GetALabelUrl(string string) string {
	string = strings.ToLower(string)
	reg := regexp.MustCompile(`href=".*?\s*"`)
	url := reg.FindAllString(string, -1)
	if url == nil {
		reg = regexp.MustCompile(`href='.*?\s*'`)
		url = reg.FindAllString(string, -1)
		if url != nil {
			return strings.Replace(strings.Replace(url[0], "href='", "", -1), "'", "", -1)
		} else {
			reg = regexp.MustCompile(`href=.*?\s`)
			url = reg.FindAllString(string, -1)

			fmt.Print("3333\n")
			fmt.Print(url)
			if url != nil {
				return strings.Replace(strings.Replace(url[0], "href=", "", -1), "'", "", -1)
			} else {
				return ""
			}
		}
	}
	return strings.Replace(strings.Replace(url[0], "href=\"", "", -1), "\"", "", -1)
}

//获取host链接地址
func GetHostNameByUrl(url string) string {
	urlString := strings.Replace(url, "http://", "", -1)
	if urlString != "" {
		array := strings.Split(urlString, "/")
		if array[0] != "" {
			return array[0]
		}
	}
	return ""
}

//去除链接中host
func RemovHostNameByUrl(url string) string {
	hostString := GetHostNameByUrl(url)
	url = strings.Replace(url, "http://", "", -1)
	url = strings.Replace(url, hostString, "", -1)
	return url
}

//获取host链接地址后缀
func GetHostUrlSuffix(url string) string {
	urlArrary := strings.Split(url, "/")
	lenCount := len(urlArrary)
	var subFix string
	if lenCount > 0 {
		subArray := strings.Split(urlArrary[lenCount-1], ".")
		subConut := len(subArray)
		if subConut > 1 {
			subFix = subArray[subConut-1]
		} else {
			subFix = ""
		}
	}
	return subFix
}

//获取链接的uri
func GetHostUri(url string) string {
	urlarray := strings.Split(url, "/")
	count := len(urlarray)
	subfixArray := strings.Split(urlarray[count-1], ".")
	subCount := len(subfixArray)
	if subCount > 1 {
		return strings.Replace(url, urlarray[count-1], "", -1)
	} else if urlarray[count-1] == "" {
		return url
	}
	return ""

}

//组合绝对地址
func GetAbsoluteUrl(host string, url string) string {
	hostName := GetHostNameByUrl(host)
	//判断如果有完整的http://或者存在主域名的情况下直接返回链接
	if strings.Contains(url, "http://") || strings.Contains(url, hostName) {
		return url
	}
	//其它情况处理
	sArray := strings.Split(url, "/")
	if len(sArray) > 0 {
		if sArray[0] == "" {
			return "http://" + hostName + url
		} else if sArray[0] == ".." {
			hostUri := GetHostUri(host)
			return hostUri + url
		} else if sArray[0] == "." {
			url = strings.Replace(url, "./", "", 1)
			hostUri := GetHostUri(host)
			return hostUri + url
		} else {
			hostUri := GetHostUri(host)
			return hostUri + url
		}
	}
	return ""
}


//获取指定 节点下相关节点所有数据
func GetHtmlDom(nodeName string, url string, charset string) (map[int]string , map[int]string) {
	s := CurrentTimeMillis()
	doc, err := goquery.NewDocument(url)
	if err != nil {
		Error(err)
	}

	// Find the review items


	contentArray:=make(map[int]string)
	c:=make(map[int]string)

	doc.Find(nodeName).Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		//a, _ := s.Find(findNode).Html()
		a,_:= s.Html()
		p,isAttr:=s.Attr("id")
		if !isAttr{
			p,isAttr=s.Attr("class")
			if !isAttr{
				p=""
			}
		}

		//fmt.Println(p)
		//fmt.Println("++++++++++++++++++++++\n")
		band := GetStringToUtf8(a, charset)
		//fmt.Println(band)
		//fmt.Println("=====================\n")
		if a != "" && p!=""{
			contentArray[i] = band
			c[i]=p
		}
		//fmt.Printf("Review %d: %s -------------------\n", i, band)

	})
	e := CurrentTimeMillis()
	ShowTime(s, e)
	return contentArray,c
}
