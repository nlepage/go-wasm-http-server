package whutil

import (
	"syscall/js"
)

// Promise is JS Promise
type Promise js.Value

// NewPromise creates a new JS Promise
func NewPromise(cb func(resolve, reject func(interface{}))) Promise {
	var cbFunc js.Func

	cbFunc = js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		defer cbFunc.Release()

		cb(
			func(v interface{}) {
				args[0].Invoke(v)
			},
			func(v interface{}) {
				args[1].Invoke(v)
			},
		)

		return js.Undefined()
	})

	return Promise(js.Global().Get("Promise").New(cbFunc))
}

// Await waits for the Promise to be resolved and returns the value
func (p Promise) Await() js.Value {
	ch := make(chan js.Value)
	p.Then(func(v js.Value) {
		ch <- v
	})
	return <-ch
}

// Then calls cb with the value when the Promise is resolved
func (p Promise) Then(cb func(js.Value)) {
	var then js.Func
	then = js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		defer then.Release()
		cb(args[0])
		return nil
	})
	js.Value(p).Call("then", then)
}

// Value returns the Promise as a js.Value
func (p Promise) Value() js.Value {
	return js.Value(p)
}
