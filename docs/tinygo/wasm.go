//go:build wasm
// +build wasm

package main

import (
	"net/http"

	wasmhttp "github.com/nlepage/go-wasm-http-server/v2"
)

func main() {
	http.HandleFunc("/tiny", goRuntimeHandler)

	wasmhttp.Serve(nil)

	select {}
}
