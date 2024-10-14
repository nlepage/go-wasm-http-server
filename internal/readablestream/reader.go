package readablestream

import (
	"io"

	promise "github.com/nlepage/go-js-promise"

	"github.com/nlepage/go-wasm-http-server/v2/internal/safejs"
)

type Reader struct {
	value safejs.Value
	buf   []byte
	off   int
}

var _ io.Reader = (*Reader)(nil)

func NewReader(r safejs.Value) *Reader {
	return &Reader{
		value: r,
	}
}

func (r *Reader) Read(p []byte) (int, error) {
	if r.off < len(r.buf) {
		n := copy(p, r.buf[r.off:])

		r.off += n

		return n, nil
	}

	r.off = 0

	pRes, err := r.value.Call("read")
	if err != nil {
		return 0, err
	}

	ures, err := promise.Await(safejs.Unsafe(pRes))
	if err != nil {
		return 0, err
	}

	res := safejs.Safe(ures)

	done, err := res.GetBool("done")
	if err != nil {
		return 0, err
	}
	if done {
		return 0, io.EOF
	}

	value, err := res.Get("value")
	if err != nil {
		return 0, err
	}

	l, err := value.GetInt("length")
	if err != nil {
		return 0, err
	}

	if cap(r.buf) < l {
		r.buf = make([]byte, l)
	}
	if len(r.buf) < cap(r.buf) {
		r.buf = r.buf[:cap(r.buf)]
	}

	n, err := safejs.CopyBytesToGo(r.buf, value)
	if err != nil {
		return 0, err
	}

	r.buf = r.buf[:n]

	n = copy(p, r.buf[r.off:])

	r.off += n

	return n, nil
}
