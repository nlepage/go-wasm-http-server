package safejs

import (
	"github.com/hack-pad/safejs"
)

type Func safejs.Func

func FuncOf(fn func(this Value, args []Value) any) (Func, error) {
	r, err := safejs.FuncOf(func(this safejs.Value, args []safejs.Value) any {
		args2 := make([]Value, len(args))
		for i, v := range args {
			args2[i] = Value(v)
		}
		return fn(Value(this), []Value(args2))
	})
	return Func(r), err
}

func (f Func) Release() {
	safejs.Func(f).Release()
}

func (f Func) Value() Value {
	return Value(safejs.Func(f).Value())
}
