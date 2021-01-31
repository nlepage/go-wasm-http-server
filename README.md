<h1 align="center">Welcome to go-wasm-http-server 👋</h1>
<p>
  <a href="https://pkg.go.dev/github.com/nlepage/go-wasm-http-server#section-documentation">
    <img src="https://pkg.go.dev/badge/github.com/nlepage/go-wasm-http-server.svg" alt="Go Reference">
  </a>
  <a href="https://github.com/nlepage/go-wasm-http-server/blob/master/LICENSE" target="_blank">
    <img alt="License: Apache 2.0" src="https://img.shields.io/badge/License-Apache--2.0-yellow.svg" />
  </a>
  <a href="https://twitter.com/njblepage" target="_blank">
    <img alt="Twitter: njblepage" src="https://img.shields.io/twitter/follow/njblepage.svg?style=social" />
  </a>
</p>

> Build your Go HTTP Server to [WebAssembly](https://mdn.io/WebAssembly/) and embed it in a ServiceWorker!

## Examples

 - [Hello example](https://nlepage.github.io/go-wasm-http-server/hello) ([sources](https://github.com/nlepage/go-wasm-http-server/tree/master/docs/hello))
 - [Hello example with state](https://nlepage.github.io/go-wasm-http-server/hello-state) ([sources](https://github.com/nlepage/go-wasm-http-server/tree/master/docs/hello-state))
 - [Hello example with state and keepalive](https://nlepage.github.io/go-wasm-http-server/hello-state-keepalive) ([sources](https://github.com/nlepage/go-wasm-http-server/tree/master/docs/hello-state-keepalive))
 - [😺 Catption generator example](https://nlepage.github.io/catption/wasm) ([sources](https://github.com/nlepage/catption/tree/wasm))

## Requirements

`go-wasm-http-server` requires you to build your Go application to WebAssembly, so you need to make sure your code is compatible:
- no C bindings
- no System dependencies such as file system or network (database server for example)

## Usage

### Step 1: Build to `js/wasm`

In your Go code, replace [`http.ListenAndServe()`](https://pkg.go.dev/net/http#ListenAndServe) (or [`net.Listen()`](https://pkg.go.dev/net#Listen) + [`http.Serve()`](https://pkg.go.dev/net/http#Serve)) by [wasmhttp.Serve()](https://pkg.go.dev/github.com/nlepage/go-wasm-http-server#Serve):

📄 `server.go`
```go
// +build !js,!wasm

package main

import (
    "net/http"
)

func main() {
    // Define handlers...

    http.ListenAndServe(":8080", nil)
}
```

becomes:

📄 `server_js_wasm.go`
```go
// +build js,wasm

package main

import (
    wasmhttp "github.com/nlepage/go-wasm-http-server"
)

func main() {
    // Define handlers...

    wasmhttp.Serve(nil)
}
```

You may want to use build tags as shown above (or file name suffixes) in order to be able to build both to WebAssembly and other targets.

Then build your WebAssembly binary:

```sh
GOOS=js GOARCH=wasm go build -o server.wasm .
```

### Step 2: Create ServiceWorker file

Create a ServiceWorker file with the following code:

📄 `sw.js`
```js
importScripts('https://cdn.jsdelivr.net/gh/nlepage/go-wasm-http-server@v1.0.0/sw.js')

registerWasmHTTPListener('path/to/server.wasm')
```

By default the server will deploy at the ServiceWorker's scope root, check [`registerWasmHTTPListener()`'s API](https://github.com/nlepage/go-wasm-http-server#registerwasmhttplistener) for more information.

You may want to add these additional event listeners in your ServiceWorker:

```js
// Skip installed stage and jump to activating stage
addEventListener('install', (event) => {
  event.waitUntil(skipWaiting())
})

// Start controlling clients as soon as the SW is activated
addEventListener('activate', event => {
  event.waitUntil(clients.claim())
})
```

### Step 3: Register the ServiceWorker

In your web page(s), register the ServiceWorker:

```html
<script>
  // By default the ServiceWorker's scope will be "server/"
  navigator.serviceWorker.register('server/sw.js')
</script>
```

Now your web page(s) may start fetching from the server:

```js
// The server will receive a request for "/path/to/resource"
fetch('server/path/to/resource').then(res => {
  // use response...
})
```

## API

### Go API

See [pkg.go.dev/github.com/nlepage/go-wasm-http-server](https://pkg.go.dev/github.com/nlepage/go-wasm-http-server#section-documentation)

### [`registerWasmHTTPListener(wasmUrl, options)`](https://github.com/nlepage/go-wasm-http-server/blob/v1.0.0/sw.js#L3)

Instantiates and runs the WebAssembly module at `wasmUrl`, and registers a fetch listener forwarding requests to the WebAssembly module.

#### `wasmUrl`

URL string of the WebAssembly module, example: `"path/to/my-module.wasm"`.

#### `options`

An optional object containing:

- `base` (`string`): Base path of the server, relative to the ServiceWorker's scope.
- `args` (`string[]`): Arguments for the WebAssembly module.

## Why?

`go-wasm-http-server` can help you put up a demonstration for a project without actually running a Go HTTP server.

## How?

If you want to know how `go-wasm-http-server` works, I will be presenting the project at [the FOSDEM 2021 Go devroom](https://fosdem.org/2021/schedule/room/dgo/).

The slides are available [here](https://nlepage.github.io/go-wasm-http-talk/).

## Author

👤 **Nicolas Lepage**

* Website: https://nicolas.lepage.dev/
* Twitter: [@njblepage](https://twitter.com/njblepage)
* Github: [@nlepage](https://github.com/nlepage)

## 🤝 Contributing

Contributions, issues and feature requests are welcome!<br />Feel free to check [issues page](https://github.com/nlepage/go-wasm-http-server/issues).

## Show your support

Give a ⭐️ if this project helped you!

## 📝 License

Copyright © 2021 [Nicolas Lepage](https://github.com/nlepage).<br />
This project is [Apache 2.0](https://github.com/nlepage/go-wasm-http-server/blob/master/LICENSE) licensed.

***
_This README was generated with ❤️ by [readme-md-generator](https://github.com/kefranabg/readme-md-generator)_