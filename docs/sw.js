importScripts('https://cdn.jsdelivr.net/gh/nlepage/go-wasm-http-server@d1781bb6931c5acda7dd8ce91a253817c517bed8/sw.js')

addEventListener('install', (event) => {
  event.waitUntil(skipWaiting())
})
  
addEventListener('activate', event => {
  event.waitUntil(clients.claim())
})

registerWasmHTTPListener('api.wasm', { base: 'api' })
