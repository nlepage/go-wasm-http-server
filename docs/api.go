package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	wasmhttp "github.com/nlepage/go-wasm-http-server"
)

func main() {
	defer fmt.Println("api stopping...")

	http.HandleFunc("/hello", func(res http.ResponseWriter, req *http.Request) {
		params := make(map[string]string)
		if err := json.NewDecoder(req.Body).Decode(&params); err != nil {
			panic(err)
		}

		if err := json.NewEncoder(res).Encode(map[string]string{
			"message": fmt.Sprintf("Hello %s!", params["name"]),
		}); err != nil {
			panic(err)
		}
	})

	wasmhttp.ServeOnce(nil)
}
