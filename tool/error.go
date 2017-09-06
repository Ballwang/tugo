package tool

import (
	"fmt"
	"os"
)

func Error(err error)  {
	if err!=nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(0)
	}
}

func ErrorPrint(err error)  {
	if err!=nil {
		fmt.Println("Fatal error ", err.Error())
	}
}
