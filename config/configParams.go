package config

//判断是否更新的集合，如果有更新优先写入该集合，然后通过UpdateList.go 定时更新到UpdateList 列表中
const UpdateListSet string = "UpdateListSet"

//UpdateList 列表用来存储更新的URL 是监控系统数据产出接口，其它系统唯一数据交互途经
const UpdateList string = "UpdateList"

//失效链接地址集合，用来排除失效链接
const BadSite string = "BadSite"

//Hash表存储监控队列详细信息
const MonitorHash string = "MonitorHash"

//需要监控的队列定时刷新推入
const MonitorList string = "MonitorList"

//用来存储网站内容
const MonitorSiteHash string = "MonitorSiteHash"
const MonitorMiddleSiteHash string = "MonitorMiddleSiteHash"

//配置参照表用来读取带有星号的链接相关配置
const MonitorShowHash = "MonitorShowHash"

//监控队列时间间隔
const MonitorTime int = 10

//单词最大处里数
const MaxProcess = 3

//采集数据前缀
const PrefixValve string = "Value:-"
//存放栏目页面链接内容
const PrefixCategory string = "cat:"

//采集链接存储集合
const ValueSet string = "ValueSet"


//敏感词存储集合
const BadWordSet = "BadWordSet"
const BadWordStoreSet = "BadWordStoreSet"

//无敏感词集合
const NullBadWordSet = "NullBadWordSet"

//链接地址集合
const HistoryUrlSet = "HistoryUrlSet"

//列表也更新集合
const CategoryUpdateSet = "CategoryUpdateSet"

//最大数据处理进程
const MaxDataProcess = 1000

//历史记录前缀
const HistoryPrefix="History:"

//查找目标链接对应栏目链接
const ContentParentHash  ="ContentParentHash"
const ContentPrefix="Content:"
//内容页采集链接集合
const ContentUrlSet="ContentUrlSet"


//数据表字段
const UrlStart  ="url_start:-"
const UrlEnd  ="url_end:-"
const Nodeid  ="nodeid:-"
const HostName ="Host:-"
const TimeHost ="Time:-"
const Sourcecharset ="sourcecharset:-"

//网页链接数要最小个数
const MinNumOfUrl=10
//链接最大差值
const MaxMinusValue=8
//栏目也无法获取链接
const NullUrlInCategory  = "NullUrlInCategory"
const NullKey  ="NullKey"

//挑取的字符串最小值
const Minstringlen=1500


const DataFilterList="node_data"
const DataFilterResult="node_data_result"
const DataBadWordList="node_data_badword_list"
const DataPrefix="node:"


type mainParams struct {
	UpdateListSet   string
	UpdateList      string
	BadSite         string
	MonitorHash     string
	MonitorList     string
	MonitorTime     int
	MonitorSiteHash string
	MaxProcess      int
}

//返回配置参数
func NewMainParams() *mainParams {
	return &mainParams{
		UpdateList:      UpdateList,
		UpdateListSet:   UpdateListSet,
		BadSite:         BadSite,
		MonitorHash:     MonitorHash,
		MonitorList:     MonitorList,
		MonitorTime:     MonitorTime,
		MonitorSiteHash: MonitorSiteHash,
		MaxProcess:      MaxProcess,

	}

}
