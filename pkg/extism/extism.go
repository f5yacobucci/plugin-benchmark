package extism

/*
#cgo LDFLAGS: -L/usr/local/lib -lextism
#include <extism.h>
EXTISM_GO_FUNCTION(callhost);
*/
import "C"

import (
	"errors"
	"runtime"
	"runtime/cgo"
	"unsafe"

	ext "github.com/extism/extism"
)

type FakePlugin struct {
	ctx       ext.Context
	inst      *ext.Plugin
	functions []ext.Function
}

type Hoster interface {
	CallHost([]byte) ([]byte, error)
}

func NewFakePlugin(manifest ext.Manifest) (*FakePlugin, error) {
	f := &FakePlugin{}
	runtime.SetFinalizer(f, free)

	f.ctx = ext.NewContext()

	f.functions = append(
		f.functions,
		ext.NewFunction(
			"callhost",
			[]ext.ValType{
				ext.I64,
				ext.I64,
			},
			[]ext.ValType{
				ext.I64,
			},
			C.callhost,
			f,
		),
	)

	var err error
	inst, err := f.ctx.PluginFromManifest(manifest, f.functions, true)
	if err != nil {
		return nil, err
	}
	f.inst = &inst

	valid := f.inst.FunctionExists("callguest")
	if !valid {
		return nil, errors.New("BAD WASM")
	}

	return f, nil
}

func free(f *FakePlugin) {
	for i := range f.functions {
		f.functions[i].Free()
	}
	if f.inst != nil {
		f.inst.Free()
	}
	f.ctx.Free()
}

//export callhost
func callhost(
	plugin unsafe.Pointer,
	inputs *C.ExtismVal, nInputs C.ExtismSize,
	outputs *C.ExtismVal, nOutputs C.ExtismSize,
	userData uintptr,
) {
	v := cgo.Handle(userData)
	output := unsafe.Slice(outputs, nOutputs)
	p, ok := v.Value().(Hoster)
	if !ok {
		ptr := unsafe.Pointer(&output[0])
		ext.ValSetI64(ptr, 1)
		return
	}

	_, _ = p.CallHost(nil)

	ptr := unsafe.Pointer(&output[0])
	ext.ValSetI64(ptr, 0)
	return
}

func (fake *FakePlugin) CallHost(_ []byte) ([]byte, error) {
	return nil, nil
}

func (fake *FakePlugin) CallGuest(payload []byte) ([]byte, error) {
	out, err := fake.inst.Call("callguest", payload)
	return out, err
}
