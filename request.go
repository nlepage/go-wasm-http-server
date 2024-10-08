package wasmhttp

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"syscall/js"

	promise "github.com/nlepage/go-js-promise"

	"github.com/nlepage/go-wasm-http-server/internal/jstype"
	"github.com/nlepage/go-wasm-http-server/internal/safejs"
)

// Request builds and returns the equivalent http.Request
func Request(ur js.Value) (*http.Request, error) {
	r := safejs.Safe(ur)

	ab, err := r.Call("arrayBuffer")
	if err != nil {
		return nil, err
	}

	u8a, err := jstype.Uint8Array.New(promise.Await(safejs.Unsafe(ab)))
	if err != nil {
		return nil, err
	}

	l, err := u8a.GetInt("length")
	if err != nil {
		return nil, err
	}

	b := make([]byte, l)

	_, err = safejs.CopyBytesToGo(b, u8a)
	if err != nil {
		return nil, err
	}

	method, err := r.GetString("method")
	if err != nil {
		return nil, err
	}

	url, err := r.GetString("url")
	if err != nil {
		return nil, err
	}

	req := httptest.NewRequest(
		method,
		url,
		bytes.NewReader(b),
	)

	headers, err := r.Get("headers")
	if err != nil {
		return nil, err
	}

	headersIt, err := headers.Call("entries")
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
