package main

import (
	"github.com/Ballwang/tugo/tool"
	"context"
	"fmt"


)

func main() {
	// Create a context
	//createIndex()
	s:=tool.CurrentTimeMillis()
	del()
	//addData()
	//GetData()
	//del()
	//GetD()
	e:=tool.CurrentTimeMillis()
	tool.ShowTime(s,e)

}

func GetD()  {
	//cxt := context.Background()
	//es := &tool.ElasticEsearch{Host:"192.168.4.80", Port:"9200"}
	//termQuery := elastic.NewMultiMatchQuery("专查食药问题警察怎办案", "Title", "Content")
	//termQuery.Type("best_fields")
	//termQuery.TieBreaker(0.5)
	//
	//h:=elastic.NewHighlight()
	//h.PreTags("<lrbstring>")
	//h.PostTags("</lrbstring>")
	//h.Field("Title")
	//h.Field("Content")
	//
	//result,err:=es.SearchData(cxt,termQuery,"artical","detial",h)
	//if err!=nil{
	//
	//}
	//fmt.Println(result.MaxScore)
}


func del()  {
	cxt := context.Background()
	es := &tool.ElasticEsearch{Host:"192.168.3.92", Port:"9200"}
	md5String:=tool.Md5String("jiangsu91186")
	fmt.Println(md5String)
	es.DeleteDataById(cxt,md5String,"artical","detial")
}

func GetData()  {
	cxt := context.Background()
	es := &tool.ElasticEsearch{Host:"192.168.4.80", Port:"9200"}
	md5String:=tool.Md5String("jiangsu10517477")
	es.GetDataExsit(cxt,md5String,"artical","detial")
}


func addData()  {
	cxt := context.Background()
	es := &tool.ElasticEsearch{Host:"192.168.4.80", Port:"9200"}
	content:=&tool.Content{
		Title: "江苏在京办宣介会 李强介绍创新驱动转型发展",
		TitleSub: "江苏在京办宣介会 李强介绍创新驱动转型发展",
		Keywords:"宣介会",
		Description:"8月31日，中共中央对外联络部在京举行“中国共产党的故事：江苏省委的实践——创新驱动、转型发展”专题宣介会，邀请中共江苏省委向访华的外国政党政要、驻华使馆高级外交官、外国商会",
		Author:"DH011",
		Time:"2017-09-01 08:25:39",
		Content:"省委书记李强作宣介，中联部部长宋涛致辞。来自美国、俄罗斯、芬兰、丹麦、挪威、西班牙、罗马尼亚、白俄罗斯、塞尔维亚、黑山、阿尔巴尼亚、乌克兰、日本、印尼、柬埔寨、斯里兰卡、印度、尼泊尔、孟加拉国、伊朗、摩洛哥、布隆迪等国家和欧洲议会的政党代表、学者以及外国驻华高级外交官、外国商会和跨国公司驻华代表等参加宣介会。中联部副部长郭业洲主持。"+
			"宋涛在致辞中表示，创新是推动发展的重要力量，更是许多国家发展战略的核心内容和主要目标。习近平总书记多次强调，惟改革者进，惟创新者强，惟改革创新者胜。党的十八大以来，以习近平同志为核心的中共中央把创新摆在国家发展全局的核心位置，大力推进以科技创新为核心的全面创新，中国已跃上创新大国之路。",
		CopyFrom:"http://xh.xhby.net/mp3/pc/c/201709/01/c370976.html",
		Url:"http://jiangsu.china.com.cn/html/jsnews/around/10517477_1.html",
		ContentID:"10517477",
	 }
	md5String:=tool.Md5String("jiangsu10517477")
	es.AddData(cxt,content,md5String,"artical","detial")
}


func createIndex()  {
	cxt := context.Background()
	es := &tool.ElasticEsearch{Host:"192.168.3.92", Port:"9200"}
	mapString := `
	    {
	       "settings":{
	           "number_of_shards":9,
		       "number_of_replicas":0
	       },
	       "mappings":{
	           "_default_": {
			      "_all": {
				    "enabled": false
			      }
		       },
		       "detial":{
		          "properties":{
		              "Title":{
                          "type": "text",
                          "analyzer": "ik_max_word",
                          "search_analyzer": "ik_max_word"
		              },
		              "TitleSub":{
                          "type": "text",
                          "analyzer": "ik_max_word",
                          "search_analyzer": "ik_max_word"

		              },
		              "Description":{
                          "type": "text",
                          "analyzer": "ik_max_word",
                          "search_analyzer": "ik_max_word"

		              },
		              "Author":{
                          "type": "text",
                          "analyzer": "ik_max_word",
                          "search_analyzer": "ik_max_word"
		              },
		              "Time":{
                          "type": "text"
		              },
		              "CopyFrom":{
                          "type": "text"
		              },
		              "Content":{
                          "type": "text",
                          "analyzer": "ik_max_word",
                          "search_analyzer": "ik_max_word"
		              },
		              "Url":{
                          "type": "text"
		              },
		              "ContentID":{
                          "type": "text"
		              }
		          }
		       }
	       }
	    }
	`

	err:=es.CreateIndex(cxt, "artical", mapString)
	if err!=nil{
		fmt.Println(err)
	}
}



