importScripts(
  'https://cdn.jsdelivr.net/gh/nlepage/go-wasm-http-server@latest/lib/wasm_exec/wasm_exec.js',
  'https://cdn.jsdelivr.net/gh/nlepage/go-wasm-http-server@latest/index.js',
)

wasmhttp.serve({
  wasm: 'api.wasm',
  base: '/api',
})
