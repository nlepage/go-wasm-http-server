window.wasmhttp = {
  register: async (wasm, { scope, base = '' } = {}) => {
    const options = {}
    if (scope) options.scope = scope
    //FIXME register once
    const registration = await navigator.serviceWorker.register('sw.js', options)
    await navigator.serviceWorker.ready
    registration.active.postMessage({
      type: 'wasmhttp.register',
      wasm,
      base,
    })
  }
}
