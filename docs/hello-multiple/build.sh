#!/bin/sh
GOOS=js GOARCH=wasm go build -o api1.wasm --ldflags="-X 'main.binaryName=api1.wasm'" .
GOOS=js GOARCH=wasm go build -o api2.wasm --ldflags="-X 'main.binaryName=api2.wasm'" .
