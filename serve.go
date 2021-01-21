package wasmhttp

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"syscall/js"

	"github.com/nlepage/go-wasm-http-server/internal/whutil"
)

// Serve serves HTTP requests using handler or http.DefaultServeMux if handler is nil.
func Serve(handler http.Handler) func() {
	var h = handler
	if h == nil {
		h = http.DefaultServeMux
	}

	var path = os.Getenv("WASMHTTP_PATH")
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	if path != "" { // FIXME always true since / suffix is added to path
		var prefix = os.Getenv("WASMHTTP_PATH")
		for strings.HasSuffix(prefix, "/") {
			prefix = strings.TrimSuffix(prefix, "/")
		}

		var mux = http.NewServeMux()
		mux.Handle(path, http.StripPrefix(prefix, h))
		h = mux
	}

	var cb = js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		var jsReq = whutil.Request{args[0]}

		var resPromise = whutil.NewPromise(func(resolve whutil.PromiseResolve, reject whutil.PromiseReject) {
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

				var req, err = jsReq.HTTPRequest()
				if err != nil {
					panic(err)
				}

				var res = whutil.NewResponseWriter()

				h.ServeHTTP(res, req)

				resolve(res)
			}()
		})

		return resPromise
	})

	js.Global().Get("wasmhttp").Call("registerHandler", os.Getenv("WASMHTTP_HANDLER_ID"), cb)

	return cb.Release
}

func ServeOnce(handler http.Handler) {
	var ctx, cancel = context.WithCancel(context.Background())

	var h = handler
	if h == nil {
		h = http.DefaultServeMux
	}

	var path = os.Getenv("WASMHTTP_PATH")
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	if path != "" { // FIXME always true since / suffix is added to path
		var prefix = os.Getenv("WASMHTTP_PATH")
		for strings.HasSuffix(prefix, "/") {
			prefix = strings.TrimSuffix(prefix, "/")
		}

		var mux = http.NewServeMux()
		mux.Handle(path, http.StripPrefix(prefix, h))
		h = mux
	}

	var cb = js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		var jsReq = whutil.Request{args[0]}

		var resPromise = whutil.NewPromise(func(resolve whutil.PromiseResolve, reject whutil.PromiseReject) {
			go func() {
				defer cancel()

				defer func() {
					if r := recover(); r != nil {
						if err, ok := r.(error); ok {
							reject(fmt.Sprintf("wasmhttp: panic: %+v\n", err))
						} else {
							reject(fmt.Sprintf("wasmhttp: panic: %v\n", r))
						}
					}
				}()

				var req, err = jsReq.HTTPRequest()
				if err != nil {
					panic(err)
				}

				var res = whutil.NewResponseWriter()

				h.ServeHTTP(res, req)

				resolve(res)
			}()
		})

		return resPromise
	})

	js.Global().Get("wasmhttp").Call("registerHandler", os.Getenv("WASMHTTP_HANDLER_ID"), cb)

	<-ctx.Done()
}
