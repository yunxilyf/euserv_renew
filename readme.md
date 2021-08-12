# 这个是个基于GOLANG开发的Euserv自动续期脚本
· 怎么用？
首先在程序目录创建user.yaml 按照下面的格式填写 在执行./euserv -file ./user.yaml 就自动续费了
请配合cron使用
```text
- username: "账号"
  password: "密码"
- username: "账号"
  password: "密码"
```
# Power BY [小幻博客](https://52xhuan.cn)