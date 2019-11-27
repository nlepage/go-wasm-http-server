importScripts(
  'https://raw.githubusercontent.com/nlepage/go-wasm-http-server/master/lib/wasm_exec/wasm_exec.js?sanitize=true',
  'https://raw.githubusercontent.com/nlepage/go-wasm-http-server/master/index.js?sanitize=true',
)

wasmhttp.serve({
  wasm: 'test.wasm',
  base: '/test',
})
