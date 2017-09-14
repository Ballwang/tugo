package tool

import (
	"github.com/axgle/mahonia"
	"strings"
	"regexp"
)

//转换字符串为Utf8 格式
func GetStringToUtf8(contentString string,charset string) string {
	if !strings.Contains(charset, "utf") {
		dec := mahonia.NewDecoder("gbk")
		contentString = dec.ConvertString(contentString)
	}
	return contentString
	
}


//返回字符中所有中文字符长度
func LenOfString(string string)int  {
	reg := regexp.MustCompile(`[\p{Han}]+`)
	contentString:=reg.FindAllString(string, -1)
	string1:=""

	for _,v:=range contentString{
		string1=string1+v
	}
	return len(string1)
}

//去何处换行空格，并替换
func TrimReplace(urlpage string,replace string) string  {
	urlpage=strings.Replace(urlpage,"\t",replace,-1)
	urlpage=strings.Replace(urlpage,"\r\n",replace,-1)
	urlpage=strings.Replace(urlpage,"\n",replace,-1)
	urlpage=strings.Replace(urlpage,"\r",replace,-1)
	return urlpage
}

//分割字符串
func GetSplitOne(v string,split string) []string  {
	string:=strings.Split(v, split)
	return string
}
