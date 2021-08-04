package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

func main() {
	byte, err := os.ReadFile("./user.yaml")
	if err != nil {
		fmt.Printf("打开用户信息失败%s", err.Error())
	}
	var user []User
	yaml.Unmarshal(byte, &user)
	for item := range user {
		LoginEuserv(user[item])
	}
	fmt.Printf("\r\n所有账号全部续期完成")
}
