package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"time"
)

type User struct {
	UserName string
	PassWord string
}
//配置
type Config struct {
	SmtpServer string
	SmtpSSL 	bool
	SmtpUserName string
	SmtpPassWord string
	ContactsMail string
	Accounts []User  `yaml:"accounts"`
}
//第一步获取token的模型
type UserLoginModel struct {
	Message  string `json:"message"`
	Code     string `json:"code"`
	TimeReal string `json:"time_real"`
	TimeUser string `json:"time_user"`
	TimeSys  string `json:"time_sys"`
	Result   struct {
		SessId struct {
			Value string `json:"value"`
		} `json:"sess_id"`
	} `json:"result"`
}

//2.登录的模型
type LoginModel struct {
	Message  string `json:"message"`
	Code     string `json:"code"`
	TimeReal string `json:"time_real"`
	TimeUser string `json:"time_user"`
	TimeSys  string `json:"time_sys"`
	Result   struct {
		SessId struct {
			Value string `json:"value"`
		} `json:"sess_id"`
	} `json:"result"`
}

//3选中订单
type selectOrder struct {
	Message  string        `json:"message"`
	Code     string        `json:"code"`
	TimeReal string        `json:"time_real"`
	TimeUser string        `json:"time_user"`
	TimeSys  string        `json:"time_sys"`
	Result   []interface{} `json:"result"`
}

//4输入密码确认订单
type checkOrder struct {
	Message  string `json:"message"`
	Code     string `json:"code"`
	TimeReal string `json:"time_real"`
	TimeUser string `json:"time_user"`
	TimeSys  string `json:"time_sys"`
	Result   struct {
		Token struct {
			Value string `json:"value"`
		} `json:"token"`
		SessId struct {
			Value string `json:"value"`
		} `json:"sess_id"`
	} `json:"result"`
}

//续期模型
type xuqiModel struct {
	Message  string `json:"message"`
	Code     string `json:"code"`
	TimeReal string `json:"time_real"`
	TimeUser string `json:"time_user"`
	TimeSys  string `json:"time_sys"`
	Result   struct {
		SessId struct {
			Value string `json:"value"`
		} `json:"sess_id"`
	} `json:"result"`
}

