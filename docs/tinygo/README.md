# Compiling with TinyGo

This example demonstrates that go-wasm-http-server can also be compiled with [TinyGo](https://www.tinygo.org), producing significantly smaller WASM blobs, though at the expense of [at least one known bug](https://github.com/tinygo-org/tinygo/issues/1140) and a [reduced standard library](https://tinygo.org/docs/reference/lang-support/stdlib/).

This example also demonstrates how the same code can be used for both server-side execution, and client-side execution in WASM (providing support for clients that cannot interpret WASM).

## Prerequisites

You'll need a version of [TinyGo installed](https://tinygo.org/getting-started/install/).  (eg. `brew install tinygo-org/tools/tinygo`)

You'll need to make sure the first line of `sw.js` here has the same tinygo version number as your TinyGo version (this was v0.35.0 at time of writing).

## Build & run

Compile the WASM blob with TinyGo (this has been done for you for this example):

```bash
GOOS=js GOARCH=wasm tinygo build -o api.wasm  .
```

Run the server (with Go, not TinyGo):

```bash
$ go run .
Server starting on http://127.0.0.1:<port>
```

## Important notes

You **must** use the TinyGo `wasm_exec.js`, specific to the version of TinyGo used to compile the WASM, in your `sw.js`. For example, if using the JSDelivr CDN:

```js
importScripts('https://cdn.jsdelivr.net/gh/tinygo-org/tinygo@0.35.0/targets/wasm_exec.js')
```

Note that the `0.35.0` within the path matches the TinyGo version used:

```sh
$ tinygo version
tinygo version 0.35.0 darwin/arm64 (using go version go1.23.4 and LLVM version 18.1.2)
#              ^----^
```
