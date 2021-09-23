# 这个是个基于GOLANG开发的Euserv自动续期脚本
· 怎么用？
首先在程序目录创建user.yaml 按照下面的格式填写 在执行./euserv -file ./user.yaml 就自动续费了
请配合cron使用
```yaml
#smtp服务器
smtpserver: "host:port"
#是否启用ssl加密
smtpssl: false
#smtp用户名
smtpusername: ""
#smtp密码
smtppassword: ""
#邮箱接收人
contactsmail: ""
#账号列表
accounts:
  - username: "用户1"
    password: "用户1"
  - username: "用户2"
    password: "用户2"
```
#打赏一杯快乐水
![打赏](https://52xhuan.cn/wp-content/uploads/2021/08/zanshang.png)
# Power BY [小幻博客](https://52xhuan.cn)