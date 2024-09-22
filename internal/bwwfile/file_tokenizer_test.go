package bwwfile

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tomvodi/limepipes-plugin-bww/internal/common"
	"github.com/tomvodi/limepipes-plugin-bww/internal/structure"
	"io"
	"os"
)

func dataFromFile(filePath string) []byte {
	bwwFile, err := os.Open(filePath)
	Expect(err).ShouldNot(HaveOccurred())
	var data []byte
	data, err = io.ReadAll(bwwFile)
	Expect(err).ShouldNot(HaveOccurred())

	return data
}

func newToken(value any, line, col int) *common.Token {
	return &common.Token{
		Value: value,
		Line:  line,
		Col:   col,
	}
}

//nolint:unused
func printTokens(tokens []*common.Token) {
	for _, tok := range tokens {
		switch t := tok.Value.(type) {
		case string:
			fmt.Printf("newToken(\"%v\", %d, %d),\n", t, tok.Line, tok.Col)
		default:
			fmt.Printf("newToken(%T(\"%v\"), %d, %d),\n", tok.Value, tok.Value, tok.Line, tok.Col)
		}
	}
}

var _ = Describe("FileTokenizer", func() {
	var ft *Tokenizer
	var err error
	var tokens []*common.Token
	var data []byte

	BeforeEach(func() {
		ft = NewTokenizer()
	})

	JustBeforeEach(func() {
		tokens, err = ft.Tokenize(data)
	})

	When("tokenize a file with a staff with 4 measures in it", func() {
		BeforeEach(func() {
			data = dataFromFile("./testfiles/four_measures.bww")
		})

		It("should tokenize the file", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(tokens).Should(BeComparableTo(
				[]*common.Token{
					newToken(structure.BagpipePlayerVersion("Bagpipe Reader:1.0"), 0, 0),
					newToken(structure.TuneTitle("Tune Title"), 2, 0),
					newToken(structure.StaffStart("&"), 4, 0),
					newToken("4_4", 4, 2),
					newToken(structure.Barline("!"), 4, 6),
					newToken("LA_4", 4, 8),
					newToken(structure.Barline("!"), 4, 13),
					newToken("B_4", 4, 15),
					newToken(structure.Barline("!"), 4, 19),
					newToken("C_4", 4, 21),
					newToken(structure.StaffEnd("!t"), 4, 25),
				}),
			)
		})
	})

	When("tokenize a file with a staff with inline comments", func() {
		BeforeEach(func() {
			data = dataFromFile("./testfiles/tune_with_inline_comments.bww")
		})

		It("should tokenize the file", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(tokens).Should(BeComparableTo(
				[]*common.Token{
					newToken(structure.BagpipePlayerVersion("Bagpipe Reader:1.0"), 0, 0),
					newToken(structure.TuneComment("just a comment"), 2, 0),
					newToken(structure.TuneTitle("Tune Title"), 4, 0),
					newToken(structure.TuneInline("tune inline text"), 5, 0),
					newToken(structure.TuneComment("and another comment"), 6, 0),
					newToken(structure.StaffStart("&"), 8, 0),
					newToken("LA_4", 8, 2),
					newToken(structure.StaffEnd("!t"), 8, 7),
					newToken(structure.InlineText("stave inline text"), 10, 0),
					newToken(structure.InlineComment("stave comment"), 11, 0),
					newToken(structure.StaffStart("&"), 13, 0),
					newToken("D_4", 13, 3),
					newToken("!I", 13, 7),
				}),
			)
		})
	})

	When("tokenize a file with a staff with symbol comments", func() {
		BeforeEach(func() {
			data = dataFromFile("./testfiles/tune_with_symbol_comment.bww")
		})

		It("should tokenize the file", func() {
			printTokens(tokens)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(tokens).Should(BeComparableTo(
				[]*common.Token{
					newToken(structure.BagpipePlayerVersion("Bagpipe Reader:1.0"), 0, 0),
					newToken(structure.TuneTitle("Tune Title"), 3, 0),
					newToken(structure.StaffStart("&"), 6, 0),
					newToken(structure.InlineComment("comment measure"), 6, 2),
					newToken("LA_4", 6, 20),
					newToken(structure.InlineComment("comment symbol"), 6, 25),
					newToken(structure.Barline("!"), 7, 0),
					newToken("B_4", 7, 2),
					newToken(structure.InlineText("comment with inline style"), 7, 6),
					newToken(structure.StaffEnd("!t"), 7, 78),
				}),
			)
		})
	})
})
