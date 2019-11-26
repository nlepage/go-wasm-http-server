package wasmhttp

import (
	"bytes"
	"net/http"
	"syscall/js"
)

// ResponseWriter implements http.ResponseWriter
type ResponseWriter struct {
	header     http.Header
	buf        *bytes.Buffer
	statusCode int
}

// NewResponseWriter creates a new ResponseWriter
func NewResponseWriter() *ResponseWriter {
	return &ResponseWriter{
		header:     http.Header{},
		buf:        bytes.NewBuffer(nil),
		statusCode: 0,
	}
}

var _ http.ResponseWriter = (*ResponseWriter)(nil)

func (rw *ResponseWriter) Header() http.Header {
	return rw.header
}

func (rw *ResponseWriter) Write(p []byte) (int, error) {
	return rw.buf.Write(p)
}

func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
}

func (rw *ResponseWriter) JSResponse() js.Value {
	init := js.Global().Get("Object").New()

	if rw.statusCode != 0 {
		init.Set("status", rw.statusCode)
	}

	//FIXME headers

	jsBody := js.Global().Get("Uint8Array").New(rw.buf.Len())
	js.CopyBytesToJS(jsBody, rw.buf.Bytes())

	return js.Global().Get("Response").New(jsBody, init)
}
