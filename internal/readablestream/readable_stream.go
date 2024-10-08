package readablestream

import (
	"context"
	"io"

	"github.com/nlepage/go-wasm-http-server/internal/jstype"
	"github.com/nlepage/go-wasm-http-server/internal/safejs"
)

type Writer struct {
	Value      safejs.Value
	controller safejs.Value
	ctx        context.Context
}

var _ io.WriteCloser = (*Writer)(nil)

func NewWriter() (*Writer, error) {
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

	return &Writer{
		Value:      value,
		controller: controller,
		ctx:        ctx,
	}, nil
}

func (rs *Writer) Write(b []byte) (int, error) {
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
	rs.controller.Call("close")
	return nil
}

func (rs *Writer) Context() context.Context {
	return rs.ctx
}
