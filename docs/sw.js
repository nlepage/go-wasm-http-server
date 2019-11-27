importScripts(
  'https://cdn.jsdelivr.net/gh/nlepage/go-wasm-http-server@master/lib/wasm_exec/wasm_exec.js',
  'https://cdn.jsdelivr.net/gh/nlepage/go-wasm-http-server@master/index.js',
)

wasmhttp.serve({
  wasm: 'test.wasm',
  base: '/test',
})
