package wasmhttp

import (
	"syscall/js"
)

// NewPromise creates a new JavaScript Promise
func NewPromise() (p js.Value, resolve func(interface{}), reject func(interface{})) {
	var cbFunc js.Func
	cbFunc = js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		cbFunc.Release()

		resolve = func(value interface{}) {
			args[0].Invoke(value)
		}

		reject = func(value interface{}) {
			args[1].Invoke(value)
		}

		return js.Undefined()
	})

	p = js.Global().Get("Promise").New(cbFunc)

	return
}

// Await waits for the Promise to be resolved and returns the value
func Await(p js.Value) (js.Value, error) {
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
