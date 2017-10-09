package tool

import (
	"encoding/json"

	"io/ioutil"
	"strings"
)

func MakeMarathonJson(Id string,cpus float64,mem float64,instances int,port int)  {
	j:=make(map[string]interface{})
	container:=make(map[string]interface{})
	docker:=make(map[string]interface{})
	portMappings:=make(map[string]interface{})
	healthChecks:=make(map[string]interface{})

	j["id"]=strings.ToLower(Id)
	j["cpus"]=cpus
	j["mem"]=mem
	j["instances"]=instances

	j["container"]="DOCKER"
	portMappings["containerPort"]=port
	portMappings["hostPort"]=0
	portMappings["servicePort"]=port
	portMappings["protocol"]="tcp"

	docker["image"]="192.168.3.54:5000/"+strings.ToLower(Id)+":1.0"
	docker["network"]="HOST"
	docker["forcePullImage"]=true
	docker["portMappings"]=[]interface{}{portMappings}

	container["type"]="DOCKER"
	container["docker"]=docker


	healthChecks["protocol"]="TCP"
	healthChecks["gracePeriodSeconds"]=3
	healthChecks["intervalSeconds"]=5
	healthChecks["port"]=port
	healthChecks["timeoutSeconds"]=6
	healthChecks["maxConsecutiveFailures"]=3
	j["container"]=container
	j["healthChecks"]=[]interface{}{healthChecks}

	b, _ := json.Marshal(j)

	ioutil.WriteFile("./config/json/"+Id+".json", b, 0755)
}
