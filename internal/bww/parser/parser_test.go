package parser

import (
	"github.com/goccy/go-yaml"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/helper"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/measure"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/musicmodel"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/messages"
	"github.com/tomvodi/limepipes-plugin-bww/internal/bww"
	"github.com/tomvodi/limepipes-plugin-bww/internal/bww/symbolmapper"
	"github.com/tomvodi/limepipes-plugin-bww/internal/bwwfile"
	"github.com/tomvodi/limepipes-plugin-bww/internal/interfaces"
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

//nolint:unused
func exportToYaml(muMo musicmodel.MusicModel, filePath string) {
	data, err := yaml.Marshal(muMo)
	Expect(err).ShouldNot(HaveOccurred())
	err = os.WriteFile(filePath, data, 0664)
	Expect(err).ShouldNot(HaveOccurred())
}

func importFromYaml(filePath string) musicmodel.MusicModel {
	muMo := make(musicmodel.MusicModel, 0)
	fileData, err := os.ReadFile(filePath)
	Expect(err).ShouldNot(HaveOccurred())
	err = yaml.Unmarshal(fileData, &muMo)
	Expect(err).ShouldNot(HaveOccurred())

	return muMo
}

func nilAllMeasureMessages(muMo musicmodel.MusicModel) {
	for _, tune := range muMo {
		for _, m := range tune.Measures {
			m.ParserMessages = nil
		}
	}
}

var _ = Describe("BWW Parser", func() {
	utils.SetupConsoleLogger()
	var err error
	var parser interfaces.BwwParser
	var parsedTunes []*messages.ParsedTune
	var musicTunesBww musicmodel.MusicModel
	var musicTunesExpect musicmodel.MusicModel
	var testFile string
	var testFileExpect string

	BeforeEach(func() {
		tok := bwwfile.NewTokenizer()
		tokConv := bwwfile.NewTokenConverter()
		sp := bwwfile.NewStructureParser(
			tok,
			tokConv,
		)
		symmap := symbolmapper.New()
		fsconv := bww.NewConverter(symmap)
		parser = New(sp, fsconv)
		musicTunesBww = make(musicmodel.MusicModel, 0)
	})

	JustBeforeEach(func() {
		data := dataFromFile(testFile)
		parsedTunes, err = parser.ParseBwwData(data)
		for _, pt := range parsedTunes {
			musicTunesBww = append(musicTunesBww, pt.Tune)
		}

		if testFileExpect != "" {
			musicTunesExpect = importFromYaml(testFileExpect)
		}
	})

	When("parsing a file with a staff with 4 measures in it", func() {
		BeforeEach(func() {
			testFile = "./testfiles/four_measures.bww"
		})

		It("should have parsed 4 measures", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).To(HaveLen(1))
			Expect(musicTunesBww[0].Measures).To(HaveLen(4))
		})
	})

	When("having a tune with title, composer, type and footer", func() {
		BeforeEach(func() {
			testFile = "./testfiles/full_tune_header.bww"
		})

		It("should have parsed 4 measures", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).To(HaveLen(1))
			Expect(musicTunesBww[0].Title).To(Equal("Tune Title"))
			Expect(musicTunesBww[0].Composer).To(Equal("Composer"))
			Expect(musicTunesBww[0].Type).To(Equal("Tune Type"))
			Expect(musicTunesBww[0].Footer).To(Equal([]string{"Tune Footer"}))
		})
	})

	When("having all possible time signatures", func() {
		BeforeEach(func() {
			testFile = "./testfiles/time_signatures.bww"
			testFileExpect = "./testfiles/time_signatures.yaml"
		})

		It("should have parsed all time signatures", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(
				BeComparableTo(musicTunesExpect, helper.MusicModelCompareOptions))
		})
	})

	When("having a tune with all kinds of melody notes", func() {
		BeforeEach(func() {
			testFile = "./testfiles/all_melody_notes.bww"
			testFileExpect = "./testfiles/all_melody_notes.yaml"
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).
				Should(BeComparableTo(musicTunesExpect, helper.MusicModelCompareOptions))
		})
	})

	When("having single grace notes following a melody note", func() {
		BeforeEach(func() {
			testFile = "./testfiles/single_graces.bww"
			testFileExpect = "./testfiles/single_graces.yaml"
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).To(Equal(musicTunesExpect))
		})
	})

	When("having only an embellishment without a following melody note", func() {
		BeforeEach(func() {
			testFile = "./testfiles/embellishment_without_following_note.bww"
			testFileExpect = "./testfiles/embellishment_without_following_note.yaml"

		})

		JustBeforeEach(func() {
			//exportToYaml(musicTunesBww, "./testfiles/embellishment_without_following_note.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(
				BeComparableTo(musicTunesExpect, helper.MusicModelCompareOptions))
		})
	})

	When("having dots for the melody note", func() {
		BeforeEach(func() {
			testFile = "./testfiles/dots.bww"
			testFileExpect = "./testfiles/dots.yaml"
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(
				BeComparableTo(musicTunesExpect, helper.MusicModelCompareOptions))
		})
	})

	When("having fermatas for melody notes", func() {
		BeforeEach(func() {
			testFile = "./testfiles/fermatas.bww"
			testFileExpect = "./testfiles/fermatas.yaml"
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(
				BeComparableTo(musicTunesExpect, helper.MusicModelCompareOptions))
		})
	})

	When("having rests", func() {
		BeforeEach(func() {
			testFile = "./testfiles/rests.bww"
			testFileExpect = "./testfiles/rests.yaml"
			//exportToYaml(musicTunesBww, "./testfiles/.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(
				BeComparableTo(musicTunesExpect, helper.MusicModelCompareOptions))
		})
	})

	When("having accidentals", func() {
		BeforeEach(func() {
			testFile = "./testfiles/accidentals.bww"
			testFileExpect = "./testfiles/accidentals.yaml"
		})

		JustBeforeEach(func() {
			//exportToYaml(musicTunesBww, "./testfiles/accidentals.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(
				BeComparableTo(musicTunesExpect, helper.MusicModelCompareOptions))
		})
	})

	When("having doublings", func() {
		BeforeEach(func() {
			testFile = "./testfiles/doublings.bww"
			testFileExpect = "./testfiles/doublings.yaml"
		})

		JustBeforeEach(func() {
			//exportToYaml(musicTunesBww, "./testfiles/doublings.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(
				BeComparableTo(musicTunesExpect, helper.MusicModelCompareOptions))
		})
	})

	When("having grips", func() {
		BeforeEach(func() {
			testFile = "./testfiles/grips.bww"
			testFileExpect = "./testfiles/grips.yaml"
		})

		JustBeforeEach(func() {
			//exportToYaml(musicTunesBww, "./testfiles/grips.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(
				BeComparableTo(musicTunesExpect, helper.MusicModelCompareOptions))
		})
	})

	When("having taorluaths", func() {
		BeforeEach(func() {
			testFile = "./testfiles/taorluaths.bww"
			testFileExpect = "./testfiles/taorluaths.yaml"
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(
				BeComparableTo(musicTunesExpect, helper.MusicModelCompareOptions))
		})
	})

	When("having bubblys", func() {
		BeforeEach(func() {
			testFile = "./testfiles/bubblys.bww"
			testFileExpect = "./testfiles/bubblys.yaml"
		})

		JustBeforeEach(func() {
			//exportToYaml(musicTunesBww, "./testfiles/bubblys.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(
				BeComparableTo(musicTunesExpect, helper.MusicModelCompareOptions))
		})
	})

	When("having throw on d", func() {
		BeforeEach(func() {
			testFile = "./testfiles/throwds.bww"
			testFileExpect = "./testfiles/throwds.yaml"
		})

		JustBeforeEach(func() {
			//exportToYaml(musicTunesBww, "./testfiles/throwds.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(
				BeComparableTo(musicTunesExpect, helper.MusicModelCompareOptions))
		})
	})

	When("having birls", func() {
		BeforeEach(func() {
			testFile = "./testfiles/birls.bww"
			testFileExpect = "./testfiles/birls.yaml"
		})

		JustBeforeEach(func() {
			//exportToYaml(musicTunesBww, "./testfiles/birls.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(
				BeComparableTo(musicTunesExpect, helper.MusicModelCompareOptions))
		})
	})

	When("having strikes", func() {
		BeforeEach(func() {
			testFile = "./testfiles/strikes.bww"
			testFileExpect = "./testfiles/strikes.yaml"
		})

		JustBeforeEach(func() {
			//exportToYaml(musicTunesBww, "./testfiles/strikes.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(
				BeComparableTo(musicTunesExpect, helper.MusicModelCompareOptions))
		})
	})

	When("having peles", func() {
		BeforeEach(func() {
			testFile = "./testfiles/peles.bww"
			testFileExpect = "./testfiles/peles.yaml"
		})

		JustBeforeEach(func() {
			//exportToYaml(musicTunesBww, "./testfiles/peles.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(
				BeComparableTo(musicTunesExpect, helper.MusicModelCompareOptions))
		})
	})

	When("having double strikes", func() {
		BeforeEach(func() {
			testFile = "./testfiles/double_strikes.bww"
			testFileExpect = "./testfiles/double_strikes.yaml"
		})

		JustBeforeEach(func() {
			//exportToYaml(musicTunesBww, "./testfiles/double_strikes.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(
				BeComparableTo(musicTunesExpect, helper.MusicModelCompareOptions))
		})
	})

	When("having triple strikes", func() {
		BeforeEach(func() {
			testFile = "./testfiles/triple_strikes.bww"
			testFileExpect = "./testfiles/triple_strikes.yaml"
		})

		JustBeforeEach(func() {
			//exportToYaml(musicTunesBww, "./testfiles/triple_strikes.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(
				BeComparableTo(musicTunesExpect, helper.MusicModelCompareOptions))
		})
	})

	When("having double graces", func() {
		BeforeEach(func() {
			testFile = "./testfiles/double_grace.bww"
			testFileExpect = "./testfiles/double_grace.yaml"
		})

		JustBeforeEach(func() {
			//exportToYaml(musicTunesBww, "./testfiles/double_grace.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(
				BeComparableTo(musicTunesExpect, helper.MusicModelCompareOptions))
		})
	})

	When("having ties", func() {
		BeforeEach(func() {
			testFile = "./testfiles/ties.bww"
			testFileExpect = "./testfiles/ties.yaml"
		})

		JustBeforeEach(func() {
			//exportToYaml(musicTunesBww, "./testfiles/ties.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(
				BeComparableTo(musicTunesExpect, helper.MusicModelCompareOptions))
		})
	})

	When("having ties in old format with error messages", func() {
		BeforeEach(func() {
			Skip("Parser does not support old format styles")
			testFile = "./testfiles/ties_old_with_errors.bww"
			testFileExpect = "./testfiles/ties_old_with_errors.yaml"
		})

		JustBeforeEach(func() {
			//exportToYaml(musicTunesBww, "./testfiles/ties_old_with_errors.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww[0].Measures[0].ParserMessages[0]).Should(
				Equal(&measure.ParserMessage{
					Symbol:   "^tla",
					Severity: measure.Severity_Warning,
					Text:     "tie in old format (^tla) must follow a note and can't be the first symbol in a measure",
					Fix:      measure.Fix_SkipSymbol,
				}))
			Expect(musicTunesBww[0].Measures[1].ParserMessages[0]).Should(
				Equal(&measure.ParserMessage{
					Symbol:   "^tlg",
					Severity: measure.Severity_Error,
					Text:     "tie in old format (^tlg) must follow a note with pitch and length",
					Fix:      measure.Fix_SkipSymbol,
				}))
			nilAllMeasureMessages(musicTunesBww)
			Expect(musicTunesBww).Should(
				BeComparableTo(musicTunesExpect, helper.MusicModelCompareOptions))
		})
	})

	When("having ties in (old format)", func() {
		BeforeEach(func() {
			Skip("Parser does not support old format styles")
			testFile = "./testfiles/ties_old.bww"
			testFileExpect = "./testfiles/ties_old.yaml"
			//exportToYaml(musicTunesBww, "./testfiles/ties_old.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(
				BeComparableTo(musicTunesExpect, helper.MusicModelCompareOptions))
		})
	})

	When("having irregular groups", func() {
		BeforeEach(func() {
			testFile = "./testfiles/irregular_groups.bww"
			testFileExpect = "./testfiles/irregular_groups.yaml"
		})

		JustBeforeEach(func() {
			//exportToYaml(musicTunesBww, "./testfiles/irregular_groups.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(
				BeComparableTo(musicTunesExpect, helper.MusicModelCompareOptions))
		})
	})

	When("having triplets (old format)", func() {
		BeforeEach(func() {
			Skip("Parser does not support old format styles")
			testFile = "./testfiles/triplets.bww"
			testFileExpect = "./testfiles/triplets.yaml"
			//exportToYaml(musicTunesBww, "./testfiles/triplets.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(
				BeComparableTo(musicTunesExpect, helper.MusicModelCompareOptions))
		})
	})

	When("having time lines", func() {
		BeforeEach(func() {
			testFile = "./testfiles/time_lines.bww"
			testFileExpect = "./testfiles/time_lines.yaml"
		})

		JustBeforeEach(func() {
			//exportToYaml(musicTunesBww, "./testfiles/time_lines.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(
				BeComparableTo(musicTunesExpect, helper.MusicModelCompareOptions))
		})
	})

	When("having space symbols", func() {
		BeforeEach(func() {
			testFile = "./testfiles/space.bww"
			testFileExpect = "./testfiles/space.yaml"
		})

		JustBeforeEach(func() {
			//exportToYaml(musicTunesBww, "./testfiles/space.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(
				BeComparableTo(musicTunesExpect, helper.MusicModelCompareOptions))
		})
	})

	When("having a file with a tune containing inline text and comments", func() {
		BeforeEach(func() {
			testFile = "./testfiles/tune_with_inline_comments.bww"
			testFileExpect = "./testfiles/tune_with_inline_comments.yaml"
		})

		JustBeforeEach(func() {
			//exportToYaml(musicTunesBww, "./testfiles/tune_with_inline_comments.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(
				BeComparableTo(musicTunesExpect, helper.MusicModelCompareOptions))
		})
	})

	When("having a file with two tunes", func() {
		BeforeEach(func() {
			testFile = "./testfiles/two_tunes.bww"
			testFileExpect = "./testfiles/two_tunes.yaml"
		})

		JustBeforeEach(func() {
			//exportToYaml(musicTunesBww, "./testfiles/two_tunes.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(
				BeComparableTo(musicTunesExpect, helper.MusicModelCompareOptions))
		})
	})

	When("having a file with the first tune without a title", func() {
		BeforeEach(func() {
			testFile = "./testfiles/first_tune_no_title.bww"
			testFileExpect = "./testfiles/first_tune_no_title.yaml"
		})

		JustBeforeEach(func() {
			//exportToYaml(musicTunesBww, "./testfiles/first_tune_no_title.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(
				BeComparableTo(musicTunesExpect, helper.MusicModelCompareOptions))
		})
	})

	When("having a tune with no proper staff ending before next staff starts", func() {
		BeforeEach(func() {
			testFile = "./testfiles/tune_with_no_staff_ending_before_next_staff.bww"
			testFileExpect = "./testfiles/tune_with_no_staff_ending_before_next_staff.yaml"
		})

		JustBeforeEach(func() {
			//exportToYaml(musicTunesBww, "./testfiles/tune_with_no_staff_ending_before_next_staff.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(
				BeComparableTo(musicTunesExpect, helper.MusicModelCompareOptions))
		})
	})

	When("having a tune staff that ends with EOF", func() {
		BeforeEach(func() {
			testFile = "./testfiles/tune_staff_ends_with_eof.bww"
			testFileExpect = "./testfiles/tune_staff_ends_with_eof.yaml"
			//exportToYaml(musicTunesBww, "./testfiles/tune_staff_ends_with_eof.yaml")
		})

		JustBeforeEach(func() {
			//exportToYaml(musicTunesBww, "./testfiles/tune_staff_ends_with_eof.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(
				BeComparableTo(musicTunesExpect, helper.MusicModelCompareOptions))
		})
	})

	When("having tune title and config with missing parameter in list", func() {
		BeforeEach(func() {
			testFile = "./testfiles/tune_with_missing_parameter_in_list.bww"
			testFileExpect = "./testfiles/tune_with_missing_parameter_in_list.yaml"
		})

		JustBeforeEach(func() {
			//exportToYaml(musicTunesBww, "./testfiles/tune_with_missing_parameter_in_list.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(
				BeComparableTo(musicTunesExpect, helper.MusicModelCompareOptions))
		})
	})

	When("having file with multiple bagpipe reader version definitions", func() {
		BeforeEach(func() {
			testFile = "./testfiles/tune_with_multiple_bagpipe_reader_version_definitions.bww"
			testFileExpect = "./testfiles/tune_with_multiple_bagpipe_reader_version_definitions.yaml"
			//exportToYaml(musicTunesBww, "./testfiles/tune_with_multiple_bagpipe_reader_version_definitions.yaml")
		})

		JustBeforeEach(func() {
			//exportToYaml(musicTunesBww, "./testfiles/tune_with_multiple_bagpipe_reader_version_definitions.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(
				BeComparableTo(musicTunesExpect, helper.MusicModelCompareOptions))
		})
	})

	When("having tune with symbol and measure comments", func() {
		BeforeEach(func() {
			testFile = "./testfiles/tune_with_symbol_comment.bww"
			testFileExpect = "./testfiles/tune_with_symbol_comment.yaml"
		})

		JustBeforeEach(func() {
			//exportToYaml(musicTunesBww, "./testfiles/tune_with_symbol_comment.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(
				BeComparableTo(musicTunesExpect, helper.MusicModelCompareOptions))
		})
	})

	When("having tune with time line end after staff end", func() {
		BeforeEach(func() {
			testFile = "./testfiles/tune_with_time_line_end_after_staff_end.bww"
			testFileExpect = "./testfiles/tune_with_time_line_end_after_staff_end.yaml"
		})

		JustBeforeEach(func() {
			//exportToYaml(musicTunesBww, "./testfiles/tune_with_time_line_end_after_staff_end.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("staff end token is not at the end of line 5"))
		})
	})

	When("having tune with inline tune tempo", func() {
		BeforeEach(func() {
			testFile = "./testfiles/tunetempo_inline.bww"
			testFileExpect = "./testfiles/tunetempo_inline.yaml"
		})

		JustBeforeEach(func() {
			//exportToYaml(musicTunesBww, "./testfiles/tunetempo_inline.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(
				BeComparableTo(musicTunesExpect, helper.MusicModelCompareOptions))
		})
	})

	When("having file with piobairached throws and doublings in it", func() {
		BeforeEach(func() {
			testFile = "./testfiles/pio_throws_and_doublings.bww"
			testFileExpect = "./testfiles/pio_throws_and_doublings.yaml"
		})

		JustBeforeEach(func() {
			//exportToYaml(musicTunesBww, "./testfiles/pio_throws_and_doublings.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(
				BeComparableTo(musicTunesExpect, helper.MusicModelCompareOptions))
		})
	})

	When("having file with all cadences in it", func() {
		BeforeEach(func() {
			Skip("Piobairached is not supported")
			testFile = "./testfiles/cadences.bww"
			testFileExpect = "./testfiles/cadences.yaml"
			//exportToYaml(musicTunesBww, "./testfiles/cadences.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(
				BeComparableTo(musicTunesExpect, helper.MusicModelCompareOptions))
		})
	})

	When("having file with piobairached grips in it", func() {
		BeforeEach(func() {
			Skip("Piobairached is not supported")
			testFile = "./testfiles/pio_grips.bww"
			testFileExpect = "./testfiles/pio_grips.yaml"
			//exportToYaml(musicTunesBww, "./testfiles/pio_grips.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(
				BeComparableTo(musicTunesExpect, helper.MusicModelCompareOptions))
		})
	})

	When("having file with piobairached echo beats in it", func() {
		BeforeEach(func() {
			Skip("Piobairached is not supported")
			testFile = "./testfiles/pio_echo_beats.bww"
			testFileExpect = "./testfiles/pio_echo_beats.yaml"
			//exportToYaml(musicTunesBww, "./testfiles/pio_echo_beats.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(
				BeComparableTo(musicTunesExpect, helper.MusicModelCompareOptions))
		})
	})

	When("having file with piobairached darodos in it", func() {
		BeforeEach(func() {
			Skip("Piobairached is not supported")
			testFile = "./testfiles/pio_darodo.bww"
			testFileExpect = "./testfiles/pio_darodo.yaml"
			//exportToYaml(musicTunesBww, "./testfiles/pio_darodo.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(
				BeComparableTo(musicTunesExpect, helper.MusicModelCompareOptions))
		})
	})

	When("having file with piobairached lemluaths in it", func() {
		BeforeEach(func() {
			Skip("Piobairached is not supported")
			testFile = "./testfiles/pio_lemluaths.bww"
			testFileExpect = "./testfiles/pio_lemluaths.yaml"
			//exportToYaml(musicTunesBww, "./testfiles/pio_lemluaths.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(
				BeComparableTo(musicTunesExpect, helper.MusicModelCompareOptions))
		})
	})

	When("having file with piobairached taorluaths in it", func() {
		BeforeEach(func() {
			Skip("Piobairached is not supported")
			testFile = "./testfiles/pio_taorluaths.bww"
			testFileExpect = "./testfiles/pio_taorluaths.yaml"
			//exportToYaml(musicTunesBww, "./testfiles/pio_taorluaths.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(
				BeComparableTo(musicTunesExpect, helper.MusicModelCompareOptions))
		})
	})

	When("having file with piobairached crunluaths in it", func() {
		BeforeEach(func() {
			Skip("Piobairached is not supported")
			testFile = "./testfiles/pio_crunluaths.bww"
			testFileExpect = "./testfiles/pio_crunluaths.yaml"
			//exportToYaml(musicTunesBww, "./testfiles/pio_crunluaths.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(
				BeComparableTo(musicTunesExpect, helper.MusicModelCompareOptions))
		})
	})

	When("having file with piobairached triplings in it", func() {
		BeforeEach(func() {
			Skip("Piobairached is not supported")
			testFile = "./testfiles/pio_triplings.bww"
			testFileExpect = "./testfiles/pio_triplings.yaml"
			//exportToYaml(musicTunesBww, "./testfiles/pio_triplings.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(
				BeComparableTo(musicTunesExpect, helper.MusicModelCompareOptions))
		})
	})

	When("having file with misc movements in it", func() {
		BeforeEach(func() {
			Skip("Piobairached is not supported")
			testFile = "./testfiles/pio_misc.bww"
			testFileExpect = "./testfiles/pio_misc.yaml"
			//exportToYaml(musicTunesBww, "./testfiles/pio_misc.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(
				BeComparableTo(musicTunesExpect, helper.MusicModelCompareOptions))
		})
	})

	When("having file with segno and dalsegno", func() {
		BeforeEach(func() {
			Skip("Piobairached is not supported")
			testFile = "./testfiles/segno_dalsegno.bww"
			testFileExpect = "./testfiles/segno_dalsegno.yaml"
			//exportToYaml(musicTunesBww, "./testfiles/segno_dalsegno.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(
				BeComparableTo(musicTunesExpect, helper.MusicModelCompareOptions))
		})
	})

	When("having file with fine and dacapoalfine", func() {
		BeforeEach(func() {
			Skip("Piobairached is not supported")
			testFile = "./testfiles/fine_dacapoalfine.bww"
			testFileExpect = "./testfiles/fine_dacapoalfine.yaml"
			//exportToYaml(musicTunesBww, "./testfiles/fine_dacapoalfine.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(
				BeComparableTo(musicTunesExpect, helper.MusicModelCompareOptions))
		})
	})

	When("having inline comment shouldn't remove measures", func() {
		BeforeEach(func() {
			testFile = "./testfiles/inline_comment_removes_first_staff_measures.bww"
			testFileExpect = "./testfiles/inline_comment_removes_first_staff_measures.yaml"
		})

		JustBeforeEach(func() {
			//exportToYaml(musicTunesBww, "./testfiles/inline_comment_removes_first_staff_measures.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(
				BeComparableTo(musicTunesExpect, helper.MusicModelCompareOptions))
		})
	})

	When("having a tune with repeats", func() {
		BeforeEach(func() {
			testFile = "./testfiles/tune_with_repeats.bww"
			testFileExpect = "./testfiles/tune_with_repeats.yaml"
			//exportToYaml(musicTunesBww, "./testfiles/tune_with_repeats.yaml")
		})

		It("should have parsed file correctly", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).Should(
				BeComparableTo(musicTunesExpect, helper.MusicModelCompareOptions))
		})
	})

	When("parsing the file with all bww symbols in it", func() {
		BeforeEach(func() {
			testFile = "./testfiles/all_symbols.bww"
		})

		It("should succeed", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).To(HaveLen(2))
		})
	})

	When("parsing the file with all piobaireached symbols in it", func() {
		BeforeEach(func() {
			Skip("Piobairached is not supported")
			testFile = "./testfiles/all_piobaireached_symbols.bww"
		})

		It("should succeed", func() {
			Expect(err).ShouldNot(HaveOccurred())
			Expect(musicTunesBww).To(HaveLen(11))
		})
	})
})
