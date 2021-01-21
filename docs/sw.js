importScripts('https://cdn.jsdelivr.net/gh/nlepage/go-wasm-http-server@cf2a3ceefa9e6180860ad449d5c53edf81dd42dc/sw.js')

addEventListener('install', (event) => {
  event.waitUntil(skipWaiting())
})
  
addEventListener('activate', event => {
  event.waitUntil(clients.claim())
})

registerWasmHTTPListener('api.wasm', 'api')
