package tool

import (
	"net"
	"fmt"
	"os"
	"strings"
)

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
