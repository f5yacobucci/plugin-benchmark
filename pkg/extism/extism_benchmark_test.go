package extism

import (
	"os"
	"testing"

	ext "github.com/extism/extism"
	. "github.com/onsi/gomega"
)

func BenchmarkInstantiate(b *testing.B) {
	guest, err := os.ReadFile("../../testdata/generic-extism.wasm")
	Expect(err).NotTo(HaveOccurred())

	manifest := ext.Manifest{
		Wasm: []ext.Wasm{
			ext.WasmData{
				Data: guest,
			},
		},
	}

	for n := 0; n < b.N; n++ {
		_, err := NewFakePlugin(manifest)
		Expect(err).NotTo(HaveOccurred())
	}
}

func BenchmarkInstantiateAndRun(b *testing.B) {
	guest, err := os.ReadFile("../../testdata/generic-extism.wasm")
	Expect(err).NotTo(HaveOccurred())

	manifest := ext.Manifest{
		Wasm: []ext.Wasm{
			ext.WasmData{
				Data: guest,
			},
		},
	}

	for n := 0; n < b.N; n++ {
		fake, err := NewFakePlugin(manifest)
		Expect(err).NotTo(HaveOccurred())

		payload := `{"topic":"fake.example.com","data":"fake"}`
		response, err := fake.CallGuest([]byte(payload))
		Expect(string(response)).To(Equal(`{"handled":true}`))
		Expect(err).NotTo(HaveOccurred())
	}
}

func BenchmarkInstantiateAndCallMulti(b *testing.B) {
	guest, err := os.ReadFile("../../testdata/generic-extism.wasm")
	Expect(err).NotTo(HaveOccurred())

	manifest := ext.Manifest{
		Wasm: []ext.Wasm{
			ext.WasmData{
				Data: guest,
			},
		},
	}

	for n := 0; n < b.N; n++ {
		fake, err := NewFakePlugin(manifest)
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
	guest, err := os.ReadFile("../../testdata/generic-extism.wasm")
	Expect(err).NotTo(HaveOccurred())

	manifest := ext.Manifest{
		Wasm: []ext.Wasm{
			ext.WasmData{
				Data: guest,
			},
		},
	}

	fake, err := NewFakePlugin(manifest)
	Expect(err).NotTo(HaveOccurred())

	payload := `{"topic":"fake.example.com","data":"fake"}`

	for m := 0; m < b.N; m++ {
		response, err := fake.CallGuest([]byte(payload))
		Expect(string(response)).To(Equal(`{"handled":true}`))
		Expect(err).NotTo(HaveOccurred())
	}
}
