self.wasmhttp = {
  Serve: async (wasm) => {
    const go = new Go()
    const { instance } = await WebAssembly.instantiateStreaming(fetch(wasm), go.importObject)
    try {
        await go.run(instance)
    } catch (e) {
        console.error(e)
    }

    addEventListener('fetch', async e => {
      if (new URL(e.request.url).pathname !== '/test') return
      e.respondWith((await fetchHandler)(e.request))
    })
  }
}
