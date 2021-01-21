importScripts('https://cdn.jsdelivr.net/gh/nlepage/go-wasm-http-server@305a51fe51f5110dc7fef94b9e1c2f95e0219811/sw.js')

addEventListener('install', (event) => {
  event.waitUntil(skipWaiting())
})
  
addEventListener('activate', event => {
  event.waitUntil(clients.claim())
})

registerWasmHTTPListener('api.wasm', 'api')
