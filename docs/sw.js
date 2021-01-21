importScripts('https://cdn.jsdelivr.net/gh/nlepage/go-wasm-http-server@0bc6fbc8259cad3fa02812be7c11967397864944/sw.js')

addEventListener('install', (event) => {
  event.waitUntil(skipWaiting())
})
  
addEventListener('activate', event => {
  event.waitUntil(clients.claim())
})

registerWasmHTTPListener('api.wasm')
