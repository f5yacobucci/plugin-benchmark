package main

import (
	"bytes"

	"github.com/valyala/fastjson"
	wapc "github.com/wapc/wapc-guest-tinygo"
)

func guest_call_(payload []byte) ([]byte, error) {
	var p fastjson.Parser
	v, err := p.ParseBytes(payload)
	if err != nil {
		return nil, err
	}

	topic := v.GetStringBytes("topic")
	if topic == nil {
		return nil, err
	}

	data := v.GetStringBytes("data")
	if data == nil {
		return nil, err
	}

	_, _ = wapc.HostCall("benchmark.wapc", "benchmark", "callhost", nil)

	var ret bytes.Buffer
	ret.Write([]byte(`{"handled":true}`))

	return ret.Bytes(), nil
}

func main() {
	wapc.RegisterFunctions(wapc.Functions{
		"guest_call_": guest_call_,
	})
}
