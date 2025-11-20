importScripts('https://cdn.jsdelivr.net/gh/golang/go@go1.25.1/lib/wasm/wasm_exec.js');
importScripts('https://cdn.jsdelivr.net/gh/nlepage/go-wasm-http-server@master/sw.js');

const wasm = 'api.wasm';

addEventListener('install', (event) => {
  event.waitUntil(caches.open('examples').then((cache) => cache.add(wasm)));
});

addEventListener('activate', (event) => {
  event.waitUntil(clients.claim());
});

registerWasmHTTPListener(wasm, { base: 'api' });
