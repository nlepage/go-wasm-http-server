importScripts('https://cdn.jsdelivr.net/gh/golang/go@go1.25.1/lib/wasm/wasm_exec.js');
importScripts('https://cdn.jsdelivr.net/gh/nlepage/go-wasm-http-server@master/sw.js');

addEventListener('install', (event) => {
  event.waitUntil(caches.open('examples').then((cache) => cache.addAll(['api1.wasm', 'api2.wasm'])));
});

addEventListener('activate', (event) => {
  event.waitUntil(clients.claim());
});

registerWasmHTTPListener('api1.wasm', { base: 'api1/1/' });
registerWasmHTTPListener('api1.wasm', { base: 'api1/2/' });
registerWasmHTTPListener('api2.wasm', { base: 'api2/1/' });
registerWasmHTTPListener('api2.wasm', { base: 'api2/2/' });
