package safejs

import (
	"syscall/js"

	"github.com/hack-pad/safejs"
)

type Value safejs.Value

func Safe(v js.Value) Value {
	return Value(safejs.Safe(v))
}

func Unsafe(v Value) js.Value {
	return safejs.Unsafe(safejs.Value(v))
}

func ValueOf(value any) (Value, error) {
	v, err := safejs.ValueOf(value)
	return Value(v), err
}

func (v Value) Call(m string, args ...any) (Value, error) {
	args = toJSValue(args).([]any)
	r, err := safejs.Value(v).Call(m, args...)
	return Value(r), err
}

func (v Value) Get(p string) (Value, error) {
	r, err := safejs.Value(v).Get(p)
	return Value(r), err
}

func (v Value) GetBool(p string) (bool, error) {
	bv, err := v.Get(p)
	if err != nil {
		return false, err
	}

	return safejs.Value(bv).Bool()
}

func (v Value) GetInt(p string) (int, error) {
	iv, err := v.Get(p)
	if err != nil {
		return 0, err
	}

	return safejs.Value(iv).Int()
}

func (v Value) GetString(p string) (string, error) {
	sv, err := v.Get(p)
	if err != nil {
		return "", err
	}

	return safejs.Value(sv).String()
}

func (v Value) Index(i int) (Value, error) {
	r, err := safejs.Value(v).Index(i)
	return Value(r), err
}

func (v Value) IndexString(i int) (string, error) {
	sv, err := v.Index(i)
	if err != nil {
		return "", err
	}

	return safejs.Value(sv).String()
}

func (v Value) IsNull() bool {
	return safejs.Value(v).IsNull()
}

func (v Value) IsUndefined() bool {
	return safejs.Value(v).IsUndefined()
}

func (v Value) New(args ...any) (Value, error) {
	args = toJSValue(args).([]any)
	r, err := safejs.Value(v).New(args...)
	return Value(r), err
}

func toJSValue(jsValue any) any {
	switch value := jsValue.(type) {
	case Value:
		return safejs.Value(value)
	case Func:
		return safejs.Func(value)
	case map[string]any:
		newValue := make(map[string]any)
		for mapKey, mapValue := range value {
			newValue[mapKey] = toJSValue(mapValue)
		}
		return newValue
	case []any:
		newValue := make([]any, len(value))
		for i, arg := range value {
			newValue[i] = toJSValue(arg)
		}
		return newValue
	default:
		return jsValue
	}
}
