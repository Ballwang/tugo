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
