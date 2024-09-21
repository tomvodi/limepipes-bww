package bwwfile_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBwwfile(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bwwfile Suite")
}
