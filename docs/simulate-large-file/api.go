package main

import (
	"bytes"
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
			bytes.NewBuffer([]byte("url must have size paramter.")).WriteTo(res)
			return
		}

		targetSize, _ := strconv.ParseInt(size[0], 0, 0)

		if targetSize < 0 || targetSize > 300 {
			bytes.NewBuffer([]byte("url must gte 0 and lte 300.")).WriteTo(res)
			return
		}

		res.Header().Set("Content-Type", "application/octet-stream")
		res.Header().Set("Content-Disposition", "attachment;filename=size-file")
		for i := int64(1); i <= 1024*targetSize; i++ {
			bytes.NewBuffer([]byte(one_kb)).WriteTo(res)
		}
		bytes.NewBuffer([]byte("")).WriteTo(res)
	})

	wasmhttp.Serve(nil)

	select {}
}
