package wapc

import (
	"bytes"
	"context"
	"errors"
	"sync"

	"github.com/bytecodealliance/wasmtime-go"
	iradix "github.com/hashicorp/go-immutable-radix/v2"
	"github.com/wapc/wapc-go"
	wapcwasmtime "github.com/wapc/wapc-go/engines/wasmtime"
)

type (
	Callback   func([]byte) ([]byte, error)
	FakeRouter struct {
		lock *sync.Mutex
		tree *iradix.Tree[Callback]
	}
)

func (r *FakeRouter) Route(_ context.Context, binding, namespace, operation string, payload []byte) ([]byte, error) {
	var key bytes.Buffer
	key.Write([]byte(namespace))
	key.Write([]byte(":"))
	key.Write([]byte(operation))
	if binding != "" {
		key.Write([]byte(":"))
		key.Write([]byte(binding))
	}

	r.lock.Lock()
	cb, found := r.tree.Get(key.Bytes())
	r.lock.Unlock()
	if !found {
		return nil, errors.New("2: ALERT ALERT ALERT")
	}
	return cb(payload)
}

func (r *FakeRouter) RegisterBinding(binding, namespace, operation string, cb Callback) (bool, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	var key bytes.Buffer
	key.Write([]byte(namespace))
	key.Write([]byte(":"))
	key.Write([]byte(operation))
	if binding != "" {
		key.Write([]byte(":"))
		key.Write([]byte(binding))
	}

	t, cb, set := r.tree.Insert(key.Bytes(), cb)
	r.tree = t
	if !set {
		return false, errors.New("ALERT ALERT ALERT")
	}
	if cb != nil {
		return true, nil
	}
	return false, nil
}

func CreateEngine() wapc.Engine {
	return wapcwasmtime.EngineWithRuntime(
		func() (*wasmtime.Engine, error) {
			config := wasmtime.NewConfig()
			config.SetWasmMemory64(true)
			return wasmtime.NewEngineWithConfig(config), nil
		},
	)
}

type FakePlugin struct {
	Inst wapc.Instance
}

func NewFakePlugin(mod wapc.Module, r *FakeRouter) (*FakePlugin, error) {
	inst, err := mod.Instantiate(context.Background())
	if err != nil {
		return nil, err
	}

	f := &FakePlugin{
		Inst: inst,
	}
	r.RegisterBinding("benchmark.wapc", "benchmark", "callhost", f.CallHost)
	return f, nil
}

func (fake *FakePlugin) CallHost(_ []byte) ([]byte, error) {
	return nil, nil
}

func (fake *FakePlugin) CallGuest(payload []byte) ([]byte, error) {
	return fake.Inst.Invoke(context.Background(), "guest_call_", payload)
}
