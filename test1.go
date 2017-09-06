package main

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"github.com/Ballwang/tugo/tool"
)

var completeCat1 chan int = make(chan int)
func main() {
	s:=tool.CurrentTimeMillis()
	for i:=0;i<=1000;i++ {
		go test(i)
	}

	for i:=0;i<=1000;i++ {
		<-completeCat1
	}
	e:=tool.CurrentTimeMillis()
	tool.ShowTime(s,e)
}

func test(i int)  {
	resp, err := http.Get("https://www.jubi.com/api/v1/ticker/?coin=btc")
	if err != nil {
		tool.ErrorPrint(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}
	fmt.Println(resp.StatusCode)
	fmt.Println(string(body))
	fmt.Println(i)
	completeCat1<-1
}


