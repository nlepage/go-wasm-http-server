importScripts('https://cdn.jsdelivr.net/gh/golang/go@go1.13.4/misc/wasm/wasm_exec.js')

const startWasm = async (wasm, WASMHTTP_HANDLER_ID, WASMHTTP_PATH) => {
  const go = new Go()
  go.env = {
    WASMHTTP_HANDLER_ID,
    WASMHTTP_PATH,
  }
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

addEventListener('activate', event => event.waitUntil(clients.claim()))

addEventListener('message', async ({ data }) => {
  if (data.type !== 'wasmhttp.register') return

  const { wasm, base } = data

  let path = new URL(registration.scope).pathname
  if (base && base !== '') path = `${trimEnd(path, '/')}/${trimStart(base, '/')}`

  const key = `${wasm}:${path}`

  if (!running.has(key)) {
    const handlerId = `${nextHandlerId++}`
    const handler = new Promise(resolve => handlerResolvers[handlerId] = resolve)

    startWasm(wasm, handlerId, path)
    running.add(key)

  // FIXME try catch
    handlers.push([path, await handler])
  }
})

addEventListener('fetch', e => {
  const { pathname } = new URL(e.request.url)
  const [, handler] = handlers.find(([path]) => pathname.startsWith(path)) || []
  if (!handler) return

  e.respondWith((handler)(e.request))
})
