const swUrl = 'https://cdn.jsdelivr.net/gh/nlepage/go-wasm-http-server@b4be701615d284296a4a1aaa852db02f139969ee/index.js'

window.wasmhttp = {
  register: async (wasm, { scope, base = '' } = {}) => {
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
