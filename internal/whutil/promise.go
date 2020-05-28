package whutil

import (
	"syscall/js"
)

// Promise is JS Promise
type Promise struct {
	js.Value
}

type PromiseResolve func(...interface{}) js.Value

type PromiseReject func(...interface{}) js.Value

// NewPromise creates a new JS Promise
func NewPromise(cb func(resolve PromiseResolve, reject PromiseReject)) Promise {
	var cbFunc js.Func
	cbFunc = js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		defer cbFunc.Release()
		cb(args[0].Invoke, args[1].Invoke)
		return js.Undefined()
	})
	return Promise{js.Global().Get("Promise").New(cbFunc)}
}

// Await waits for the Promise to be resolved and returns the value
func (p Promise) Await() js.Value {
	ch := make(chan js.Value)
	var then js.Func
	then = js.FuncOf(func(_ js.Value, args []js.Value) interface{} {
		defer then.Release()
		ch <- args[0]
		return nil
	})
	p.Call("then", then)
	return <-ch
}
