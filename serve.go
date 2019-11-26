package wasmhttp

import (
	"net/http"
	"os"
	"syscall/js"
)

// Serve serves HTTP requests using handler or http.DefaultServeMux if handler is nil
func Serve(handler http.Handler) func() {
	h := handler
	if h == nil {
		h = http.DefaultServeMux
	}

	cb := js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		jsReq := Request(args[0])

		var resolveRes func(interface{})
		var res = NewPromise(func(resolve, _ func(interface{})) {
			resolveRes = resolve
		})

		go func() {
			req, err := jsReq.HTTPRequest()
			if err != nil {
				//FIXME reject
				panic(err)
			}

			res := NewResponseWriter()

			handler.ServeHTTP(res, req)

			resolveRes(res.JSResponse())
		}()

		return res.Value()
	})

	js.Global().Get("wasmhttp").Call("RegisterHandler", os.Getenv("WASMHTTP_HANDLER_ID"), cb)

	return cb.Release
}
