package wasmhttp_test

import (
	"encoding/json"
	"fmt"
	"net/http"

	wasmhttp "github.com/nlepage/go-wasm-http-server"
)

// Example_JSON demostrates a simple hello JSON service.
func Example_JSON() {
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

	defer wasmhttp.Serve(nil)()

	// Wait for webpage event or use empty select{}
}
