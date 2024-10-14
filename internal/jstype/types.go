package jstype

import (
	"syscall/js"

	"github.com/nlepage/go-wasm-http-server/v2/internal/safejs"
)

var (
	ReadableStream = safejs.Safe(js.Global().Get("ReadableStream"))
	Response       = safejs.Safe(js.Global().Get("Response"))
	Uint8Array     = safejs.Safe(js.Global().Get("Uint8Array"))
)
