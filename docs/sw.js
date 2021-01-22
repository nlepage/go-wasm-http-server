importScripts('https://cdn.jsdelivr.net/gh/nlepage/go-wasm-http-server@4cf9aa09614470c46e9cd9be25e03da4f780238c/sw.js')

addEventListener('install', (event) => {
  event.waitUntil(skipWaiting())
})
  
addEventListener('activate', event => {
  event.waitUntil(clients.claim())
})

registerWasmHTTPListener('api.wasm', { base: 'api' })
