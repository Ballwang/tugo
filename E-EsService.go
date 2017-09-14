package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/Ballwang/tugo/tool"
	"context"
	"gopkg.in/olivere/elastic.v5"
	"github.com/Ballwang/tugo/config"
	"github.com/Ballwang/tugo/soft/softServer"
	consulapi "github.com/hashicorp/consul/api"
	"strconv"
)

var token = "eo99s001144999999381111"

//推送请求结果
func ShowRequest(w http.ResponseWriter, r interface{}) {
	bytes, _ := json.Marshal(r)
	fmt.Fprint(w, string(bytes))
}

//权限验证
func CheckRight(w http.ResponseWriter, req *http.Request) {
	is := true
	key, isToken := req.Form["token"]
	if !isToken {
		is = false
	} else {
		if key[0] != token {
			is = false
		}
	}

	result := make(map[string]interface{})
	if !is {
		result["success"] = "false"
		result["code"]="401"
		result["message"] = "权限验证失败！"
		ShowRequest(w, result)
		return
	}
}

//检查特定字段是否为空
func CheckValue(name string, w http.ResponseWriter, req *http.Request) bool {
	_, is := req.Form[name]
	result := make(map[string]interface{})
	if !is {
		result["success"] = "false"
		result["code"]="418"
		result["message"] = name + " 必须填写！"
		ShowRequest(w, result)
		return false
	}
	return true
}

//添加数据
func AddEsData(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	result := make(map[string]interface{})
	_, isToken := req.Form["token"]
	if !isToken {
		result["success"] = "false"
		result["code"]="401"
		result["message"] = "权限验证失败！"
		ShowRequest(w, result)
		return
	}

	//判断请求类型
	//if req.Method == "POST" {
	if req.Method == "POST" {
		//data:=req.PostForm
		data := req.PostForm

		siteName, isSiteName := req.Form["siteName"]
		if !isSiteName {
			result["success"] = "false"
			result["code"]="418"
			result["message"] = "SiteName 必须填写！"
			ShowRequest(w, result)
			return
		}
		contentID, _ := req.Form["ContentID"]

		////获取参数并且映射到结构体中
		c := &tool.Content{}
		c.InitContentByReq(data)
		r, e := c.Check()
		if !r {
			result["success"] = "false"
			result["code"]="419"
			result["message"] = e
			ShowRequest(w, result)
			return
		}

		cxt := context.Background()
		es := tool.NewESFromConfig()
		idString := tool.Md5String(siteName[0] + contentID[0])
		s, err := es.AddData(cxt, c, idString, es.Index, es.Type)
		if err != nil {
			result["success"] = "false"
			result["code"]="510"
			result["message"] = "数据添加失败！"
			ShowRequest(w, result)
			return
		} else {
			result["success"] = "true"
			result["code"]="200"
			result["message"] = s
			ShowRequest(w, result)
			return
		}

	} else {
		result["success"] = "false"
		result["code"]="420"
		result["message"] = "提交类型错误！"
		ShowRequest(w, result)
		return
	}

}

func DelEsData(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	CheckRight(w, req)
	result := make(map[string]interface{})

	siteName, isSiteName := req.Form["siteName"]
	if !isSiteName {
		result["success"] = "false"
		result["code"]="418"
		result["message"] = "SiteName 必须填写！"
		ShowRequest(w, result)
		return
	}

	contentID, isContentID := req.Form["ContentID"]
	if !isContentID {
		result["success"] = "false"
		result["code"]="418"
		result["message"] = "ContentID 必须填写！"
		ShowRequest(w, result)
		return
	}

	idString := tool.Md5String(siteName[0] + contentID[0])
	cxt := context.Background()
	es := tool.NewESFromConfig()
	isSuccess := es.DeleteDataById(cxt, idString, es.Index, es.Type)
	if isSuccess {
		result["success"] = "true"
		result["code"] = "200"
		result["message"] = "文章ID：" + contentID[0] + " 删除成功！"
		ShowRequest(w, result)
		return
	} else {
		result["success"] = "false"
		result["code"] = "421"
		result["message"] = "文章ID：" + contentID[0] + " 不存在！"
		ShowRequest(w, result)
		return
	}

}

//查询搜索结果采用IK 分词查询
func SearchEsData(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	CheckRight(w, req)

	//result:=make(map[string]interface{})

	if is := CheckValue("query", w, req); !is {
		return
	}
	if is := CheckValue("field", w, req); !is {
		return
	}

	re, isfield2 := req.Form["field2"]
	field2 := ""
	query := req.Form["query"][0]
	field1 := req.Form["field"][0]
	if isfield2 {
		field2 = re[0]
	}

	cxt := context.Background()
	es := tool.NewESFromConfig()

	termQuery := elastic.NewMultiMatchQuery(query, field1, field2)
	termQuery.Type("best_fields")
	termQuery.TieBreaker(0.5)
	params := config.NewConfig()
	preTags := params.GetConfig("elasticsearch", "preTags")
	postTags := params.GetConfig("elasticsearch", "postTags")

	h := elastic.NewHighlight()
	h.PreTags(preTags)
	h.PostTags(postTags)
	h.Field(field1)
	h.Field(field2)

	hits, err := es.SearchData(cxt, termQuery, es.Index, es.Type, h)
	if err != nil {

	}

	if hits.TotalHits > 0 {
		r, _ := json.Marshal(hits)
		fmt.Fprint(w, string(r))
	}

}

