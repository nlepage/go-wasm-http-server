package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	wasmhttp "github.com/nlepage/go-wasm-http-server"
)

var no = 1

func main() {
	http.HandleFunc("/hello", func(res http.ResponseWriter, req *http.Request) {
		h, m, s := time.Now().Clock()
		fmt.Printf("hello at %02d:%02d:%02d\n", h, m, s)

		params := make(map[string]string)
		if err := json.NewDecoder(req.Body).Decode(&params); err != nil {
			panic(err)
		}

		if err := json.NewEncoder(res).Encode(map[string]string{
			"message": fmt.Sprintf("Hello %s! (%d)", params["name"], no),
		}); err != nil {
			panic(err)
		}
		no++
	})

	wasmhttp.Serve(nil)

	go func() {
		for range time.Tick(time.Second) {
			h, m, s := time.Now().Clock()
			fmt.Printf("tick at %02d:%02d:%02d\n", h, m, s)
		}
	}()

	select {}
}
