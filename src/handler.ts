import { EUservAutoCheck, User } from './EUservAutoCheck/EUservAutoCheck'
export async function handleRequest(request: Request): Promise<Response> {
  let User: User[] = JSON.parse(await Mydb.get('UserList') || "{}") as User[];
  let LogList:any = [];
  User.forEach(async item => {
    console.log(item.email+'Log')
    let log=await Mydb.get(item.email+'Log');   
    console.log("读取日志",log);
    LogList.push(log);
  })
  return new Response(`${JSON.stringify(LogList)}`)
}
/**
 * 获取队列中的用户数据 
 */
async function GetQueueUser() {
  //取出队列中的8个数据
  let queue: User[] = [];
  //全部队列
  let UserQueue: User[] = JSON.parse(await Mydb.get('queue') || "[]") as User[];
  console.log("全部队列数据", UserQueue);
  if (UserQueue.length == 0) {
    console.log("总任务队列为空")
    return undefined;
  }
  let getQueueLength = UserQueue.length > 8 ? 8 : UserQueue.length;
  for (var i = 0; i < getQueueLength; i++) {
    queue.push(UserQueue[i]);
    delete UserQueue[i]
  }
  if (queue.length == 0) {
    console.log("执行队列为空")
    return undefined;
  }
  console.log("处理完的队列数据", UserQueue);
  //保存回队列  
  await Mydb.put('queue', JSON.stringify(UserQueue.filter(item => !!item)))
  return queue;
}
/**
 * 续期
 * @param corn 定时表达式
 * @returns 
 */
export async function xuqi(corn: string) {
  if (corn == "59 23 * *  *" || !corn) {
    let User: User[] = JSON.parse(await Mydb.get('UserList') || "{}") as User[];
    await Mydb.put('queue', JSON.stringify(User))
  }
  let queue = await GetQueueUser();

  console.log("执行任务", queue)
  let auto = new EUservAutoCheck(queue, async function (user: User, log: string,flag?:boolean) {
    new Promise(async () => {
      console.log(log);
      let LogDB = JSON.parse(await Mydb.get(user.email + 'Log') || "{}");
      let time = new Date().toJSON();
      LogDB[time] = log;
      await Mydb.put(user.email + 'Log', JSON.stringify(LogDB));
    });
  });
  auto.AutoLoginAndRenewal();
}
