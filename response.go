package wasmhttp

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"syscall/js"

	promise "github.com/nlepage/go-js-promise"

	"github.com/nlepage/go-wasm-http-server/internal/jstype"
	"github.com/nlepage/go-wasm-http-server/internal/readablestream"
	"github.com/nlepage/go-wasm-http-server/internal/safejs"
)

type Response interface {
	http.ResponseWriter
	io.StringWriter
	http.Flusher
	io.Closer
	Context() context.Context
	WriteError(string)
	JSValue() js.Value
}

type response struct {
	header      http.Header
	wroteHeader bool

	promise js.Value
	resolve func(any)

	rs   *readablestream.Writer
	body *bufio.Writer
}

func NewResponse() (Response, error) {
	rs, err := readablestream.NewWriter()
	if err != nil {
		return nil, err
	}

	promise, resolve, _ := promise.New()

	return &response{
		promise: promise,
		resolve: resolve,

		rs:   rs,
		body: bufio.NewWriter(rs),
	}, nil
}

var _ Response = (*response)(nil)

// Header implements [http.ResponseWriter].
func (r *response) Header() http.Header {
	if r.header == nil {
		r.header = make(http.Header)
	}
	return r.header
}

func (r *response) headerValue() map[string]any {
	h := r.Header()
	hh := make(map[string]any, len(h)+1)
	for k := range h {
		hh[k] = h.Get(k)
	}
	return hh
}

// Write implements http.ResponseWriter.
func (r *response) Write(buf []byte) (int, error) {
	r.writeHeader(buf, "")
	return r.body.Write(buf)
}

// WriteHeader implements [http.ResponseWriter].
func (r *response) WriteHeader(code int) {
	if r.wroteHeader {
		return
	}

	checkWriteHeaderCode(code)

	init, err := safejs.ValueOf(map[string]any{
		"code":    code,
		"headers": r.headerValue(),
	})
	if err != nil {
		panic(err)
	}

	res, err := jstype.Response.New(r.rs.Value, init)
	if err != nil {
		panic(err)
	}

	r.wroteHeader = true

	r.resolve(safejs.Unsafe(res))
}

// WriteString implements [io.StringWriter].
func (r *response) WriteString(str string) (int, error) {
	r.writeHeader(nil, str)
	return r.body.WriteString(str)
}

// Flush implements [http.Flusher]
func (r *response) Flush() {
	if !r.wroteHeader {
		r.WriteHeader(200)
	}
	if err := r.body.Flush(); err != nil {
		panic(err)
	}
}

// Close implements [io.Closer]
func (r *response) Close() error {
	if err := r.body.Flush(); err != nil {
		return err
	}
	return r.rs.Close()
}

func (r *response) Context() context.Context {
	return r.rs.Context()
}

func (r *response) WriteError(str string) {
	slog.Error(str)
	if !r.wroteHeader {
		r.WriteHeader(500)
		_, _ = r.WriteString(str)
	}
}

func (r *response) JSValue() js.Value {
	return r.promise
}

func (r *response) writeHeader(b []byte, str string) {
	if r.wroteHeader {
		return
	}

	m := r.Header()

	_, hasType := m["Content-Type"]
	hasTE := m.Get("Transfer-Encoding") != ""
	if !hasType && !hasTE {
		if b == nil {
			if len(str) > 512 {
				str = str[:512]
			}
			b = []byte(str)
		}
		m.Set("Content-Type", http.DetectContentType(b))
	}

	r.WriteHeader(200)
}

func checkWriteHeaderCode(code int) {
	if code < 100 || code > 999 {
		panic(fmt.Sprintf("invalid WriteHeader code %v", code))
	}
}