var client = http.Client{}
var useTime int
var Log string
func LoginEuserv(user User) {
	startTime := time.Now()
	Log+=fmt.Sprintf("========================%s=============================\r\n",user.UserName)
	fmt.Printf("\r\n==========%s==========\r\n", user.UserName)
	request, _ := http.NewRequest("GET", "https://support.euserv.com/?method=json", nil)
	res, err := client.Do(request)
	if err != nil {
		fmt.Printf("登录%s账号失败！%s\r\n", user.UserName, err.Error())
		Log+=fmt.Sprintf("登录%s账号失败！%s\r\n", user.UserName, err.Error())
		return
	}
	defer res.Body.Close()
	//判断获取tokne是否成功
	if res.StatusCode == 200 {
		tokenInfo, _ := io.ReadAll(res.Body)
		var TokenInfo UserLoginModel
		tokenErr := json.Unmarshal(tokenInfo, &TokenInfo)
		if tokenErr != nil {
			fmt.Printf("解析登录Token错误\r\n")
			Log+=fmt.Sprintf("解析登录Token错误\r\n")
			return
		}
		url := fmt.Sprintf("https://support.euserv.com/?subaction=login&method=json&sess_id=%s&email=%s&password=%s", TokenInfo.Result.SessId.Value, user.UserName, user.PassWord)
		//fmt.Printf("%s\r\n", url)
		loginRequest, err := http.NewRequest("GET", url, nil)
		loginRes, err := client.Do(loginRequest)
		if err != nil {
			fmt.Printf("登录Euser获取响应失败!\r\n")
			Log+=fmt.Sprintf("登录Euser获取响应失败!\r\n")
			return
		}
		if loginRes.StatusCode == 200 {
			loginInfoByte, _ := io.ReadAll(loginRes.Body)
			var loginInfo LoginModel
			json.Unmarshal(loginInfoByte, &loginInfo)
			if loginInfo.Code != "100" {
				fmt.Printf("登录账号失败%s", loginInfo.Message)
				Log+=fmt.Sprintf("登录账号失败%s", loginInfo.Message)
				return
			}
			//订单的url
			orderUrl := fmt.Sprintf("https://support.euserv.com/?action=showorders&sess_id=%s", loginInfo.Result.SessId.Value)
			orderRequest, _ := http.NewRequest("GET", orderUrl, nil)
			orderRes, err := client.Do(orderRequest)
			if err != nil {
				fmt.Printf("获取订单失败%s", err.Error())
				Log+=fmt.Sprintf("获取订单失败%s", err.Error())
				return
			}
			if orderRes.StatusCode == 200 {
				fmt.Printf("正在获取订单...\r\n")
				orderDoc, err := goquery.NewDocumentFromReader(orderRes.Body)
				if err != nil {
					fmt.Printf("解析订单html失败")
					Log+=fmt.Sprintf("解析订单html失败")
					return
				}
				orderDoc.Find(".kc2_order_table tr").Each(func(i int, s *goquery.Selection) {
					if i > 0 {
						orderId := s.Find(".td-z1-sp1-kc").Text()
						fmt.Printf("正在续费订单：%s\r\n", orderId)
						Log+=fmt.Sprintf("正在续费订单：%s\r\n", orderId)
						Renew(loginInfo.Result.SessId.Value, orderId, user)
					}
				})
				fmt.Printf("续费完成总共用时：%v", time.Since(startTime))
				Log+=fmt.Sprintf("续费完成总共用时：%v", time.Since(startTime))
			} else {
				fmt.Printf("获取订单失败\r\n")
				Log+=fmt.Sprintf("获取订单失败\r\n")
				return
			}
		} else {
			fmt.Printf("登录Euser登录失败!")
			Log+=fmt.Sprintf("登录Euser登录失败!")
		}

	} else {
		fmt.Printf("获取Token失败")
		Log+=fmt.Sprintf("获取Token失败")
	}

	//defer res.Body.Close()
}
func Renew(token string, order string, user User) {
	selectOrderUrl := fmt.Sprintf("https://support.euserv.com/?subaction=choose_order&method=json&sess_id=%s&choose_order_subaction=show_contract_details&ord_no=%s", token, order)
	//fmt.Printf("%s\r\n", selectOrderUrl)
	selectOrderReq, _ := http.NewRequest("GET", selectOrderUrl, nil)
	selectOrderRes, err := client.Do(selectOrderReq)
	if err != nil {
		fmt.Printf("选择订单(%s)失败%s\r\n", order, err.Error())
		Log+=fmt.Sprintf("选择订单(%s)失败%s\r\n", order, err.Error())
		return
	}
	if selectOrderRes.StatusCode == 200 {
		var selOrder selectOrder
		selecrOrderByte, _ := io.ReadAll(selectOrderRes.Body)
		json.Unmarshal(selecrOrderByte, &selOrder)
		checkOrderURL := fmt.Sprintf("https://support.euserv.com/?subaction=kc2_security_password_get_token&method=json&sess_id=%s&prefix=kc2_customer_contract_details_extend_contract_&password=%s", token, user.PassWord)
		//fmt.Printf(checkOrderURL)
		checkOrderUrlReq, _ := http.NewRequest("GET", checkOrderURL, nil)
		checkOrderURLRes, err := client.Do(checkOrderUrlReq)
		if err != nil {
			fmt.Printf("确认订单(%s)错误:%s", order, err.Error())
			Log+=fmt.Sprintf("确认订单(%s)错误:%s", order, err.Error())
			return
		}
		var checdOrder checkOrder
		orderByte, _ := io.ReadAll(checkOrderURLRes.Body)
		json.Unmarshal(orderByte, &checdOrder)
		if checdOrder.Code != "100" {
			fmt.Printf("确认订单(%s)错误:%s", order, checdOrder.Message)
			Log+=fmt.Sprintf("确认订单(%s)错误:%s", order, checdOrder.Message)
			return
		}
		xufeiUrl := fmt.Sprintf("https://support.euserv.com/?method=json&sess_id=%s&subaction=kc2_customer_contract_details_extend_contract_term&token=%s&ord_id=%s", token, checdOrder.Result.Token.Value, order)
		//fmt.Printf("%s\r\n", xufeiUrl)
		xufeiReq, _ := http.NewRequest("GET", xufeiUrl, nil)
		xufeiRes, err := client.Do(xufeiReq)
		if err != nil {
			fmt.Printf("确认订单(%s)错误:%s", order, err.Error())
			Log+=fmt.Sprintf("确认订单(%s)错误:%s", order, err.Error())
			return
		}
		if xufeiRes.StatusCode == 200 {
			xufeiByte, _ := io.ReadAll(xufeiRes.Body)
			var xuqi xuqiModel
			json.Unmarshal(xufeiByte, &xuqi)
			fmt.Printf("订单（%s）续期完毕结果：%s", order, xuqi.Message)
			Log+=fmt.Sprintf("订单（%s）续期完毕结果：%s", order, xuqi.Message)
		} else {
			fmt.Printf("续费确认订单(%s)错误:%s", order, xufeiRes.Status)
			Log+=fmt.Sprintf("续费确认订单(%s)错误:%s", order, xufeiRes.Status)
			return
		}
	} else {
		fmt.Printf("选择订单(%s)错误:%s", order, err.Error())
		Log+=fmt.Sprintf("选择订单(%s)错误:%s", order, err.Error())
		return

	}
}
