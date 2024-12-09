package readablestream

import (
	"context"
	"io"

	"github.com/nlepage/go-wasm-http-server/v2/internal/jstype"
	"github.com/nlepage/go-wasm-http-server/v2/internal/safejs"
)

type Writer struct {
	Value      safejs.Value
	controller safejs.Value
	ctx        context.Context
	cancelled  bool
}

var _ io.WriteCloser = (*Writer)(nil)

func NewWriter() (*Writer, error) {
	var rs *Writer

	var start safejs.Func
	var controller safejs.Value

	start, err := safejs.FuncOf(func(_ safejs.Value, args []safejs.Value) any {
		defer start.Release()
		controller = args[0]
		return nil
	})
	if err != nil {
		return nil, err
	}

	var cancel safejs.Func
	ctx, cancelCtx := context.WithCancel(context.Background())

	cancel, err = safejs.FuncOf(func(_ safejs.Value, _ []safejs.Value) any {
		defer cancel.Release()
		rs.cancelled = true
		cancelCtx()
		return nil
	})
	if err != nil {
		return nil, err
	}

	source, err := safejs.ValueOf(map[string]any{
		"start":  safejs.Unsafe(start.Value()),
		"cancel": safejs.Unsafe(cancel.Value()),
	})
	if err != nil {
		return nil, err
	}

	value, err := jstype.ReadableStream.New(source)
	if err != nil {
		return nil, err
	}

	rs = &Writer{
		Value:      value,
		controller: controller,
		ctx:        ctx,
	}

	return rs, nil
}

func (rs *Writer) Write(b []byte) (int, error) {
	if rs.cancelled {
		return 0, nil
	}

	chunk, err := jstype.Uint8Array.New(len(b)) // FIXME reuse same Uint8Array?
	if err != nil {
		return 0, err
	}

	n, err := safejs.CopyBytesToJS(chunk, b)
	if err != nil {
		return 0, err
	}

	_, err = rs.controller.Call("enqueue", chunk)

	return n, err
}

func (rs *Writer) Close() error {
	if rs.cancelled {
		return nil
	}

	_, err := rs.controller.Call("close")
	return err
}

func (rs *Writer) Context() context.Context {
	return rs.ctx
}
