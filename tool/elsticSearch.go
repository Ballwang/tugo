package tool

import (
	"gopkg.in/olivere/elastic.v5"
	"golang.org/x/net/context"
	"fmt"
	"github.com/Ballwang/tugo/config"

	"github.com/mitchellh/mapstructure"
	"reflect"

)

type ElasticEsearch struct {
	Host string
	Port string "9200"
	Index string "artical"
	Type string "detail"
}

//创建ES客户端
func (c *ElasticEsearch) NewElasticEsearch() (*elastic.Client, error) {
	client, err := elastic.NewClient(elastic.SetURL("http://" + c.Host + ":" + c.Port))
	return client, err
}

//通过配置文件创建客户端
func NewESFromConfig() *ElasticEsearch {
	params := config.NewConfig()
	eslHost := params.GetConfig("elasticsearch", "eshost")
	eslPort := params.GetConfig("elasticsearch", "esport")
	eslIndex := params.GetConfig("elasticsearch", "index")
	eslType := params.GetConfig("elasticsearch", "type")
	return &ElasticEsearch{Host: eslHost, Port: eslPort,Index:eslIndex,Type:eslType}
}

//创建索引
func (c *ElasticEsearch) CreateIndex(cxt context.Context, indexString string, mapString string) error {
	client, err := elastic.NewClient(elastic.SetURL("http://" + c.Host + ":" + c.Port))
	if err != nil {
		return err
	}
	exists, err := client.IndexExists(indexString).Do(cxt)
	if !exists {
		// Index does not exist yet.
		createIndex, err := client.CreateIndex(indexString).Body(mapString).Do(cxt)
		if err != nil {
			return err
		}
		if !createIndex.Acknowledged {
			fmt.Println(indexString + "索引创建失败！")
		} else {
			fmt.Println(indexString + " 索引创建成功- mapping 内容为：" + mapString)
		}
	} else {
		fmt.Println(indexString + " 该索引已经被创建!\n")
	}
	return nil
}

//增加数据
func (c *ElasticEsearch) AddData(cxt context.Context, content *Content, idString string, indexString string, typeString string)(string,error)  {
	client, err := c.NewElasticEsearch()
	if err != nil {

		return "",err
	}
	put, err := client.Index().
		Index(indexString).
		Type(typeString).
		Id(idString).
		BodyJson(content).
		Do(cxt)
	if err != nil {
		// Handle error

		return "",err
	}

	string:=content.ContentID+" 添加成功！全文系统ID："+put.Id
	return string,nil
}

//更新数据 与 add 函数一样，这里只做区分
func (c *ElasticEsearch) UpdateData(cxt context.Context, content *Content, idString string, indexString string, typeString string) error {
	client, err := c.NewElasticEsearch()
	if err != nil {
		fmt.Println(err)
		return err
	}
	put, err := client.Index().
		Index(indexString).
		Type(typeString).
		Id(idString).
		BodyJson(content).
		Do(cxt)
	if err != nil {
		// Handle error
		fmt.Println(err)
		return err
	}
	fmt.Printf(indexString+" "+typeString+" %s to index %s, type %s\n", put.Id, put.Index, put.Type)
	return err
}

//获取指定数据是否纯在
func (c *ElasticEsearch) GetDataExsit(cxt context.Context, IdString string, indexString string, typeString string) (bool, error) {
	client, err := c.NewElasticEsearch()
	if err != nil {
		fmt.Println(err)
	}

	get, err := client.Get().
		Index(indexString).
		Type(typeString).
		Id(IdString).
		Do(cxt)
	if err != nil {
		// Handle error
		fmt.Println(err)
	}
	if get.Found {
		fmt.Printf("Got document %s in version %d from index %s, type %s\n", get.Id, get.Version, get.Index, get.Type)
		fmt.Println(get.Fields)
		return true, nil
	}
	return false, nil
}

//删除特定数据
func (c *ElasticEsearch) DeleteDataById(cxt context.Context, IdString string, indexString string, typeString string) bool {
	client, err := c.NewElasticEsearch()
	if err != nil {
		fmt.Println(err)
	}
	del, err := client.Delete().
		Index(indexString).
		Type(typeString).
		Id(IdString).
		Do(cxt)
	return del.Found
}

//查询数据
func (c *ElasticEsearch) SearchData(cxt context.Context, termQuery *elastic.MultiMatchQuery, indexString string, typeString string, highlight *elastic.Highlight) (*elastic.SearchHits, error) {
	client, err := c.NewElasticEsearch()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	searchResult, err := client.Search().
		Index(indexString).
		Query(termQuery).
		Highlight(highlight).
		Do(cxt)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	//返回查询结果
	if searchResult.Hits.TotalHits > 0 {
		return searchResult.Hits, err
	} else {
		// No hits
		fmt.Println(err)
		return nil, err
	}

}

type Content struct {
	Title       string
	TitleSub    string
	Keywords    string
	Description string
	Author      string
	Time        string
	CopyFrom    string
	Content     string
	Url         string
	ContentID   string
}

//赋值map 到结构体 map中多出的字段自动抛弃
func (c *Content) InitContentByMap(m map[string]interface{}) {
	mapstructure.Decode(m, c)
}

//赋值map 到结构体 map中多出的字段自动抛弃
func (c *Content) InitContentByReq(m map[string][]string) {
	r:=make(map[string]string)
	for k1,v1:=range m{
		for k2,v2:=range v1{
			if k2==0{
				r[k1]=v2
			}
		}
	}
	fmt.Println(r)
	mapstructure.Decode(r, c)
}

//检查结构体是否赋值
func (c *Content) Check() (bool, string) {
	resultString := ""
	isPass := true
	t := reflect.TypeOf(*c)
	v := reflect.ValueOf(*c)

	for i := 0; i < t.NumField(); i++ { //NumField取出这个接口所有的字段数量
		f := t.Field(i)               //取得结构体的第i个字段
		val := v.Field(i).Interface() //取得字段的值
		if val == "" {
			resultString = resultString + "\"" + f.Name + "\"" + ":\"不能为空!\""
			if i != (t.NumField() - 1) {
				resultString = resultString + ","
			}
			isPass = false
		}
	}
	resultString = "{" + resultString + "}"

	if isPass {
		return isPass, ""
	} else {
		return isPass, resultString
	}

}
