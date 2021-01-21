importScripts('https://cdn.jsdelivr.net/gh/nlepage/go-wasm-http-server@cea18c5f76930149be3dd059e2dde795bbdeb2da/sw.js')

addEventListener('install', (event) => {
  event.waitUntil(skipWaiting())
})
  
addEventListener('activate', event => {
  event.waitUntil(clients.claim())
})

registerWasmHTTPListener('api.wasm', 'api')
