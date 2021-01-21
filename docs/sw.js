importScripts('https://cdn.jsdelivr.net/gh/nlepage/go-wasm-http-server@211427cb34ab93d81c79778674951746f7ab636f/sw.js')

addEventListener('install', (event) => {
  event.waitUntil(skipWaiting())
})
  
addEventListener('activate', event => {
  event.waitUntil(clients.claim())
})

registerWasmHTTPListener('api.wasm', 'api')
