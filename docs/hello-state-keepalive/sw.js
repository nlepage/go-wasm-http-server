importScripts('https://cdn.jsdelivr.net/gh/golang/go@go1.23.1/misc/wasm/wasm_exec.js')
importScripts('https://cdn.jsdelivr.net/gh/nlepage/go-wasm-http-server@v2.0.2/sw.js')

addEventListener('install', event => {
  event.waitUntil(skipWaiting())
})
  
addEventListener('activate', event => {
  event.waitUntil(clients.claim())
})

addEventListener('message', () => {})

registerWasmHTTPListener('../hello-state/api.wasm', { base: 'api' })
