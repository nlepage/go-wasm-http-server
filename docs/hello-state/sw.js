importScripts('https://cdn.jsdelivr.net/gh/nlepage/go-wasm-http-server@v1.0.0/sw.js')

addEventListener('install', (event) => {
  event.waitUntil(skipWaiting())
})
  
addEventListener('activate', event => {
  event.waitUntil(clients.claim())
})

registerWasmHTTPListener('api.wasm', { base: 'api' })
