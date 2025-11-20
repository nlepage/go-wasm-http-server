package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	wasmhttp "github.com/nlepage/go-wasm-http-server/v2"
)

var binaryName = ""

func main() {
	http.HandleFunc("/hello", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(res).Encode(map[string]string{
			"message": fmt.Sprintf("Hello from %s at path %s", binaryName, os.Getenv("WASM_HTTP_PATH")),
		}); err != nil {
			panic(err)
		}
	})

	wasmhttp.Serve(nil)

	select {}
}
