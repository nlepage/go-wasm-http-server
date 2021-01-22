importScripts('https://cdn.jsdelivr.net/gh/golang/go@go1.15.7/misc/wasm/wasm_exec.js')

function registerWasmHTTPListener(wasm, { base, args = [], timeout = 25 } = {}) {
  let path = new URL(registration.scope).pathname
  if (base && base !== '') path = `${trimEnd(path, '/')}/${trimStart(base, '/')}`

  let timeoutPromise
  let resetTimeout
  if (timeout) {
    let resolveTimeout
    timeoutPromise = new Promise(resolve => { resolveTimeout = resolve })
    let timeoutId
    resetTimeout = () => {
      clearTimeout(timeoutId)
      timeoutId = setTimeout(resolveTimeout, timeout * 1000)
    }
    resetTimeout()
  }

  const handlerPromise = new Promise(setHandler => {
    self.wasmhttp = {
      path,
      setHandler,
      timeoutPromise,
    }
  })

  const go = new Go()
  go.argv = [wasm, ...args]
  WebAssembly.instantiateStreaming(fetch(wasm), go.importObject).then(({ instance }) => go.run(instance))

  addEventListener('fetch', e => {
    resetTimeout?.()

    const { pathname } = new URL(e.request.url)
    if (!pathname.startsWith(path)) return

    e.respondWith(handlerPromise.then(handler => handler(e.request)))
  })
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
