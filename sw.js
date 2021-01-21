importScripts('https://cdn.jsdelivr.net/gh/golang/go@go1.15.7/misc/wasm/wasm_exec.js')

const startWasm = async (wasm, { env, args = [] }) => {
  const go = new Go()
  go.env = env
  go.argv = [wasm, ...args]
  const { instance } = await WebAssembly.instantiateStreaming(fetch(wasm), go.importObject)
  return go.run(instance)
}

const trimStart = (s, c) => {
  let r = s
  while (r.startsWith(c)) r = r.slice(c.length)
  return r
}

const trimEnd = (s, c) => {
  let r = s
  while (r.endsWith(c)) r = r.slice(0, -c.length)
  return r
}

// addEventListener('install', (event) => {
//   event.waitUntil(skipWaiting())
// })

const running = new Set()
let nextHandlerId = 1
const handlerResolvers = {}
const handlers = []

self.wasmhttp = {
  registerHandler: (handlerId, handler) => {
    handlerResolvers[handlerId](handler)
    delete handlerResolvers[handlerId]
  },
}

console.log('aha!')

addEventListener('activate', event => {
  console.log('activate!')
  event.waitUntil(clients.claim())
})

addEventListener('message', async ({ data }) => {
  console.log('message!', data)

  if (data.type !== 'wasmhttp.register') return

  const { wasm, base, args } = data

  let path = new URL(registration.scope).pathname
  if (base && base !== '') path = `${trimEnd(path, '/')}/${trimStart(base, '/')}`

  const key = `${wasm}:${path}`

  if (!running.has(key)) {
    const handlerId = `${nextHandlerId++}`
    const handler = new Promise(resolve => handlerResolvers[handlerId] = resolve)

    startWasm(wasm, { env: { WASMHTTP_HANDLER_ID: handlerId, WASMHTTP_PATH: path }, args })
    running.add(key)

    // FIXME try catch
    handlers.push([path, await handler])
  }
})

addEventListener('fetch', e => {
  console.log('fetch')

  const { pathname } = new URL(e.request.url)
  const [, handler] = handlers.find(([path]) => pathname.startsWith(path)) || []
  if (!handler) return

  e.respondWith(handler(e.request))
})

function registerWasmHTTPListener(wasm, base, args) {
  let path = new URL(registration.scope).pathname
  if (base && base !== '') path = `${trimEnd(path, '/')}/${trimStart(base, '/')}`

  addEventListener('fetch', async e => {
    const { pathname } = new URL(e.request.url)
    if (!pathname.startsWith(path)) return

    const handlerId = `${nextHandlerId++}`
    const handlerPromise = new Promise(resolve => handlerResolvers[handlerId] = resolve)
  
    // FIXME await ? catch ?
    startWasm(wasm, { env: { WASMHTTP_HANDLER_ID: handlerId, WASMHTTP_PATH: path }, args })
  
    const handler = await handlerPromise
    
    e.respondWith(handler(e.request))
  })
}