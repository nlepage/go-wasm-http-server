package wasmhttp

import (
	"fmt"
	"net/http"
	"strings"
	"syscall/js"

	promise "github.com/nlepage/go-js-promise"
)

// Serve serves HTTP requests using handler or http.DefaultServeMux if handler is nil.
func Serve(handler http.Handler) func() {
	var h = handler
	if h == nil {
		h = http.DefaultServeMux
	}

	var prefix = js.Global().Get("wasmhttp").Get("path").String()
	for strings.HasSuffix(prefix, "/") {
		prefix = strings.TrimSuffix(prefix, "/")
	}

	if prefix != "" {
		var mux = http.NewServeMux()
		mux.Handle(prefix+"/", http.StripPrefix(prefix, h))
		h = mux
	}

	var cb = js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		var resPromise, resolve, reject = promise.New()

		go func() {
			defer func() {
				if r := recover(); r != nil {
					if err, ok := r.(error); ok {
						reject(fmt.Sprintf("wasmhttp: panic: %+v\n", err))
					} else {
						reject(fmt.Sprintf("wasmhttp: panic: %v\n", r))
					}
				}
			}()

			var res = NewResponseRecorder()

			h.ServeHTTP(res, Request(args[0]))

			resolve(res.JSResponse())
		}()

		return resPromise
	})

	js.Global().Get("wasmhttp").Call("setHandler", cb)

	return cb.Release
}
