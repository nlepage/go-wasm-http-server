package wasmhttp

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"syscall/js"

	"github.com/nlepage/go-wasm-http-server/internal/whutil"
)

// Serve serves HTTP requests using handler or http.DefaultServeMux if handler is nil.
func Serve(handler http.Handler) func() {
	h := handler
	if h == nil {
		h = http.DefaultServeMux
	}

	path := os.Getenv("WASMHTTP_PATH")
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	if path != "" {
		prefix := os.Getenv("WASMHTTP_PATH")
		for strings.HasSuffix(prefix, "/") {
			prefix = strings.TrimSuffix(prefix, "/")
		}

		mux := http.NewServeMux()
		mux.Handle(path, http.StripPrefix(prefix, h))
		h = mux
	}

	cb := js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		jsReq := whutil.Request(args[0])

		var resolveRes func(interface{})
		var res = whutil.NewPromise(func(resolve, _ func(interface{})) {
			resolveRes = resolve
		})

		go func() {
			defer func() {
				r := recover()
				if r != nil {
					if err, ok := r.(error); ok {
						fmt.Printf("wasmhttp: panic: %+v\n", err)
					} else {
						fmt.Printf("wasmhttp: panic: %v\n", r)
					}

					res := whutil.NewResponseWriter()
					res.WriteHeader(500)
					resolveRes(res.JSResponse())
				}
			}()

			req, err := jsReq.HTTPRequest()
			if err != nil {
				panic(err)
			}

			res := whutil.NewResponseWriter()

			h.ServeHTTP(res, req)

			resolveRes(res.JSResponse())
		}()

		return res.Value()
	})

	js.Global().Get("wasmhttp").Call("registerHandler", os.Getenv("WASMHTTP_HANDLER_ID"), cb)

	return cb.Release
}
