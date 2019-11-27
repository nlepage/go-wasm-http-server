package wasmhttp

import (
	"fmt"
	"net/http"
	"os"
	"syscall/js"

	"github.com/nlepage/go-wasm-http-server/internal/whutil"
)

// Serve serves HTTP requests using handler or http.DefaultServeMux if handler is nil.
func Serve(handler http.Handler) func() {
	h := handler
	if h == nil {
		h = http.DefaultServeMux
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
						fmt.Fprintf("wasmhttp: panic: %+v", err)
					} else {
						fmt.Fprintf("wasmhttp: panic: %v", r)
					}

					res := whutil.NewResponseWriter()
					res.WriteHeader(500)
					resolveRes(res)
				}
			}()

			req, err := jsReq.HTTPRequest()
			if err != nil {
				panic(err)
			}

			res := whutil.NewResponseWriter()

			handler.ServeHTTP(res, req)

			resolveRes(res.JSResponse())
		}()

		return res.Value()
	})

	js.Global().Get("wasmhttp").Call("registerHandler", os.Getenv("WASMHTTP_HANDLER_ID"), cb)

	return cb.Release
}
