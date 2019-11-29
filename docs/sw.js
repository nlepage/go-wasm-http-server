importScripts('https://cdn.jsdelivr.net/gh/golang/go@go1.13.4/misc/wasm/wasm_exec.js')

addEventListener('install', (event) => {
  console.log('install!')
  // wasmhttp.serve({
  //   wasm: 'api.wasm',
  //   base: '/api',
  // })
  event.waitUntil(skipWaiting())
})

addEventListener('activate', event => {
  console.log('activate!')
  event.waitUntil(clients.claim())
})

addEventListener('fetch', () => {
  console.log('fetch!')
})

addEventListener('message', ({ data }) => {
  console.log('message', data)
})
