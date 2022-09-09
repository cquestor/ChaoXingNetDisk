package main

import (
	"ChaoXingNetDisk/apis"
	"fmt"
)

func main() {
	cookie, err := apis.Login("13914308903", "app5896302")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	apis.NewFolder(cookie, "测试2")
}
