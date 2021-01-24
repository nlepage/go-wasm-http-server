package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	wasmhttp "github.com/nlepage/go-wasm-http-server"
)

func main() {
	var no = 1

	http.HandleFunc("/hello", func(res http.ResponseWriter, req *http.Request) {
		params := make(map[string]string)
		if err := json.NewDecoder(req.Body).Decode(&params); err != nil {
			panic(err)
		}

		res.Header().Add("Content-Type", "application/json")
		if err := json.NewEncoder(res).Encode(map[string]string{
			"message": fmt.Sprintf("Hello %s! (request n°%d)", params["name"], no),
		}); err != nil {
			panic(err)
		}

		no++
	})

	wasmhttp.Serve(nil)

	select {}
}
