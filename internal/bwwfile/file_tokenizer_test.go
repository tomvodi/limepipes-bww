package bwwfile

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tomvodi/limepipes-plugin-bww/internal/common"
	"github.com/tomvodi/limepipes-plugin-bww/internal/filestructure"
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
			_, err := fmt.Printf("newToken(\"%v\", %d, %d),\n", t, tok.Line, tok.Col)
			Expect(err).ShouldNot(HaveOccurred())
		default:
			_, err := fmt.Printf("newToken(%T(\"%v\"), %d, %d),\n", tok.Value, tok.Value, tok.Line, tok.Col)
			Expect(err).ShouldNot(HaveOccurred())
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
					newToken(filestructure.BagpipePlayerVersion("Bagpipe Reader:1.0"), 0, 0),
					newToken(filestructure.TuneTitle("Tune Title"), 2, 0),
					newToken(filestructure.StaffStart("&"), 4, 0),
					newToken("4_4", 4, 2),
					newToken(filestructure.Barline("!"), 4, 6),
					newToken("LA_4", 4, 8),
					newToken(filestructure.Barline("!"), 4, 13),
					newToken("B_4", 4, 15),
					newToken(filestructure.Barline("!"), 4, 19),
					newToken("C_4", 4, 21),
					newToken(filestructure.StaffEnd("!t"), 4, 25),
				}),
			)
		})
	})

	When("tokenize a file with a staff with inline comments", func() {
		BeforeEach(func() {
			data = dataFromFile("./testfiles/tune_with_inline_comments.bww")
		})

		It("should tokenize the file", func() {
			printTokens(tokens)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(tokens).Should(BeComparableTo(
				[]*common.Token{
					newToken(filestructure.BagpipePlayerVersion("Bagpipe Reader:1.0"), 0, 0),
					newToken(filestructure.TuneComment("just a comment"), 2, 0),
					newToken(filestructure.TuneTitle("Tune Title"), 4, 0),
					newToken(filestructure.TuneInline("tune inline text"), 5, 0),
					newToken(filestructure.TuneComment("and another comment"), 6, 0),
					newToken(filestructure.StaffInline("staff inline text"), 8, 0),
					newToken(filestructure.StaffComment("staff comment"), 9, 0),
					newToken(filestructure.StaffStart("&"), 11, 0),
					newToken("LA_4", 11, 2),
					newToken(filestructure.StaffEnd("!t"), 11, 7),
					newToken(filestructure.StaffInline("staff inline text"), 13, 0),
					newToken(filestructure.StaffComment("staff comment"), 14, 0),
					newToken(filestructure.StaffStart("&"), 16, 0),
					newToken("D_4", 16, 3),
					newToken(filestructure.StaffEnd("!I"), 16, 7),
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
					newToken(filestructure.BagpipePlayerVersion("Bagpipe Reader:1.0"), 0, 0),
					newToken(filestructure.TuneTitle("Tune Title"), 3, 0),
					newToken(filestructure.StaffStart("&"), 6, 0),
					newToken(filestructure.InlineComment("comment measure"), 6, 2),
					newToken("LA_4", 6, 20),
					newToken(filestructure.InlineComment("comment symbol"), 6, 25),
					newToken(filestructure.Barline("!"), 7, 0),
					newToken("B_4", 7, 2),
					newToken(filestructure.InlineText("comment with inline style"), 7, 6),
					newToken(filestructure.StaffEnd("!t"), 7, 78),
				}),
			)
		})
	})

	When("tokenize a file with metadata", func() {
		BeforeEach(func() {
			data = []byte(`Bagpipe Reader:1.0
MIDINoteMappings,(54,56,58,59,61,63,64,66,68,56,58,60,61,63,65,66,68,70,55,57,59,60,62,64,65,67,69)
FrequencyMappings,(370,415,466,494,554,622,659,740,831,415,466,523,554,622,699,740,831,932,392,440,494,523,587,659,699,784,880)
InstrumentMappings,(71,71,45,33,1000,60,70)
GracenoteDurations,(20,40,30,50,100,200,800,1200,250,250,250,500,200)
FontSizes,(90,100,100,80,250)
TuneFormat,(1,0,M,L,500,500,500,500,P,0,0)
`)
		})

		It("should tokenize the file and skip metadata", func() {
			printTokens(tokens)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(tokens).Should(BeComparableTo(
				[]*common.Token{
					newToken(filestructure.BagpipePlayerVersion("Bagpipe Reader:1.0"), 0, 0),
				}),
			)
		})
	})

	When("tokenize a file with tune tempo", func() {
		BeforeEach(func() {
			data = []byte(`Bagpipe Reader:1.0
"Tune Title",(T,L,0,0,Times New Roman,15,700,0,1,2,0,0,32768)
TuneTempo,105
& TuneTempo,80 C_4  !t
`)
		})

		It("should tokenize the file", func() {
			printTokens(tokens)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(tokens).Should(BeComparableTo(
				[]*common.Token{
					newToken(filestructure.BagpipePlayerVersion("Bagpipe Reader:1.0"), 0, 0),
					newToken(filestructure.TuneTitle("Tune Title"), 1, 0),
					newToken(filestructure.TuneTempo(105), 2, 0),
					newToken(filestructure.StaffStart("&"), 3, 0),
					newToken(filestructure.TempoChange(80), 3, 2),
					newToken("C_4", 3, 15),
					newToken(filestructure.StaffEnd("!t"), 3, 20)}),
			)
		})
	})

	When("tokenize a file with dalsegno", func() {
		BeforeEach(func() {
			data = []byte(`Bagpipe Reader:1.0
"Tune Title",(T,L,0,0,Times New Roman,15,700,0,1,2,0,0,32768)
& segno 
! C_4 !t dalsegno
`)
		})

		It("should tokenize the file", func() {
			//printTokens(tokens)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(tokens).Should(BeComparableTo(
				[]*common.Token{
					newToken(filestructure.BagpipePlayerVersion("Bagpipe Reader:1.0"), 0, 0),
					newToken(filestructure.TuneTitle("Tune Title"), 1, 0),
					newToken(filestructure.StaffStart("&"), 2, 0),
					newToken("segno", 2, 2),
					newToken(filestructure.Barline("!"), 3, 0),
					newToken("C_4", 3, 2),
					newToken(filestructure.StaffEnd("!t"), 3, 6),
					newToken(filestructure.DalSegno("dalsegno"), 3, 9),
				}),
			)
		})
	})

	When("tokenize a file with dacapoalfine", func() {
		BeforeEach(func() {
			data = []byte(`Bagpipe Reader:1.0
"Tune Title",(T,L,0,0,Times New Roman,15,700,0,1,2,0,0,32768)
& segno 
! C_4 !t dacapoalfine
`)
		})

		It("should tokenize the file", func() {
			printTokens(tokens)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(tokens).Should(BeComparableTo(
				[]*common.Token{
					newToken(filestructure.BagpipePlayerVersion("Bagpipe Reader:1.0"), 0, 0),
					newToken(filestructure.TuneTitle("Tune Title"), 1, 0),
					newToken(filestructure.StaffStart("&"), 2, 0),
					newToken("segno", 2, 2),
					newToken(filestructure.Barline("!"), 3, 0),
					newToken("C_4", 3, 2),
					newToken(filestructure.StaffEnd("!t"), 3, 6),
					newToken(filestructure.DacapoAlFine("dacapoalfine"), 3, 9),
				}),
			)
		})
	})
})
