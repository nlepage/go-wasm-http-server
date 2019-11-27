importScripts(
  'https://cdn.jsdelivr.net/gh/nlepage/go-wasm-http-server@96677e251874f074906c61ccebd283c63cdec54d/lib/wasm_exec/wasm_exec.js',
  'https://cdn.jsdelivr.net/gh/nlepage/go-wasm-http-server@96677e251874f074906c61ccebd283c63cdec54d/index.js',
)

console.log('scope', registration.scope)

addEventListener('install', event => {
  console.log('install!')
  wasmhttp.serve({
    wasm: 'api.wasm',
    base: '/api',
  })
  skipWaiting()
})

addEventListener('activate', event => {
  console.log('activate!')
  event.waitUntil(clients.claim())
})
