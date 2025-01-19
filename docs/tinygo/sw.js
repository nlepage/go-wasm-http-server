importScripts('https://cdn.jsdelivr.net/gh/tinygo-org/tinygo@0.35.0/targets/wasm_exec.js')
importScripts('https://cdn.jsdelivr.net/gh/nlepage/go-wasm-http-server@v2.0.5/sw.js')

const wasm = 'api.wasm'

addEventListener('install', (event) => {
  event.waitUntil(caches.open('examples').then((cache) => cache.add(wasm)))
})

addEventListener('activate', (event) => {
  event.waitUntil(clients.claim())
})

registerWasmHTTPListener(wasm, { base: 'api' })
