package wapc

import (
	"context"
	"os"
	"runtime"
	"sync"
	"testing"

	iradix "github.com/hashicorp/go-immutable-radix/v2"

	. "github.com/onsi/gomega"
	"github.com/wapc/wapc-go"
)

func BenchmarkInstantiate(b *testing.B) {
	engine := CreateEngine()

	guest, err := os.ReadFile("../../testdata/generic-wapc.wasm")
	Expect(err).NotTo(HaveOccurred())

	mod, err := engine.New(
		context.Background(),
		nil,
		guest,
		&wapc.ModuleConfig{},
	)
	Expect(err).NotTo(HaveOccurred())
	runtime.SetFinalizer(
		mod,
		func(m wapc.Module) {
			m.Close(context.Background())
		},
	)

	router := &FakeRouter{
		lock: &sync.Mutex{},
		tree: iradix.New[Callback](),
	}
	for n := 0; n < b.N; n++ {
		_, err := NewFakePlugin(mod, router)
		Expect(err).NotTo(HaveOccurred())
	}
}

func BenchmarkInstantiateAndCallSingle(b *testing.B) {
	engine := CreateEngine()

	guest, err := os.ReadFile("../../testdata/generic-wapc.wasm")
	Expect(err).NotTo(HaveOccurred())

	mod, err := engine.New(
		context.Background(),
		nil,
		guest,
		&wapc.ModuleConfig{},
	)
	Expect(err).NotTo(HaveOccurred())
	runtime.SetFinalizer(
		mod,
		func(m wapc.Module) {
			m.Close(context.Background())
		},
	)

	router := &FakeRouter{
		lock: &sync.Mutex{},
		tree: iradix.New[Callback](),
	}
	for n := 0; n < b.N; n++ {
		fake, err := NewFakePlugin(mod, router)
		Expect(err).NotTo(HaveOccurred())

		payload := `{"topic":"fake.example.com","data":"fake"}`
		response, err := fake.CallGuest([]byte(payload))
		Expect(string(response)).To(Equal(`{"handled":true}`))
		Expect(err).NotTo(HaveOccurred())
	}
}

func BenchmarkInstantiateAndCallMulti(b *testing.B) {
	engine := CreateEngine()

	guest, err := os.ReadFile("../../testdata/generic-wapc.wasm")
	Expect(err).NotTo(HaveOccurred())

	mod, err := engine.New(
		context.Background(),
		nil,
		guest,
		&wapc.ModuleConfig{},
	)
	Expect(err).NotTo(HaveOccurred())
	runtime.SetFinalizer(
		mod,
		func(m wapc.Module) {
			m.Close(context.Background())
		},
	)

	router := &FakeRouter{
		lock: &sync.Mutex{},
		tree: iradix.New[Callback](),
	}
	for n := 0; n < b.N; n++ {
		fake, err := NewFakePlugin(mod, router)
		Expect(err).NotTo(HaveOccurred())

		payload := `{"topic":"fake.example.com","data":"fake"}`
		for m := 0; m < b.N; m++ {
			response, err := fake.CallGuest([]byte(payload))
			Expect(string(response)).To(Equal(`{"handled":true}`))
			Expect(err).NotTo(HaveOccurred())
		}
	}
}

func BenchmarkIsolatedCalls(b *testing.B) {
	engine := CreateEngine()

	guest, err := os.ReadFile("../../testdata/generic-wapc.wasm")
	Expect(err).NotTo(HaveOccurred())

	mod, err := engine.New(
		context.Background(),
		nil,
		guest,
		&wapc.ModuleConfig{},
	)
	Expect(err).NotTo(HaveOccurred())
	runtime.SetFinalizer(
		mod,
		func(m wapc.Module) {
			m.Close(context.Background())
		},
	)

	router := &FakeRouter{
		lock: &sync.Mutex{},
		tree: iradix.New[Callback](),
	}

	fake, err := NewFakePlugin(mod, router)
	Expect(err).NotTo(HaveOccurred())

	payload := `{"topic":"fake.example.com","data":"fake"}`

	for m := 0; m < b.N; m++ {
		response, err := fake.CallGuest([]byte(payload))
		Expect(string(response)).To(Equal(`{"handled":true}`))
		Expect(err).NotTo(HaveOccurred())
	}
}
