package wasmhttp

import (
	"bytes"
	"net/http"
	"syscall/js"
)

// Request is a JS Request
type Request js.Value

// HTTPRequest builds and returns this equivalent http.Request
func (r *Request) HTTPRequest() (*http.Request, error) {
	rValue := js.Value(*r)

	jsBody := js.Global().Get("Uint8Array").New(Promise(rValue.Call("arrayBuffer")).Await())
	body := make([]byte, jsBody.Get("length").Int())
	js.CopyBytesToGo(body, jsBody)

	req, err := http.NewRequest(
		rValue.Get("method").String(),
		rValue.Get("url").String(),
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}

	headersIt := rValue.Get("headers").Call("entries")
	for {
		v := headersIt.Call("next")
		if v.Get("done").Bool() {
			break
		}
		req.Header.Set(v.Index(0).String(), v.Index(1).String())
	}

	return req, nil
}
