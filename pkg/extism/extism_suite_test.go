package extism

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestExtism(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Extism Suite")
}
