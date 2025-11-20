class WasmHTTP {
  #listeners = new Map();

  setHandler(path, handler) {
    const listener = this.#listeners.get(path);
    if (!listener) {
      throw new Error(`no listener for path "${path}"`);
    }

    listener.setHandler(handler);
  }

  handle(e) {
    const { pathname } = new URL(e.request.url);

    for (const [path, listener] of this.#listeners) {
      if (pathname.startsWith(path)) {
        listener.handle(e);
        return;
      }
    }
  }

  addListener({ wasm, base = '', cacheName, passthrough, args = [], env = {} }) {
    const path = new URL(trimStart(base, '/'), registration.scope).pathname;

    if (this.#listeners.has(path)) {
      throw new Error(`a listener is already registered for path "${path}"`);
    }

    this.#listeners.set(path, new WasmHTTPListener({ wasm, path, cacheName, passthrough, args, env }));
  }
}

class WasmHTTPListener {
  #handlerPromise;
  #resolveHandler;
  #passthrough;

  constructor({ wasm, path, cacheName, passthrough, args, env }) {
    this.#handlerPromise = new Promise((resolve) => {
      this.#resolveHandler = resolve;
    });

    this.#passthrough = passthrough;

    this.#run({ wasm, path, cacheName, passthrough, args, env });
  }

  async #run({ wasm, path, cacheName, args, env }) {
    try {
      const go = new Go();
      go.argv = [wasm, ...args];
      go.env = { ...env, WASM_HTTP_PATH: path };

      const cache = cacheName ? await caches.open(cacheName) : caches;
      const source = (await cache.match(wasm)) ?? (await fetch(wasm));

      const { instance } = await WebAssembly.instantiateStreaming(source, go.importObject);

      await go.run(instance);
    } catch (err) {
      console.error(`error while running ${wasm} for path "${path}"`, err);
    }
  }

  setHandler(handler) {
    this.#resolveHandler(handler);
  }

  handle(e) {
    if (this.#passthrough?.(e.request)) return;

    // FIXME return 500 if run has thrown

    e.respondWith(this.#handlerPromise.then((handler) => handler(e.request)));
  }
}

self.wasmhttp = new WasmHTTP();

addEventListener('fetch', (e) => {
  self.wasmhttp.handle(e);
});

function registerWasmHTTPListener(wasm, { base, cacheName, passthrough, args, env } = {}) {
  self.wasmhttp.addListener({ wasm, base, cacheName, passthrough, args, env });
}

function trimStart(s, c) {
  let r = s;
  while (r.startsWith(c)) r = r.slice(c.length);
  return r;
}
