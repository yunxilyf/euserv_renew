import { handleRequest,xuqi } from './handler'

addEventListener('fetch', (event) => {
  event.respondWith(handleRequest(event.request))
})
addEventListener('scheduled', (event:any)=>{
  event.waitUntil(new Promise((resolve, reject) => {
    console.log(event)
     xuqi(event.cron);

  }))
})