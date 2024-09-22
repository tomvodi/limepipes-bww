package bwwfile

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tomvodi/limepipes-plugin-bww/internal/common"
	"github.com/tomvodi/limepipes-plugin-bww/internal/interfaces"
	"github.com/tomvodi/limepipes-plugin-bww/internal/interfaces/mocks"
	"github.com/tomvodi/limepipes-plugin-bww/internal/structure"
	"github.com/tomvodi/limepipes-plugin-bww/internal/utils"
)

var _ = Describe("StructureParser", func() {
	utils.SetupConsoleLogger()
	var err error
	var bwwFile *structure.BwwFile
	var data []byte
	var tokens []*common.Token
	var parser interfaces.StructureParser
	var conv *mocks.TokenStructureConverter
	var tokenizer *mocks.FileTokenizer

	BeforeEach(func() {
		data = []byte("test data")
		tokens = []*common.Token{
			{Value: "test", Line: 1, Col: 1},
		}

		tokenizer = mocks.NewFileTokenizer(GinkgoT())
		conv = mocks.NewTokenStructureConverter(GinkgoT())
		parser = NewStructureParser(tokenizer, conv)
	})

	JustBeforeEach(func() {
		bwwFile, err = parser.ParseDocumentStructure(data)
	})

	When("tokenizer returns an error", func() {
		BeforeEach(func() {
			tokenizer.EXPECT().Tokenize(data).
				Return(nil, fmt.Errorf("tokenizer error"))
		})

		It("should return an error", func() {
			Expect(err).Should(HaveOccurred())
		})
	})

	When("converter returns an error", func() {
		BeforeEach(func() {
			tokenizer.EXPECT().Tokenize(data).
				Return(tokens, nil)
			conv.EXPECT().Convert(tokens).
				Return(nil, fmt.Errorf("converter error"))
		})

		It("should return an error", func() {
			Expect(err).Should(HaveOccurred())
		})
	})

	When("parsing succeeds", func() {
		cFile := &structure.BwwFile{
			BagpipePlayerVersion: "123",
		}
		BeforeEach(func() {
			tokenizer.EXPECT().Tokenize(data).
				Return(tokens, nil)
			conv.EXPECT().Convert(tokens).
				Return(cFile, nil)
		})

		It("should return the parsed and converted file", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(bwwFile).Should(Equal(cFile))
		})
	})

})
