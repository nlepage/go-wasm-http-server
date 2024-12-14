package wasmhttp

import (
	"io"
	"net/http"
	"net/url"
	"syscall/js"

	promise "github.com/nlepage/go-js-promise"
	"github.com/nlepage/go-wasm-http-server/v2/internal/readablestream"
	"github.com/nlepage/go-wasm-http-server/v2/internal/safejs"
)

// Request builds and returns the equivalent http.Request
func Request(uvalue js.Value) (*http.Request, error) {
	value := safejs.Safe(uvalue)

	method, err := value.GetString("method")
	if err != nil {
		return nil, err
	}

	rawURL, err := value.GetString("url")
	if err != nil {
		return nil, err
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	body, err := value.Get("body")
	if err != nil {
		return nil, err
	}

	var bodyReader io.Reader

	if !body.IsNull() {
		// WORKAROUND: Firefox does not have request.body ReadableStream
		if body.IsUndefined() {
			blobp, err := value.Call("blob")
			if err != nil {
				return nil, err
			}

			blob, err := promise.Await(safejs.Unsafe(blobp))
			if err != nil {
				return nil, err
			}

			body, err = safejs.Safe(blob).Call("stream")
			if err != nil {
				return nil, err
			}
		}

		r, err := body.Call("getReader")
		if err != nil {
			return nil, err
		}

		bodyReader = readablestream.NewReader(r)
	}

	req := &http.Request{
		Method: method,
		URL:    u,
		Body:   io.NopCloser(bodyReader),
		Header: make(http.Header),
	}

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
