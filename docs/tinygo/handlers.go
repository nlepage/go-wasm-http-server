package main

import (
	"encoding/json"
	"net/http"
	"runtime"
)

func goRuntimeHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(res).Encode(map[string]string{
		"os":       runtime.GOOS,
		"arch":     runtime.GOARCH,
		"compiler": runtime.Compiler,
		"version":  runtime.Version(),
	}); err != nil {
		panic(err)
	}
}
