package whutil

import (
	"bytes"
	"net/http"
	"syscall/js"
)

// Request is a JS Request
type Request struct {
	js.Value
}

// HTTPRequest builds and returns the equivalent http.Request
func (r Request) HTTPRequest() (*http.Request, error) {
	jsBody := js.Global().Get("Uint8Array").New(Promise{r.Call("arrayBuffer")}.Await())
	body := make([]byte, jsBody.Get("length").Int())
	js.CopyBytesToGo(body, jsBody)

	req, err := http.NewRequest(
		r.Get("method").String(),
		r.Get("url").String(),
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}

	headersIt := r.Get("headers").Call("entries")
	for {
		v := headersIt.Call("next")
		if v.Get("done").Bool() {
			break
		}
		req.Header.Set(v.Index(0).String(), v.Index(1).String())
	}

	return req, nil
}
