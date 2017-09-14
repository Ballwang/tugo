package tool

import (
	"time"
	"github.com/chasex/redis-go-cluster"

)

//创建集群对象
func NewRedisCluster() (*redis.Cluster, error) {
	c, err := redis.NewCluster(
		&redis.Options{
			StartNodes:   []string{"192.168.4.83:7001", "192.168.4.83:7002", "192.168.4.84:7003", "192.168.4.85:7004", "192.168.4.85:7005", "192.168.4.85:7005"},
			ConnTimeout:  50 * time.Millisecond,
			ReadTimeout:  50 * time.Millisecond,
			WriteTimeout: 50 * time.Millisecond,
			KeepAlive:    16,
			AliveTime:    60 * time.Second,
		})
	return c, err
}

//获取hash列表所有数据
func RedisClusterHGETALL(Key string) []string {
	c, _ := NewRedisCluster()
	defer c.Close()
	listString := []string{}
	reply, err := redis.StringMap(c.Do("HGETALL", Key))
	if err != nil {

	} else {
		for _, v := range reply {
			listString = append(listString, string(v))
		}
	}
	return listString
}

//获取集合中所有元素
func RedisClusterSMEMBERS(Key string) []string {
	c, _ := NewRedisCluster()
	defer c.Close()
	listString := []string{}
	reply, err := redis.Values(c.Do("SMEMBERS", Key))
	if err != nil {

	} else {
		i := 0
		for _, v := range reply {
			listString = append(listString, string(v.([]byte)))
			i++
		}
	}
	return listString

}

//获取单个hash表值 并且删除原有记录
func RedisClusterGetHashValueAndDel(hashName string, key string) (string, bool) {
	c, _ := NewRedisCluster()
	defer c.Close()
	reply, err := redis.String(c.Do("HGET", hashName, key))
	if err != nil {

	} else {
		if key!=""{
			c.Do("HDEL",hashName,key)
		}
		return reply, true
	}
	return "", false

}

//获取单个hash表值 并且删除原有记录
func RedisClusterGetHashValueAndDelWithReids(c *redis.Cluster,hashName string, key string) (string, bool) {

	reply, err := redis.String(c.Do("HGET", hashName, key))
	if err != nil {

	} else {
		if key!=""{
			c.Do("HDEL",hashName,key)
		}
		return reply, true
	}
	return "", false

}


//获取单个hash表值
func RedisClusterGetHashValue(hashName string, key string) (string, bool) {
	c, _ := NewRedisCluster()
	defer c.Close()
	reply, err := redis.String(c.Do("HGET", hashName, key))
	if err != nil {

	} else {
		return reply, true
	}
	return "", false

}

//获取单个hash表值
func RedisClusterGetHashValueWithRedis(c *redis.Cluster,hashName string, key string) (string, bool) {

	reply, err := redis.String(c.Do("HGET", hashName, key))
	if err != nil {

	} else {
		return reply, true
	}
	return "", false

}
