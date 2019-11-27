let nextHandlerId = 1
const handlerResolvers = {}

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

self.wasmhttp = {
  serve: async ({ wasm, base } = {}) => {
    try {
      if (!wasm) throw TypeError('options.wasm must be defined')

      const handlerId = `${nextHandlerId++}`
      const handler = new Promise(resolve => handlerResolvers[handlerId] = resolve)

      let path = new URL(registration.scope).pathname
      if (base && base !== '') path = `${trimEnd(path, '/')}/${trimStart(base, '/')}`

      startWasm(wasm, handlerId, path)

      addEventListener('fetch', async e => {
        if (!new URL(e.request.url).pathname.startsWith(path)) return

        // FIXME try catch
        e.respondWith((await handler)(e.request))
      })
    } catch (e) {
      console.error('wasmhttp: error:', e)
      throw e
    }
  },

  registerHandler: (handlerId, handler) => {
    handlerResolvers[handlerId](handler)
    delete handlerResolvers[handlerId]
  },
}
