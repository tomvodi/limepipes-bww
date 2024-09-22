package bwwfile

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tomvodi/limepipes-plugin-bww/internal/common"
	"github.com/tomvodi/limepipes-plugin-bww/internal/structure"
)

var _ = Describe("TokenStructureConverter", func() {
	var err error
	var bwwFile *structure.BwwFile
	var tokens []*common.Token
	var tc *TokenConverter

	BeforeEach(func() {
		tc = NewTokenConverter()
	})

	JustBeforeEach(func() {
		bwwFile, err = tc.Convert(tokens)
	})

	When("converting a file with a tune and two measures", func() {
		BeforeEach(func() {
			tokens = []*common.Token{
				newToken(structure.BagpipePlayerVersion("Bagpipe Reader:1.0"), 0, 0),
				newToken(structure.TuneTitle("Tune Title"), 2, 0),
				newToken(structure.StaffStart("&"), 4, 0),
				newToken(structure.Barline("!"), 4, 6),
				newToken(structure.StaffEnd("!t"), 4, 25),
			}
		})

		It("should convert the tokens to a BwwFile", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(bwwFile).Should(BeComparableTo(&structure.BwwFile{
				BagpipePlayerVersion: "Bagpipe Reader:1.0",
				TuneDefs: []structure.TuneDefinition{
					{
						Data: []byte(`Bagpipe Reader:1.0
"Tune Title"
&
! !t
`),
						Tune: structure.Tune{
							Header: &structure.TuneHeader{
								Title: "Tune Title",
							},
							Measures: []*structure.Measure{
								{},
								{},
							},
						},
					},
				},
			}))
		})
	})

	When("converting a file with a tune and two measures with symbols", func() {
		BeforeEach(func() {
			tokens = []*common.Token{
				newToken(structure.BagpipePlayerVersion("Bagpipe Reader:1.0"), 0, 0),
				newToken(structure.TuneTitle("Tune Title"), 2, 0),
				newToken(structure.StaffStart("&"), 4, 0),
				newToken("4_4", 4, 2),
				newToken(structure.Barline("!"), 4, 6),
				newToken("LA_4", 4, 8),
				newToken(structure.StaffEnd("!t"), 4, 25),
			}
		})

		It("should convert the tokens to a BwwFile", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(bwwFile).Should(BeComparableTo(&structure.BwwFile{
				BagpipePlayerVersion: "Bagpipe Reader:1.0",
				TuneDefs: []structure.TuneDefinition{
					{
						Data: []byte(`Bagpipe Reader:1.0
"Tune Title"
& 4_4
! LA_4 !t
`),
						Tune: structure.Tune{
							Header: &structure.TuneHeader{
								Title: "Tune Title",
							},
							Measures: []*structure.Measure{
								{
									Symbols: []*structure.MusicSymbol{
										{
											Pos:  structure.Position{Line: 4, Column: 2},
											Text: "4_4",
										},
									},
								},
								{
									Symbols: []*structure.MusicSymbol{
										{
											Pos:  structure.Position{Line: 4, Column: 8},
											Text: "LA_4",
										},
									},
								},
							},
						},
					},
				},
			}))
		})
	})

	When("converting a file with measure and symbol comments", func() {
		BeforeEach(func() {
			tokens = []*common.Token{
				newToken(structure.BagpipePlayerVersion("Bagpipe Reader:1.0"), 0, 0),
				newToken(structure.TuneTitle("Tune Title"), 2, 0),
				newToken(structure.StaffStart("&"), 4, 0),
				newToken(structure.InlineText("measure inline comment"), 4, 1),
				newToken("4_4", 4, 2),
				newToken(structure.InlineText("symbol inline comment"), 4, 1),
				newToken(structure.Barline("!"), 4, 6),
				newToken(structure.InlineComment("measure comment"), 4, 1),
				newToken("LA_4", 4, 2),
				newToken(structure.InlineComment("symbol comment"), 4, 1),
				newToken(structure.StaffEnd("!t"), 4, 25),
			}
		})

		It("should convert the tokens to a BwwFile", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(bwwFile).Should(BeComparableTo(&structure.BwwFile{
				BagpipePlayerVersion: "Bagpipe Reader:1.0",
				TuneDefs: []structure.TuneDefinition{
					{
						Data: []byte(`Bagpipe Reader:1.0
"Tune Title"
& "measure inline comment",(I,L,0,0,Times New Roman,11,700,0,0,0,0,0,0) 4_4 "symbol inline comment",(I,L,0,0,Times New Roman,11,700,0,0,0,0,0,0)
! "measure comment" LA_4 "symbol comment" !t
`),
						Tune: structure.Tune{
							Header: &structure.TuneHeader{
								Title: "Tune Title",
							},
							Measures: []*structure.Measure{
								{
									InlineTexts: []structure.InlineText{
										"measure inline comment",
									},
									Symbols: []*structure.MusicSymbol{
										{
											Pos:  structure.Position{Line: 4, Column: 2},
											Text: "4_4",
											InlineTexts: []structure.InlineText{
												"symbol inline comment",
											},
										},
									},
								},
								{
									Comments: []structure.InlineComment{
										"measure comment",
									},
									Symbols: []*structure.MusicSymbol{
										{
											Comments: []structure.InlineComment{
												"symbol comment",
											},
											Pos:  structure.Position{Line: 4, Column: 2},
											Text: "LA_4",
										},
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
