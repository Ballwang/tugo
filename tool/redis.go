package tool

import (
	"github.com/garyburd/redigo/redis"
	"github.com/Ballwang/tugo/config"
)

//创建redis
func NewRedis() (c redis.Conn, err error) {
	cf := config.NewConfig()
	c, err = redis.Dial("tcp", cf.GetConfig("redis", "redisHost")+":"+cf.GetConfig("redis", "redisPort"))
	Error(err)
	return
}

//异步调用的时候需要创建和关闭后期需要优化
func DoRedis(commandName string, args ...interface{}) {
	c, err := NewRedis()
	_, err = c.Do(commandName, args[0], args[1], args[2])
	Error(err)
	defer c.Close()
}

//获取hash列表所有数据
func RedisHGETALL(Key string) []string {
	c, _ := NewRedis()
	defer c.Close()
	listString := []string{}
	reply, err := redis.Values(c.Do("HGETALL", Key))
	if err != nil {
		Error(err)
	} else {
		i := 0
		for _, v := range reply {
			if i%2 != 0 {
				listString = append(listString, string(v.([]byte)))
			}
			i++
		}
	}
	return listString
}

//获取集合中所有元素
func RedisSMEMBERS(Key string) []string {
	c, _ := NewRedis()
	defer c.Close()
	listString := []string{}
	reply, err := redis.Values(c.Do("SMEMBERS", Key))
	if err != nil {
		Error(err)
	} else {
		i := 0
		for _, v := range reply {
			listString = append(listString, string(v.([]byte)))
			i++
		}
	}
	return listString

}
