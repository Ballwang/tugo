package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"os/exec"
	"bytes"
)

var Cipwd  = "imooc"

func Githook(w http.ResponseWriter,req *http.Request)  {
	var hook map[string]interface{}
	//获取返回结果
	body,err:=ioutil.ReadAll(req.Body)
	if (err!=nil){
		fmt.Println("Read body failed!")
	}

	//映射到hook返回结果映射到map中
	json.Unmarshal(body,&hook)

	//验证pass  是否合法
	if(Cipwd !=hook["password"]){
		fmt.Println("Pass word not allowed!")
	}

	//制作shell 脚本
	cmd := exec.Command("/bin/bash", "")
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s", out.String())

	
}




func  main ()  {
	ip  := "139.162.120.224"
	port :="7010"
	http.HandleFunc("/",Githook)
	err:=http.ListenAndServe(ip+port,nil)
	if (err !=nil){
		fmt.Println("Server is listening on:"+ip+":"+port)
	}
}
