importScripts('https://cdn.jsdelivr.net/gh/golang/go@go1.23.1/misc/wasm/wasm_exec.js')
importScripts('../sw.js')

const wasm = '../hello-state/api.wasm'

addEventListener('install', event => {
  event.waitUntil(caches.open('hello-state').then((cache) => cache.add(wasm)))
})
  
addEventListener('activate', event => {
  event.waitUntil(clients.claim())
})

addEventListener('message', () => {})

registerWasmHTTPListener(wasm, { base: 'api' })
