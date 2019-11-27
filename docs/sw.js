importScripts(
  'https://raw.githubusercontent.com/nlepage/go-wasm-http-server/master/lib/wasm_exec/wasm_exec.js',
  'https://raw.githubusercontent.com/nlepage/go-wasm-http-server/master/index.js',
)

wasmhttp.serve({
  wasm: 'test.wasm',
  base: '/test',
})
