package wasmhttp

import (
	"net/http"
	"net/http/httptest"
	"syscall/js"

	"github.com/nlepage/go-wasm-http-server/internal/readablestream"
	"github.com/nlepage/go-wasm-http-server/internal/safejs"
)

// Request builds and returns the equivalent http.Request
func Request(uvalue js.Value) (*http.Request, error) {
	value := safejs.Safe(uvalue)

	body, err := value.Get("body")
	if err != nil {
		return nil, err
	}

	r, err := body.Call("getReader")
	if err != nil {
		return nil, err
	}

	method, err := value.GetString("method")
	if err != nil {
		return nil, err
	}

	url, err := value.GetString("url")
	if err != nil {
		return nil, err
	}

	req := httptest.NewRequest(
		method,
		url,
		readablestream.NewReader(r),
	)

	headers, err := value.Get("headers")
	if err != nil {
		return nil, err
	}

	headersIt, err := headers.Call("entries")
	if err != nil {
		return nil, err
	}
	for {
		e, err := headersIt.Call("next")
		if err != nil {
			return nil, err
		}

		done, err := e.GetBool("done")
		if err != nil {
			return nil, err
		}

		if done {
			break
		}

		v, err := e.Get("value")
		if err != nil {
			return nil, err
		}

		key, err := v.IndexString(0)
		if err != nil {
			return nil, err
		}

		value, err := v.IndexString(1)
		if err != nil {
			return nil, err
		}

		req.Header.Set(key, value)
	}

	return req, nil
}
