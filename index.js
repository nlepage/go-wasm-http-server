window.wasmhttp = {
  register: async (wasm, { scope, base = '', swUrl = 'sw.js', args = [] } = {}) => {
    const options = {}
    if (scope) options.scope = scope
    //FIXME register once (beware of changing scope ?)
    const registration = await navigator.serviceWorker.register(swUrl, options)
    await navigator.serviceWorker.ready
    registration.active.postMessage({
      type: 'wasmhttp.register',
      wasm,
      base,
      args,
    })
  }
}
