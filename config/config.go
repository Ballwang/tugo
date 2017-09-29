package config

import (
	"bufio"
	"io"
	"os"
	"strings"
)

const middle = "-"

type Config struct {
	ConfigMap map[string]string
	strcet    string
}

//初始化Config，配置文件地址是写程序的相对地址
func NewConfig() *Config {
	//分配内存
	c := new(Config)
	c.ConfigMap = make(map[string]string)
	path := "./config/config.ini"
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		s := strings.TrimSpace(string(b))
		//fmt.Println(s)
		if strings.Index(s, "#") == 0 {
			continue
		}

		n1 := strings.Index(s, "[")
		n2 := strings.LastIndex(s, "]")
		if n1 > -1 && n2 > -1 && n2 > n1+1 {
			c.strcet = strings.TrimSpace(s[n1+1: n2])
			continue
		}

		if len(c.strcet) == 0 {
			continue
		}
		index := strings.Index(s, "=")
		if index < 0 {
			continue
		}

		frist := strings.TrimSpace(s[:index])
		if len(frist) == 0 {
			continue
		}
		second := strings.TrimSpace(s[index+1:])

		pos := strings.Index(second, "\t#")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, " #")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, "\t//")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, " //")
		if pos > -1 {
			second = second[0:pos]
		}

		if len(second) == 0 {
			continue
		}
		key := c.strcet + middle + frist
		c.ConfigMap[key] = strings.TrimSpace(second)
	}
	return c
}

//获取配置文件详细信息
func (c *Config) GetConfig(section string, mapKey string) string {
	key := section + middle + mapKey
	v, found := c.ConfigMap[key]
	if !found {
		return ""
	}

	return v
}
