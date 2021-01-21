importScripts('https://cdn.jsdelivr.net/gh/nlepage/go-wasm-http-server@9ee7a321b6d6c00036da1c6a47bd5c1b78c011e2/sw.js')

addEventListener('install', (event) => {
  event.waitUntil(skipWaiting())
})
  
addEventListener('activate', event => {
  event.waitUntil(clients.claim())
})

registerWasmHTTPListener('api.wasm')
