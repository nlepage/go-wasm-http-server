package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"

	wasmhttp "github.com/nlepage/go-wasm-http-server"
)

func main() {
	var counter int32

	http.HandleFunc("/hello", func(res http.ResponseWriter, req *http.Request) {
		params := make(map[string]string)
		if err := json.NewDecoder(req.Body).Decode(&params); err != nil {
			panic(err)
		}

		res.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(res).Encode(map[string]string{
			"message": fmt.Sprintf("Hello %s! (request %d)", params["name"], atomic.AddInt32(&counter, 1)),
		}); err != nil {
			panic(err)
		}
	})

	wasmhttp.Serve(nil)

	select {}
}
