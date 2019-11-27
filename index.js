let nextHandlerId = 1
const handlerResolvers = {}

const startWasm = async (wasm, WASMHTTP_HANDLER_ID, WASMHTTP_BASE) => {
  const go = new Go()
  go.env = {
    WASMHTTP_HANDLER_ID,
    WASMHTTP_BASE,
  }
  const { instance } = await WebAssembly.instantiateStreaming(fetch(wasm), go.importObject)
  return go.run(instance)
}

self.wasmhttp = {
  serve: async ({ wasm, base } = {
    base: '',
  }) => {
    try {
      if (!wasm) throw TypeError('option.wasm must be defined')

      const handlerId = `${nextHandlerId++}`
      const handler = new Promise(resolve => handlerResolvers[handlerId] = resolve)

      startWasm(wasm, handlerId, base)

      addEventListener('fetch', async e => e.respondWith((await handler)(e.request)))
    } catch (e) {
      console.error('wasmhttp: error:', e)
    }
  },

  registerHandler: (handlerId, handler) => {
    handlerResolvers[handlerId](handler)
    delete handlerResolvers[handlerId]
  },
}
