package wasmhttp

import (
	"io/ioutil"
	"net/http/httptest"
	"syscall/js"
)

// ResponseRecorder uses httptest.ResponseRecorder to build a JS Response
type ResponseRecorder struct {
	*httptest.ResponseRecorder
}

// NewResponseRecorder returns a new ResponseRecorder
func NewResponseRecorder() ResponseRecorder {
	return ResponseRecorder{httptest.NewRecorder()}
}

// JSResponse builds and returns the equivalent JS Response
func (rr ResponseRecorder) JSResponse() js.Value {
	var res = rr.Result()

	var body js.Value = js.Undefined()
	if res.ContentLength != 0 {
		var b, err = ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}
		body = js.Global().Get("Uint8Array").New(len(b))
		js.CopyBytesToJS(body, b)
	}

	var init = make(map[string]interface{}, 2)

	if res.StatusCode != 0 {
		init["status"] = res.StatusCode
	}

	if len(res.Header) != 0 {
		var headers = make(map[string]interface{}, len(res.Header))
		for k := range res.Header {
			headers[k] = res.Header.Get(k)
		}
		init["headers"] = headers
	}

	return js.Global().Get("Response").New(body, init)
}
