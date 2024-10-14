package main

import (
	"net/http"
	"time"

	"github.com/tmaxmax/go-sse"

	wasmhttp "github.com/nlepage/go-wasm-http-server"
)

func main() {
	s := &sse.Server{}
	t, _ := sse.NewType("ping")

	go func() {
		m := &sse.Message{
			Type: t,
		}
		m.AppendData("Hello world")

		for range time.Tick(time.Second) {
			_ = s.Publish(m)
		}
	}()

	http.Handle("/events", s)

	if _, err := wasmhttp.Serve(nil); err != nil {
		panic(err)
	}

	select {}
}
