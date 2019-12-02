window.wasmhttp = {
  register: async (wasm, { scope, base = '', swUrl = 'sw.js' } = {}) => {
    const options = {}
    if (scope) options.scope = scope
    //FIXME register once
    const registration = await navigator.serviceWorker.register(swUrl, options)
    await navigator.serviceWorker.ready
    registration.active.postMessage({
      type: 'wasmhttp.register',
      wasm,
      base,
    })
  }
}
