importScripts('https://cdn.jsdelivr.net/gh/nlepage/go-wasm-http-server@b17438900520578427d51505142202a9cf2997d6/sw.js')

addEventListener('install', (event) => {
  event.waitUntil(skipWaiting())
})
  
addEventListener('activate', event => {
  event.waitUntil(clients.claim())
})

registerWasmHTTPListener('api.wasm', 'api')
