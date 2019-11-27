importScripts(
  'https://cdn.jsdelivr.net/gh/nlepage/go-wasm-http-server@96677e251874f074906c61ccebd283c63cdec54d/lib/wasm_exec/wasm_exec.js',
  'https://cdn.jsdelivr.net/gh/nlepage/go-wasm-http-server@96677e251874f074906c61ccebd283c63cdec54d/index.js',
)

addEventListener('install', () => {
  wasmhttp.serve({
    wasm: 'api.wasm',
    base: '/api',
  })
  skipWaiting()
})

addEventListener('activate', event => event.waitUntil(clients.claim()))
