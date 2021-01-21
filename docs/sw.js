importScripts('https://cdn.jsdelivr.net/gh/nlepage/go-wasm-http-server@0e553499b390775d1361a2d0f492fdcaa6973ae7/sw.js')

addEventListener('install', (event) => {
  event.waitUntil(skipWaiting())
})
  
addEventListener('activate', event => {
  event.waitUntil(clients.claim())
})

registerWasmHTTPListener('api.wasm', 'api')
