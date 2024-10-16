function registerWasmHTTPListener(wasm, { base, cacheName, args = [] } = {}) {
  let path = new URL(registration.scope).pathname
  if (base && base !== '') path = `${trimEnd(path, '/')}/${trimStart(base, '/')}`

  const handlerPromise = new Promise(setHandler => {
    self.wasmhttp = {
      path,
      setHandler,
    }
  })

  const go = new Go()
  go.argv = [wasm, ...args]
  const source = cacheName
    ? caches.open(cacheName).then((cache) => cache.match(wasm)).then((response) => response ?? fetch(wasm))
    : caches.match(wasm).then(response => (response) ?? fetch(wasm))
  WebAssembly.instantiateStreaming(source, go.importObject).then(({ instance }) => go.run(instance))

  addEventListener('fetch', e => {
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
