package main

import (
	"bytes"

	pdk "github.com/extism/go-pdk"
	"github.com/valyala/fastjson"
)

/*
#include "runtime/extism-pdk.h"
*/
import "C"

/*
//go:wasm-module env
//export callhost
func callhost(uint64, uint64) uint64
*/

//export callguest
func callguest() int32 {
	input := pdk.Input()

	var p fastjson.Parser
	v, err := p.ParseBytes(input)
	if err != nil {
		mem := pdk.AllocateString(err.Error())
		defer mem.Free()
		C.extism_error_set(mem.Offset())
		return -1
	}

	topic := v.GetStringBytes("topic")
	if topic == nil {
		mem := pdk.AllocateString("CANNOT FIND TOPIC")
		defer mem.Free()
		C.extism_error_set(mem.Offset())
		return -1
	}

	data := v.GetStringBytes("data")
	if data == nil {
		mem := pdk.AllocateString("CANNOT FIND DATA")
		defer mem.Free()
		C.extism_error_set(mem.Offset())
		return -1
	}

	/*
		ret := callhost(0, 0)
		if ret > 0 {
			mem := pdk.FindMemory(ret)
			buf := make([]byte, mem.Length())
			mem.Load(buf)

			errMem := pdk.AllocateBytes(buf)
			defer errMem.Free()
			C.extism_error_set(errMem.Offset())

			return -1
		}
	*/

	var b bytes.Buffer
	b.Write([]byte(`{"handled":true}`))
	mem := pdk.AllocateBytes(b.Bytes())
	defer mem.Free()
	pdk.OutputMemory(mem)

	return 0
}

func main() {}
