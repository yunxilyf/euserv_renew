export class EUservAutoCheck {
    /**
     * 用户数据
     */
    public UserList?: User[];
    /**
     * 处理日志
     */
    public LogFunction?: Function;
    public close?: Function;
    /**
     * 初始化续费器
     * @param userList 用户数据
     */
    constructor(userList?: User[], log?: Function,close?:Function) {
        this.UserList = userList;
        this.LogFunction = log;
        this.close=close;
    }
    async LoginEuserv(user: User): Promise<string> {
        var self = this;
        return new Promise<string>((resolve, reject) => {
            fetch('https://support.euserv.com/?method=json', {
                method: "GET"
            }).then(res => res.json()).then(tokenInfo => {
                let userToken = tokenInfo.result.sess_id.value;
                fetch(`https://support.euserv.com/?subaction=login&method=json&sess_id=${userToken}&email=${user.email}&password=${user.passWord}`).then(res => res.json()).then(res => {
                    if (res.message != "success") {
                        reject("登录失败[" + JSON.stringify(user) + "]：" + JSON.stringify(res))
                    }
                    fetch(`https://support.euserv.com/?action=showorders&sess_id=${userToken}`).then(res => res.text()).then(res => {
                        let pattern: RegExp = /Customer ID: (\d+)/;
                        let customerId = pattern.exec(res);
                        if (customerId == null) {
                            if (!!self.LogFunction) {
                                self.LogFunction(user,`[${user.email}]登录失败`,false);
                            }
                            reject(JSON.stringify(user) + "登录失败账号或密码错误")
                        } else {
                            let serverListPattern = /class="td-z1-sp1-kc">(\d+)<\/td>/g;
                            let serverList = res.match(serverListPattern);
                            serverList?.forEach(item => {
                                let order = serverListPattern.exec(item) || '';
                                self.Renewal(user, userToken, order[1]).then(log => {
                                    if (!!self.LogFunction) {
                                        self.LogFunction(user,`[${user.email}]` + log,true);
                                    }
                                })


                            })
                        }
                    })

                })

            })



        });
    }
    /**
     * 续期
     * @param sess_id 
     * @param ord_no 
     * @returns 
     */
    Renewal(user: User, sess_id: string, ord_no: string): Promise<string> {
        return new Promise<string>((resolve, reject) => {
            fetch(`https://support.euserv.com/?subaction=choose_order&method=json&sess_id=${sess_id}&choose_order_subaction=show_contract_details&ord_no=${ord_no}`).then(res => {
                //查看订单输入密码确认
                fetch(`https://support.euserv.com/?subaction=kc2_security_password_get_token&method=json&sess_id=${sess_id}&prefix=kc2_customer_contract_details_extend_contract_&password=${user.passWord}`).then(res => res.json()).then(res => {
                    let token = res.result.token.value;
                    if (!!token) {
                        fetch(`https://support.euserv.com/?method=json&sess_id=${sess_id}&subaction=kc2_customer_contract_details_extend_contract_term&token=${token}&ord_id=${ord_no}`).then(res => res.json()).then(res => {
                            resolve(res.message);
                        })
                    }

                })
            })

        })

    }
    async AutoLoginAndRenewal():Promise<undefined> {
        return new Promise<undefined>(()=>{
            this.UserList?.forEach(async item=>{
                console.log("执行任务",item)
                this.LoginEuserv(item);
            });
        });
    }
}

export interface User {
    email: string,
    passWord: string
}

