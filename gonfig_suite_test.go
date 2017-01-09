package gonfig

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGonfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gonfig Suite")
}
