package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"time"
)

func main() {

	var path string
	flag.StringVar(&path, "file", "", "账户配置")
	flag.Parse()
	if len(path)==0 {
		fmt.Printf("没有配置账户文件！")
		os.Exit(1)
	}
	byte, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("打开用户信息失败%s", err.Error())
	}
	fmt.Printf("==================%s================\r\n",time.Now())
	fmt.Printf("续期任务开始拉~")
	
	var user []User
	yaml.Unmarshal(byte, &user)
	for item := range user {
		LoginEuserv(user[item])
	}
	fmt.Printf("\r\n所有账号全部续期完成")
}
