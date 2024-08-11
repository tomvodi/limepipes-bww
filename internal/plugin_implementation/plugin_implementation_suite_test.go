package plugin_implementation_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPluginImplementation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PluginImplementation Suite")
}
