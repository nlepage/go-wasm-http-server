package wasmhttp

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"syscall/js"
)

// Request builds and returns the equivalent http.Request
func Request(r js.Value) *http.Request {
	jsBody := js.Global().Get("Uint8Array").New(Promise{r.Call("arrayBuffer")}.Await())
	body := make([]byte, jsBody.Get("length").Int())
	js.CopyBytesToGo(body, jsBody)

	req := httptest.NewRequest(
		r.Get("method").String(),
		r.Get("url").String(),
		bytes.NewBuffer(body),
	)

	headersIt := r.Get("headers").Call("entries")
	for {
		e := headersIt.Call("next")
		if e.Get("done").Bool() {
			break
		}
		v := e.Get("value")
		req.Header.Set(v.Index(0).String(), v.Index(1).String())
	}

	return req
}
