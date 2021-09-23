package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/smtp"
)

func SendMail(user string ,pwd string,server string,Subject string,body string,toMail string)  {
	host,port,_:=net.SplitHostPort(server)

	auth := smtp.PlainAuth("", user, pwd, host)
	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	to := []string{toMail}
	msg := []byte("To: "+toMail+"\r\n" +
		"Subject:"+Subject+" \r\n" +
		"\r\n" +body+ ".\r\n")
	err := smtp.SendMail(host+":"+port, auth, user, to, msg)
	if err != nil {
		log.Fatal(err)
	}
}

func SendMailTls(user string ,pwd string,server string,Subject string,body string,toMail string){
	host,_,_:=net.SplitHostPort(server)
	header := make(map[string]string)
	header["From"] = user
	header["To"] = toMail
	header["Subject"] = Subject
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body
	auth := smtp.PlainAuth(
		"",
		user,
		pwd,
		host,
	)
	SendMailUsingTLS(server,auth,user,[]string{toMail},[]byte(message))
}
func Dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		log.Println("Dialing Error:", err)
		return nil, err
	}
	//分解主机端口字符串
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}
func SendMailUsingTLS(addr string, auth smtp.Auth, from string,
	to []string, msg []byte) (err error) {
	//create smtp client
	c, err := Dial(addr)
	if err != nil {
		log.Println("Create smpt client error:", err)
		return err
	}
	defer c.Close()
	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				log.Println("Error during AUTH", err)
				return err
			}
		}
	}
	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}