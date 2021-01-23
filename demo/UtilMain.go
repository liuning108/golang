package main

import (
	"fmt"
	"p2pserve/util"
)

func main() {
	m := map[string]interface{}{
		//消息类型
		"type": "heartPackage",
		//空数据包
		"data": "",
	}
	fmt.Println(util.Marshal(m))
	l, _ := util.Unmarshal(util.Marshal(m))
	for key := range l {
		fmt.Println(key, "首都是", l[key])
	}
	fmt.Println(l["type"])
}
