package wasmhttp

import (
	"syscall/js"
)

// Promise is a JavaScript Promise
type Promise struct {
	js.Value
}

// NewPromise creates a new JavaScript Promise
func NewPromise(cb func(resolve func(interface{}), reject func(interface{}))) Promise {
	var cbFunc js.Func
	cbFunc = js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		defer cbFunc.Release()
		cb(
			func(value interface{}) {
				args[0].Invoke(value)
			},
			func(value interface{}) {
				args[1].Invoke(value)
			},
		)
		return js.Undefined()
	})
	return Promise{js.Global().Get("Promise").New(cbFunc)}
}

// Await waits for the Promise to be resolved and returns the value
func (p Promise) Await() (js.Value, error) {
	resCh := make(chan js.Value)
	var then js.Func
	then = js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		resCh <- args[0]
		return nil
	})
	defer then.Release()

	errCh := make(chan error)
	var catch js.Func
	catch = js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		errCh <- js.Error{args[0]}
		return nil
	})
	defer catch.Release()

	p.Call("then", then).Call("catch", catch)

	select {
	case res := <-resCh:
		return res, nil
	case err := <-errCh:
		return js.Undefined(), err
	}
}
