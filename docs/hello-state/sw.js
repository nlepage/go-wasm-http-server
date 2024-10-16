importScripts('https://cdn.jsdelivr.net/gh/golang/go@go1.23.1/misc/wasm/wasm_exec.js')
importScripts('https://cdn.jsdelivr.net/gh/nlepage/go-wasm-http-server@v2.0.2/sw.js')

addEventListener('install', (event) => {
  event.waitUntil(caches.open('hello-state').then((cache) => cache.add('api.wasm')))
  event.waitUntil(skipWaiting())
})
  
addEventListener('activate', (event) => {
  event.waitUntil(clients.claim())
})

registerWasmHTTPListener('api.wasm', { base: 'api' })
