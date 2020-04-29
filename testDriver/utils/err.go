package utils

import (
	"fmt"
	"os"
)

//HandleErr 处理一般错误
func HandleErr(err error, str string) {
	if err != nil {
		fmt.Println("ERROR:", str, err)
		os.Exit(1)
	}
}