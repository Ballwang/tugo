package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"os/exec"
	"bytes"
	"os"
	"net"
	"strings"
)

var pass = "lrb123"

var rootDir = "/server/"

var marathonHost="192.168.3.21:8080"
//delete by Id "curl -X DELETE  192.168.3.21:8080/v2/apps/nginxweb"

//开始获取网站内容

func GetHook(w http.ResponseWriter, req *http.Request) {

	var hook map[string]interface{}
	body, _ := ioutil.ReadAll(req.Body)

	json.Unmarshal(body, &hook)

	if pass != hook["password"] {
		return
	}

	var name = ""
	var git_ssh_url = ""
	var cmd = "#!/bin/sh\n"

	if hook["project"] != "" {
		//fmt.Println(hook["project"].(map[string]interface{})["path_with_namespace"])
		//path_with_namespace:=hook["project"].(map[string]interface{})["path_with_namespace"]
		git_ssh_url = hook["project"].(map[string]interface{})["git_ssh_url"].(string)
		name = hook["project"].(map[string]interface{})["name"].(string)
	}
	path := rootDir + name+"/"
	ok, _ := pathExists(path)
	if !ok {
		cmd += "mkdir " + rootDir + "\n"
		cmd += "cd " + rootDir + "\n"
		cmd += "git clone " + git_ssh_url + "\n"
	}
	cmd += "cd " + path + "\n"
	cmd += "git pull\n"

	d1 := []byte(cmd)
	ioutil.WriteFile(name+".sh", d1, 0755)

	run_shell(name)

	//写shell
	start := "#!/bin/sh\n"
	start += "cd /\n"
	start += "cd /go/src/\n" //这里使用name变量必须和版本控制名称和执行脚本名称相同
	start += "./" + name + "\n"

	dstart := []byte(start)
	ioutil.WriteFile(path+"start.sh", dstart, 0755)

	//写shell
	docker := "FROM 192.168.3.54:5000/centos7-go1.9 \n"
	docker += "MAINTAINER Ballwang  ballwang@foxmail.com\n" //这里使用name变量必须和版本控制名称和执行脚本名称相同
	docker += "ADD ./config/ /go/src/config/\n"
	docker += "ADD ./config/json/ /go/src/config/json/\n"
	docker += "ADD " + name + " /go/src/\n"
	docker += "ADD start.sh /go/src/\n"
	docker += "RUN chmod 755 /go/src/\n"

	docker += "CMD /go/src/start.sh"

	ddocker := []byte(docker)
	ioutil.WriteFile(path+"Dockerfile", ddocker, 0755)

	//构建docker镜像
	dockerImage:="cd "+path+"\n"
	dockerImage+="docker build -t 192.168.3.54:5000/"+strings.ToLower(name)+":1.0 . \n"
	dockerImage+="docker push 192.168.3.54:5000/"+strings.ToLower(name)+":1.0\n"
	//删除tag为none的镜像
	dockerImage+="docker rmi $(docker images -f \"dangling=true\" -q)\n"
	dockerImageStart := []byte(dockerImage)
	ioutil.WriteFile(name+"-docker.sh", dockerImageStart, 0755)
	run_shell(name+"-docker")

	//创建marathon
	marathon:="cd "+path+"\n"
	marathon+="curl -i -X DELETE "+marathonHost+"/v2/apps/"+strings.ToLower(name)+"\n"
	marathon+="sleep 6\n"
	marathon+="curl -i -H 'Content-Type: application/json' -d@config/json/"+name+".json "+marathonHost+"/v2/apps\n"

	marathonStart := []byte(dockerImage)
	ioutil.WriteFile(name+"-marathon.sh", marathonStart, 0755)
	run_shell(name+"-marathon")


	//生成marathon create json
	//使用 json 创建镜像 curl -i -H 'Content-Type: application/json' -d@1.json 192.168.3.21:8080/v2/apps

	fmt.Println("11111111111111")

}

//判断文件夹是否存在

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func run_shell(name string) {
	cmdDO := exec.Command("/bin/bash", "./"+name+".sh")
	var out bytes.Buffer
	cmdDO.Stdout = &out
	cmdDO.Run()
}

//执行shell

func exec_shell(arg []string) {
	cmd := exec.Command("/bin/bash", arg...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s", out.String())
}

//返回本机IP地址
func GetIP() string {
	conn, err := net.Dial("udp", "google.com:80")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	defer conn.Close()
	return strings.Split(conn.LocalAddr().String(), ":")[0]
}

func main() {

	ip := GetIP()
	http.HandleFunc("/StartGetContent", GetHook)
	http.HandleFunc("/", GetHook)

	err := http.ListenAndServe(ip+":7010", nil)

	if err != nil {
		fmt.Println("Listen And Serve error: ", err.Error())
	}

}