//查询搜索结果采用IK 分词查询
func SearchEsDataHighScore(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	CheckRight(w, req)

	result := make(map[string]interface{})

	if is := CheckValue("query", w, req); !is {
		return
	}
	if is := CheckValue("field", w, req); !is {
		return
	}

	query := req.Form["query"][0]
	field1 := req.Form["field"][0]

	cxt := context.Background()
	es := tool.NewESFromConfig()

	termQuery := elastic.NewMultiMatchQuery(query, field1)
	termQuery.Type("best_fields")
	termQuery.TieBreaker(0.5)
	params := config.NewConfig()
	preTags := params.GetConfig("elasticsearch", "preTags")
	postTags := params.GetConfig("elasticsearch", "postTags")

	h := elastic.NewHighlight()
	h.PreTags(preTags)
	h.PostTags(postTags)
	h.Field(field1)

	hits, err := es.SearchData(cxt, termQuery, es.Index, es.Type, h)
	if err != nil {

	}



	if hits!=nil&&hits.TotalHits > 0 {
		// Iterate through results
		for _, hit := range hits.Hits {
			// hit.Index contains the name of the index

			// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
			var t tool.Content
			err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				// Deserialization failed
			}

			if *hit.Score == *hits.MaxScore {
				result["success"] = "true"
				result["code"] = "200"
				result["message"] = "记录查找成功！"
				result["score"] = *hits.MaxScore
				result["hightlight"] = hit.Highlight
				result["data"] = t
				resultWithData, _ := json.Marshal(result)
				fmt.Fprint(w, string(resultWithData))
				return
			}
		}


	} else {
		result["success"] = "false"
		result["message"] = "未查找到相关记录！"
		r, _ := json.Marshal(result)
		fmt.Fprint(w, string(r))
	}
}

func UpdateEsData(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	result := make(map[string]interface{})
	_, isToken := req.Form["token"]
	if !isToken {
		result["success"] = "false"
		result["code"]="401"
		result["message"] = "权限验证失败！"
		ShowRequest(w, result)
		return
	}

	//判断请求类型
	//if req.Method == "POST" {
	if req.Method == "POST" {
		//data:=req.PostForm
		data := req.PostForm

		siteName, isSiteName := req.Form["siteName"]
		if !isSiteName {
			result["success"] = "false"
			result["code"]="418"
			result["message"] = "SiteName 必须填写！"
			ShowRequest(w, result)
			return
		}
		contentID, _ := req.Form["ContentID"]

		////获取参数并且映射到结构体中
		c := &tool.Content{}
		c.InitContentByReq(data)
		r, e := c.Check()
		if !r {
			result["success"] = "false"
			result["code"]="419"
			result["message"] = e
			ShowRequest(w, result)
			return
		}

		cxt := context.Background()
		es := tool.NewESFromConfig()
		idString := tool.Md5String(siteName[0] + contentID[0])
		s, err := es.AddData(cxt, c, idString, es.Index, es.Type)
		if err != nil {
			result["success"] = "false"
			result["code"]="510"
			result["message"] = "数据添加失败！"
			ShowRequest(w, result)
			return
		} else {
			result["success"] = "true"
			result["code"]="200"
			result["message"] = s
			ShowRequest(w, result)
			return
		}

	} else {
		result["success"] = "false"
		result["code"]="420"
		result["message"] = "提交类型错误！"
		ShowRequest(w, result)
		return
	}
}

func main() {
	//第一个参数为客户端发起http请求时的接口名，第二个参数是一个func，负责处理这个请求。
	http.HandleFunc("/Add", AddEsData)
	http.HandleFunc("/Del", DelEsData)
	http.HandleFunc("/Search", SearchEsData)
	http.HandleFunc("/SearchOne", SearchEsDataHighScore)
	http.HandleFunc("/Update", UpdateEsData)
	//服务器要监听的主机地址和端口号
	//配置注册服务器信息
	ip:=tool.GetIP()
	registration :=new(consulapi.AgentServiceRegistration)
	registration.ID="es-service"
	registration.Name="ES 全文搜索引擎接口服务"
	registration.Address=ip
	registration.Port=8081
	registration.Tags=[]string{"Elasticsearch 接口服务器地址！"}
	registration.Check=&consulapi.AgentServiceCheck{
		//TCP:                          fmt.Sprintf("http://%s:%d", registration.Address, checkPort),
		TCP:                            ip+":"+strconv.Itoa(registration.Port),
		Timeout:                        "3s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "5s", //check失败后5秒删除本服务
	}
	softServer.RegisterMcService(registration)
	fmt.Println("Server starting at :"+ip+":"+strconv.Itoa(registration.Port))
	err := http.ListenAndServe(ip+":"+strconv.Itoa(registration.Port), nil)

	if err != nil {
		fmt.Println("ListenAndServe error: ", err.Error())
	}
}
