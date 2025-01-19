//go:build !wasm
// +build !wasm

package main

import (
	"embed"
	"fmt"
	"log"
	"net"
	"net/http"
)

//go:embed *.html *.js *.wasm
var thisDir embed.FS

func main() {
	// Serve all files in this directory statically
	http.Handle("/", http.FileServer(http.FS(thisDir)))

	// Note that this needs to be mounted at /api/tiny, rather than just /tiny (like in wasm.go)
	// because the service worker mounts the WASM server at /api (at the end of sw.js)
	http.HandleFunc("/api/tiny", goRuntimeHandler)

	// Pick any available port. Note that ServiceWorkers _require_ localhost for non-SSL serving (so other LAN/WAN IPs will prevent the service worker from loading)
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("Unable to claim a port to start server on: %v", err)
	}

	// Share the port being used & start
	fmt.Printf("Server starting on http://127.0.0.1:%d\n", listener.Addr().(*net.TCPAddr).Port)
	panic(http.Serve(listener, nil))
}
