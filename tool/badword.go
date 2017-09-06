package tool

import "strings"

//查找字符串是否存在，如果存在返回true
func FindBadWord(content string,badword string) bool  {
	return strings.Contains(content,badword)
}

//多个关键词查找
func FindBadWordAny(content string,badword string) bool  {
	return strings.ContainsAny(content,badword)
}

//根据关键词分片查找敏感词是否存在
func FindBadWordSet(content string,badwordSet []string) bool  {
	if len(badwordSet)>0{
		for _,v:=range badwordSet{
			if FindBadWord(content,v){
				return true
			}
		}
	}
	return false
}

