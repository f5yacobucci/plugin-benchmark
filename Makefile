.PHONY: wapc extism testbin test test-wapc test-extism
SHELL=/bin/bash

testbin: wapc extism

wapc:
	tinygo build -o testdata/generic-wapc.wasm -scheduler=none -target=wasi wasm/generic-wapc.go

extism:
	tinygo build -o testdata/generic-extism.wasm -scheduler=none -target=wasi wasm/generic-extism.go
