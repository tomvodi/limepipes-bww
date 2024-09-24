package bwwfile

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tomvodi/limepipes-plugin-bww/internal/common"
	"github.com/tomvodi/limepipes-plugin-bww/internal/filestructure"
)

var _ = Describe("TokenStructureConverter", func() {
	var err error
	var bwwFile *filestructure.BwwFile
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
				newToken(filestructure.BagpipePlayerVersion("Bagpipe Reader:1.0"), 0, 0),
				newToken(filestructure.TuneTitle("Tune Title"), 2, 0),
				newToken(filestructure.StaffStart("&"), 4, 0),
				newToken(filestructure.Barline("!"), 4, 6),
				newToken(filestructure.StaffEnd("!t"), 4, 25),
			}
		})

		It("should convert the tokens to a BwwFile", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(bwwFile).Should(BeComparableTo(&filestructure.BwwFile{
				BagpipePlayerVersion: "Bagpipe Reader:1.0",
				TuneDefs: []filestructure.TuneDefinition{
					{
						Data: []byte(`Bagpipe Reader:1.0
"Tune Title",(T,L,0,0,Times NewConverter Roman,11,700,0,0,0,0,0,0)
&
! !t
`),
						Tune: &filestructure.Tune{
							Header: &filestructure.TuneHeader{
								Title: "Tune Title",
							},
							Measures: []*filestructure.Measure{
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
				newToken(filestructure.BagpipePlayerVersion("Bagpipe Reader:1.0"), 0, 0),
				newToken(filestructure.TuneTitle("Tune Title"), 2, 0),
				newToken(filestructure.StaffStart("&"), 4, 0),
				newToken("4_4", 4, 2),
				newToken(filestructure.Barline("!"), 4, 6),
				newToken("LA_4", 4, 8),
				newToken(filestructure.StaffEnd("!t"), 4, 25),
			}
		})

		It("should convert the tokens to a BwwFile", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(bwwFile).Should(BeComparableTo(&filestructure.BwwFile{
				BagpipePlayerVersion: "Bagpipe Reader:1.0",
				TuneDefs: []filestructure.TuneDefinition{
					{
						Data: []byte(`Bagpipe Reader:1.0
"Tune Title",(T,L,0,0,Times NewConverter Roman,11,700,0,0,0,0,0,0)
& 4_4
! LA_4 !t
`),
						Tune: &filestructure.Tune{
							Header: &filestructure.TuneHeader{
								Title: "Tune Title",
							},
							Measures: []*filestructure.Measure{
								{
									Symbols: []*filestructure.MusicSymbol{
										{
											Pos:  filestructure.Position{Line: 4, Column: 2},
											Text: "4_4",
										},
									},
								},
								{
									Symbols: []*filestructure.MusicSymbol{
										{
											Pos:  filestructure.Position{Line: 4, Column: 8},
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
				newToken(filestructure.BagpipePlayerVersion("Bagpipe Reader:1.0"), 0, 0),
				newToken(filestructure.TuneTitle("Tune Title"), 2, 0),
				newToken(filestructure.StaffStart("&"), 4, 0),
				newToken(filestructure.InlineText("measure inline comment"), 4, 1),
				newToken("4_4", 4, 2),
				newToken(filestructure.InlineText("symbol inline comment"), 4, 1),
				newToken(filestructure.Barline("!"), 4, 6),
				newToken(filestructure.InlineComment("measure comment"), 4, 1),
				newToken("LA_4", 4, 2),
				newToken(filestructure.InlineComment("symbol comment"), 4, 1),
				newToken(filestructure.StaffEnd("!t"), 4, 25),
			}
		})

		It("should convert the tokens to a BwwFile", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(bwwFile).Should(BeComparableTo(&filestructure.BwwFile{
				BagpipePlayerVersion: "Bagpipe Reader:1.0",
				TuneDefs: []filestructure.TuneDefinition{
					{
						Data: []byte(`Bagpipe Reader:1.0
"Tune Title",(T,L,0,0,Times NewConverter Roman,11,700,0,0,0,0,0,0)
& "measure inline comment",(I,L,0,0,Times NewConverter Roman,11,700,0,0,0,0,0,0) 4_4 "symbol inline comment",(I,L,0,0,Times NewConverter Roman,11,700,0,0,0,0,0,0)
! "measure comment" LA_4 "symbol comment" !t
`),
						Tune: &filestructure.Tune{
							Header: &filestructure.TuneHeader{
								Title: "Tune Title",
							},
							Measures: []*filestructure.Measure{
								{
									InlineTexts: []filestructure.InlineText{
										"measure inline comment",
									},
									Symbols: []*filestructure.MusicSymbol{
										{
											Pos:  filestructure.Position{Line: 4, Column: 2},
											Text: "4_4",
											InlineTexts: []filestructure.InlineText{
												"symbol inline comment",
											},
										},
									},
								},
								{
									InlineComments: []filestructure.InlineComment{
										"measure comment",
									},
									Symbols: []*filestructure.MusicSymbol{
										{
											Comments: []filestructure.InlineComment{
												"symbol comment",
											},
											Pos:  filestructure.Position{Line: 4, Column: 2},
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

	When("converting tune with staff comments (comments that appear right before a starting staff", func() {
		BeforeEach(func() {
			tokens = []*common.Token{
				newToken(filestructure.BagpipePlayerVersion("Bagpipe Reader:1.0"), 0, 0),
				newToken(filestructure.TuneComment("just a comment"), 2, 0),
				newToken(filestructure.TuneTitle("Tune Title"), 4, 0),
				newToken(filestructure.TuneInline("tune inline text"), 5, 0),
				newToken(filestructure.TuneComment("and another tune comment"), 6, 0),
				newToken(filestructure.StaffInline("staff inline text"), 8, 0),
				newToken(filestructure.StaffComment("staff comment"), 9, 0),
				newToken(filestructure.StaffStart("&"), 11, 0),
				newToken("LA_4", 11, 2),
				newToken(filestructure.StaffEnd("!t"), 11, 7),
				newToken(filestructure.StaffInline("staff inline comment in between"), 13, 0),
				newToken(filestructure.StaffComment("staff comment in between"), 14, 0),
				newToken(filestructure.StaffStart("&"), 16, 0),
				newToken("D_4", 16, 3),
				newToken(filestructure.StaffEnd("!I"), 16, 7),
			}
		})

		It("should convert the tokens to a BwwFile", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(bwwFile).Should(BeComparableTo(&filestructure.BwwFile{
				BagpipePlayerVersion: "Bagpipe Reader:1.0",
				TuneDefs: []filestructure.TuneDefinition{
					{
						Data: []byte(`Bagpipe Reader:1.0
"just a comment"
"Tune Title",(T,L,0,0,Times NewConverter Roman,11,700,0,0,0,0,0,0)
"tune inline text",(I,L,0,0,Times NewConverter Roman,11,700,0,0,0,0,0,0)
"and another tune comment"
"staff inline text",(I,L,0,0,Times NewConverter Roman,11,700,0,0,0,0,0,0)
"staff comment"
& LA_4 !t
"staff inline comment in between",(I,L,0,0,Times NewConverter Roman,11,700,0,0,0,0,0,0)
"staff comment in between"
& D_4 !I
`),
						Tune: &filestructure.Tune{
							Header: &filestructure.TuneHeader{
								Title: "Tune Title",
								Comments: []filestructure.TuneComment{
									"just a comment",
									"and another tune comment",
								},
								InlineTexts: []filestructure.TuneInline{
									"tune inline text",
								},
							},
							Measures: []*filestructure.Measure{
								{
									StaffComments: []filestructure.StaffComment{
										"staff comment",
									},
									StaffInlineTexts: []filestructure.StaffInline{
										"staff inline text",
									},
									Symbols: []*filestructure.MusicSymbol{
										{
											Pos:  filestructure.Position{Line: 11, Column: 2},
											Text: "LA_4",
										},
									},
								},
								{
									StaffComments: []filestructure.StaffComment{
										"staff comment in between",
									},
									StaffInlineTexts: []filestructure.StaffInline{
										"staff inline comment in between",
									},
									Symbols: []*filestructure.MusicSymbol{
										{
											Pos:  filestructure.Position{Line: 16, Column: 3},
											Text: "D_4",
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

	When("converting file with two tunes", func() {
		BeforeEach(func() {
			tokens = []*common.Token{
				newToken(filestructure.BagpipePlayerVersion("Bagpipe Reader:1.0"), 0, 0),
				newToken(filestructure.TuneTitle("Tune 1 Title"), 2, 0),
				newToken(filestructure.StaffStart("&"), 5, 0),
				newToken("LA_4", 5, 2),
				newToken(filestructure.StaffEnd("!t"), 5, 7),
				newToken(filestructure.TuneTitle("Tune 2 Title"), 7, 0),
				newToken(filestructure.StaffStart("&"), 9, 0),
				newToken("B_4", 9, 2),
				newToken(filestructure.StaffEnd("!t"), 9, 6),
			}
		})

		It("should convert the tokens to a BwwFile", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(bwwFile).Should(BeComparableTo(&filestructure.BwwFile{
				BagpipePlayerVersion: "Bagpipe Reader:1.0",
				TuneDefs: []filestructure.TuneDefinition{
					{
						Data: []byte(`Bagpipe Reader:1.0
"Tune 1 Title",(T,L,0,0,Times NewConverter Roman,11,700,0,0,0,0,0,0)
& LA_4 !t
`),
						Tune: &filestructure.Tune{
							Header: &filestructure.TuneHeader{
								Title: "Tune 1 Title",
							},
							Measures: []*filestructure.Measure{
								{
									Symbols: []*filestructure.MusicSymbol{
										{
											Pos:  filestructure.Position{Line: 5, Column: 2},
											Text: "LA_4",
										},
									},
								},
							},
						},
					},
					{
						Data: []byte(`Bagpipe Reader:1.0
"Tune 2 Title",(T,L,0,0,Times NewConverter Roman,11,700,0,0,0,0,0,0)
& B_4 !t
`),
						Tune: &filestructure.Tune{
							Header: &filestructure.TuneHeader{
								Title: "Tune 2 Title",
							},
							Measures: []*filestructure.Measure{
								{
									Symbols: []*filestructure.MusicSymbol{
										{
											Pos:  filestructure.Position{Line: 9, Column: 2},
											Text: "B_4",
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

	When("converting file with two tunes where the first one doesn't have a title", func() {
		BeforeEach(func() {
			tokens = []*common.Token{
				newToken(filestructure.BagpipePlayerVersion("Bagpipe Reader:1.0"), 0, 0),
				//newToken(structure.TuneTitle("Tune 1 Title"), 2, 0),
				newToken(filestructure.StaffStart("&"), 5, 0),
				newToken("LA_4", 5, 2),
				newToken(filestructure.StaffEnd("!t"), 5, 7),
				newToken(filestructure.TuneTitle("Tune 2 Title"), 7, 0),
				newToken(filestructure.StaffStart("&"), 9, 0),
				newToken("B_4", 9, 2),
				newToken(filestructure.StaffEnd("!t"), 9, 6),
			}
		})

		It("should convert the tokens to a BwwFile", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(bwwFile).Should(BeComparableTo(&filestructure.BwwFile{
				BagpipePlayerVersion: "Bagpipe Reader:1.0",
				TuneDefs: []filestructure.TuneDefinition{
					{
						Data: []byte(`Bagpipe Reader:1.0
"No Name",(T,L,0,0,Times NewConverter Roman,11,700,0,0,0,0,0,0)
& LA_4 !t
`),
						Tune: &filestructure.Tune{
							Header: &filestructure.TuneHeader{
								Title: "No Name",
							},
							Measures: []*filestructure.Measure{
								{
									Symbols: []*filestructure.MusicSymbol{
										{
											Pos:  filestructure.Position{Line: 5, Column: 2},
											Text: "LA_4",
										},
									},
								},
							},
						},
					},
					{
						Data: []byte(`Bagpipe Reader:1.0
"Tune 2 Title",(T,L,0,0,Times NewConverter Roman,11,700,0,0,0,0,0,0)
& B_4 !t
`),
						Tune: &filestructure.Tune{
							Header: &filestructure.TuneHeader{
								Title: "Tune 2 Title",
							},
							Measures: []*filestructure.Measure{
								{
									Symbols: []*filestructure.MusicSymbol{
										{
											Pos:  filestructure.Position{Line: 9, Column: 2},
											Text: "B_4",
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
