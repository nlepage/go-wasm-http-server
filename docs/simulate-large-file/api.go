package main

import (
	"net/http"
	"strconv"

	wasmhttp "github.com/nlepage/go-wasm-http-server"
)

func main() {

	one_kb := ""
	for i := 0; i < 1024; i++ {
		one_kb += "0"
	}

	http.HandleFunc("/create-file-by-size", func(res http.ResponseWriter, req *http.Request) {
		vars := req.URL.Query()
		size, ok := vars["size"]
		if !ok {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte("URL must have size parameter."))
			return
		}

		targetSize, err := strconv.ParseInt(size[0], 0, 0)

		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte("Size must be 0 or a positive integer."))
			return
		}

		if targetSize < 0 || targetSize > 300 {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte("Size must be >= 0 and <= 300."))
			return
		}

		res.Header().Set("Content-Type", "application/octet-stream")
		res.Header().Set("Content-Disposition", "attachment;filename=size-file")
		for i := int64(1); i <= 1024*targetSize; i++ {
			res.Write([]byte(one_kb))
		}
	})

	wasmhttp.Serve(nil)

	select {}
}
