package safejs

import "github.com/hack-pad/safejs"

func CopyBytesToGo(dst []byte, src Value) (int, error) {
	return safejs.CopyBytesToGo(dst, safejs.Value(src))
}

func CopyBytesToJS(dst Value, src []byte) (int, error) {
	return safejs.CopyBytesToJS(safejs.Value(dst), src)
}
