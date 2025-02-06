<h1 align="center">Welcome to go-wasm-http-server 👋</h1>
<p>
  <a href="https://pkg.go.dev/github.com/nlepage/go-wasm-http-server#section-documentation">
    <img src="https://pkg.go.dev/badge/github.com/nlepage/go-wasm-http-server.svg" alt="Go Reference">
  </a>
  <a href="https://github.com/nlepage/go-wasm-http-server/blob/master/LICENSE" target="_blank">
    <img alt="License: Apache 2.0" src="https://img.shields.io/badge/License-Apache--2.0-yellow.svg" />
  </a>
</p>

> Embed your Go HTTP handlers in a ServiceWorker (using [WebAssembly](https://mdn.io/WebAssembly/)) and emulate an HTTP server!

## Examples

 - [Hello example](https://nlepage.github.io/go-wasm-http-server/hello) ([sources](https://github.com/nlepage/go-wasm-http-server/tree/master/docs/hello))
 - [Hello example with state](https://nlepage.github.io/go-wasm-http-server/hello-state) ([sources](https://github.com/nlepage/go-wasm-http-server/tree/master/docs/hello-state))
 - [Hello example with state and keepalive](https://nlepage.github.io/go-wasm-http-server/hello-state-keepalive) ([sources](https://github.com/nlepage/go-wasm-http-server/tree/master/docs/hello-state-keepalive))
 - [Hello example with Server Sent Events](https://nlepage.github.io/go-wasm-http-server/hello-sse/) ([sources](https://nlepage.github.io/go-wasm-http-server/hello-sse/))
 - [😺 Catption generator example](https://nlepage.github.io/catption/wasm) ([sources](https://github.com/nlepage/catption/tree/wasm))
 - [Random password generator web server](https://nlepage.github.io/random-password-please/) ([sources](https://github.com/nlepage/random-password-please) forked from [jbarham/random-password-please](https://github.com/jbarham/random-password-please))
 - [Server fallbacks, and compiling with TinyGo](https://nlepage.github.io/go-wasm-http-server/tinygo/) (runs locally; see [sources & readme](https://github.com/nlepage/go-wasm-http-server/tree/master/docs/tinygo#readme) for how to run this example)


## How?

Talk given at the Go devroom of FOSDEM 2021 explaining how `go-wasm-http-server` works:

[![Deploy a Go HTTP server in your browser Youtube link](https://raw.githubusercontent.com/nlepage/go-wasm-http-talk/main/youtube.png)](https://youtu.be/O2RB_8ircdE)

The slides are available [here](https://nlepage.github.io/go-wasm-http-talk/).

## Why?

`go-wasm-http-server` can help you put up a demonstration for a project without actually running a Go HTTP server.

## Requirements

`go-wasm-http-server` requires you to build your Go application to WebAssembly, so you need to make sure your code is compatible:
- no C bindings
- no System dependencies such as file system or network (database server for example)
- For smaller WASM blobs, your code may also benefit from being compatible with, and compiled by, [TinyGo](https://tinygo.org/docs/reference/lang-support/stdlib/). See the TinyGo specific details below.

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
    wasmhttp "github.com/nlepage/go-wasm-http-server/v2"
)

func main() {
    // Define handlers...

    wasmhttp.Serve(nil)
}
```

You may want to use build tags as shown above (or file name suffixes) in order to be able to build both to WebAssembly and other targets.

Then build your WebAssembly binary:

```sh
# To compile with Go
GOOS=js GOARCH=wasm go build -o server.wasm .

# To compile with TinyGo, if your code is compatible
GOOS=js GOARCH=wasm tinygo build -o server.wasm  .
```

### Step 2: Create ServiceWorker file

First, check the version of Go/TinyGo you compiled your wasm with:

```sh
$ go version
go version go1.23.4 darwin/arm64
#          ^------^

$ tinygo version
tinygo version 0.35.0 darwin/arm64 (using go version go1.23.4 and LLVM version 18.1.2)
#              ^----^
```

Create a ServiceWorker file with the following code:

📄 `sw.js`
```js
// Note the 'go.1.23.4' below, that matches the version you just found:
importScripts('https://cdn.jsdelivr.net/gh/golang/go@go1.23.4/misc/wasm/wasm_exec.js')
// If you compiled with TinyGo then, similarly, use:
importScripts('https://cdn.jsdelivr.net/gh/tinygo-org/tinygo@0.35.0/targets/wasm_exec.js')

importScripts('https://cdn.jsdelivr.net/gh/nlepage/go-wasm-http-server@v2.1.0/sw.js')

registerWasmHTTPListener('path/to/server.wasm')
```

By default the server will deploy at the ServiceWorker's scope root, check [`registerWasmHTTPListener()`'s API](https://github.com/nlepage/go-wasm-http-server#registerwasmhttplistenerwasmurl-options) for more information.

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

For Go API see [pkg.go.dev/github.com/nlepage/go-wasm-http-server](https://pkg.go.dev/github.com/nlepage/go-wasm-http-server#section-documentation)

### JavaScript API

### `registerWasmHTTPListener(wasmUrl, options)`

Instantiates and runs the WebAssembly module at `wasmUrl`, and registers a fetch listener forwarding requests to the WebAssembly module's server.

⚠ This function must be called only once in a ServiceWorker, if you want to register several servers you must use several ServiceWorkers.

The server will be "deployed" at the root of the ServiceWorker's scope by default, `base` may be used to deploy the server at a subpath of the scope.

See [ServiceWorkerContainer.register()](https://developer.mozilla.org/en-US/docs/Web/API/ServiceWorkerContainer/register) for more information about the scope of a ServiceWorker.

#### `wasmUrl`

URL string of the WebAssembly module, example: `"path/to/my-module.wasm"`.

#### `options`

An optional object containing:

- `base` (`string`): Base path of the server, relative to the ServiceWorker's scope.
- `cacheName` (`string`): Name of the [Cache](https://developer.mozilla.org/en-US/docs/Web/API/Cache) to store the WebAssembly binary.
- `args` (`string[]`): Arguments for the WebAssembly module.
- `passthrough` (`(request: Request): boolean`): Optional callback to allow passing the request through to network.

## Contributors ✨

Thanks goes to these wonderful people ([emoji key](https://allcontributors.org/docs/en/emoji-key)):

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tbody>
    <tr>
      <td align="center" valign="top" width="14.28%"><a href="https://byjp.me/"><img src="https://avatars.githubusercontent.com/u/42999?v=4?s=100" width="100px;" alt="JP Hastings-Edrei"/><br /><sub><b>JP Hastings-Edrei</b></sub></a><br /><a href="https://github.com/nlepage/go-wasm-http-server/commits?author=jphastings" title="Code">💻</a> <a href="https://github.com/nlepage/go-wasm-http-server/commits?author=jphastings" title="Documentation">📖</a> <a href="#example-jphastings" title="Examples">💡</a></td>
      <td align="center" valign="top" width="14.28%"><a href="https://recolude.com/"><img src="https://avatars.githubusercontent.com/u/9094977?v=4?s=100" width="100px;" alt="Eli Davis"/><br /><sub><b>Eli Davis</b></sub></a><br /><a href="https://github.com/nlepage/go-wasm-http-server/commits?author=EliCDavis" title="Code">💻</a></td>
    </tr>
  </tbody>
</table>

<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->

<!-- ALL-CONTRIBUTORS-LIST:END -->

This project follows the [all-contributors](https://github.com/all-contributors/all-contributors) specification. Contributions of any kind welcome!

## 🤝 Contributing

Contributions, issues and feature requests are welcome!<br />Feel free to check [issues page](https://github.com/nlepage/go-wasm-http-server/issues).

## Show your support

Give a ⭐️ if this project helped you!

## 📝 License

Copyright © 2025 [Nicolas Lepage](https://github.com/nlepage).<br />
This project is [Apache 2.0](https://github.com/nlepage/go-wasm-http-server/blob/master/LICENSE) licensed.

***
_This README was generated with ❤️ by [readme-md-generator](https://github.com/kefranabg/readme-md-generator)_
