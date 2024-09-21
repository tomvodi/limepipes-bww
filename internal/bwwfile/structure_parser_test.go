package bwwfile_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tomvodi/limepipes-plugin-bww/internal/bwwfile"
	"github.com/tomvodi/limepipes-plugin-bww/internal/bwwfile/interfaces"
	"github.com/tomvodi/limepipes-plugin-bww/internal/structure"
	"github.com/tomvodi/limepipes-plugin-bww/internal/utils"
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

var _ = Describe("StructureParser", func() {
	utils.SetupConsoleLogger()
	var err error
	var bwwFile *structure.BwwFile
	var parser interfaces.StructureParser

	BeforeEach(func() {
		parser = bwwfile.NewStructureParser()
	})

	When("parsing a file with a staff with 4 measures in it", func() {
		BeforeEach(func() {
			data := dataFromFile("./testfiles/four_measures.bww")
			bwwFile, err = parser.ParseDocumentStructure(data)
		})

		It("should have parsed 4 measures", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(bwwFile.TuneDefs).To(HaveLen(1))
			Expect(bwwFile.TuneDefs[0].Tune.Staffs).To(HaveLen(1))
			Expect(bwwFile.TuneDefs[0].Tune.Staffs[0].Measures).To(HaveLen(4))
			Expect(bwwFile.TuneDefs[0].Tune).To(Equal(
				structure.Tune{
					Header: structure.TuneHeader{
						Title: "Tune Title",
					},
					Staffs: []structure.Staff{
						{
							Measures: []structure.Measure{
								{
									Components: []any{
										&structure.MusicSymbol{
											Pos: structure.Position{
												Line:   0,
												Column: 0,
											},
											Text: "4_4",
										},
									},
								},
								{
									Components: []any{
										&structure.MusicSymbol{
											Pos: structure.Position{
												Line:   0,
												Column: 0,
											},
											Text: "LA_4",
										},
									},
								},
							},
						},
					},
				}))
		})
	})

})
