importScripts('https://cdn.jsdelivr.net/gh/nlepage/go-wasm-http-server@48880aadf0163ec0d688dc41f121bd143ab6c305/sw.js')

addEventListener('install', (event) => {
  event.waitUntil(skipWaiting())
})
  
addEventListener('activate', event => {
  event.waitUntil(clients.claim())
})

registerWasmHTTPListener('api.wasm')
