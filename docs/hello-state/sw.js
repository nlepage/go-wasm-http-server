importScripts('https://cdn.jsdelivr.net/gh/nlepage/go-wasm-http-server@078ff3547ebe2abfbee1fd5af9ca5ad64be480c0/sw.js')

addEventListener('install', (event) => {
  event.waitUntil(skipWaiting())
})
  
addEventListener('activate', event => {
  event.waitUntil(clients.claim())
})

registerWasmHTTPListener('api.wasm', { base: 'api' })
