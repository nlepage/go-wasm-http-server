package wasmhttp

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"syscall/js"

	"github.com/nlepage/go-wasm-http-server/v2/internal/safejs"
)

var (
	wasmhttp = safejs.Safe(js.Global().Get("wasmhttp"))
)

// Serve serves HTTP requests using handler or http.DefaultServeMux if handler is nil.
func Serve(handler http.Handler) (func(), error) {
	h := handler
	if h == nil {
		h = http.DefaultServeMux
	}

	path := os.Getenv("WASM_HTTP_PATH")
	prefix := path

	for strings.HasSuffix(prefix, "/") {
		prefix = strings.TrimSuffix(prefix, "/")
	}

	if prefix != "" {
		mux := http.NewServeMux()
		mux.Handle(prefix+"/", http.StripPrefix(prefix, h))
		h = mux
	}

	handlerValue, err := safejs.FuncOf(func(_ safejs.Value, args []safejs.Value) interface{} {
		res, err := NewResponse()
		if err != nil {
			panic(err)
		}

		go func() {
			ctx, cancel := context.WithCancel(res.Context())

			defer func() {
				cancel()
			}()

			defer func() {
				if err := res.Close(); err != nil {
					panic(err)
				}
			}()

			defer func() {
				if r := recover(); r != nil {
					res.WriteError(fmt.Sprintf("%+v", r))
				}
			}()

			req, err := Request(safejs.Unsafe(args[0]))
			if err != nil {
				res.WriteError(fmt.Sprintf("%+v", err))
				return
			}

			req = req.WithContext(ctx)

			h.ServeHTTP(res, req)
		}()

		return res.JSValue()
	})
	if err != nil {
		return nil, err
	}

	if _, err = wasmhttp.Call("setHandler", path, handlerValue); err != nil {
		return nil, err
	}

	return handlerValue.Release, nil
}
