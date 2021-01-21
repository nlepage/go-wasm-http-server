importScripts('https://cdn.jsdelivr.net/gh/golang/go@go1.15.7/misc/wasm/wasm_exec.js')

let nextHandlerId = 1
const handlerResolvers = {}
const handlers = []

self.wasmhttp = {
  registerHandler: (handlerId, handler) => {
    handlerResolvers[handlerId](handler)
    delete handlerResolvers[handlerId]
  },
}

function registerWasmHTTPListener(wasm, base, args) {
  let path = new URL(registration.scope).pathname
  if (base && base !== '') path = `${trimEnd(path, '/')}/${trimStart(base, '/')}`

  addEventListener('fetch', e => {
    const { pathname } = new URL(e.request.url)
    if (!pathname.startsWith(path)) return

    console.log(`FetchEvent ${pathname}`)

    const handlerId = `${nextHandlerId++}`
    const handlerPromise = new Promise(resolve => handlerResolvers[handlerId] = resolve)
  
    // FIXME await ? catch ?
    startWasm(wasm, { env: { WASMHTTP_HANDLER_ID: handlerId, WASMHTTP_PATH: path }, args })
  
    e.respondWith(handlerPromise.then(handler => handler(e.request)))
  })
}

async function startWasm(wasm, { env, args = [] }) {
  const go = new Go()
  go.env = env
  go.argv = [wasm, ...args]
  const { instance } = await WebAssembly.instantiateStreaming(fetch(wasm), go.importObject)
  return go.run(instance)
}

function trimStart(s, c) {
  let r = s
  while (r.startsWith(c)) r = r.slice(c.length)
  return r
}

function trimEnd(s, c) {
  let r = s
  while (r.endsWith(c)) r = r.slice(0, -c.length)
  return r
}
