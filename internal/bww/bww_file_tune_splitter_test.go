package bww

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tomvodi/limepipes-plugin-bww/internal/common"
	"os"
)

var _ = Describe("BwwFileTuneSplitter", func() {
	var fileData []byte
	var err error
	var tuneData *common.BwwFileTuneData
	var splitter *FileSplitter

	BeforeEach(func() {
		splitter = &FileSplitter{}
	})

	JustBeforeEach(func() {
		tuneData, err = splitter.SplitFileData(fileData)
	})

	When("having a file with two tunes in it", func() {
		BeforeEach(func() {
			fileData, err = os.ReadFile("./testfiles/two_tunes.bww")
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should return the tuneData for both files", func() {
			Expect(tuneData.TuneTitles()).
				To(Equal([]string{
					"Tune 1 Title",
					"Tune 2 Title",
				}))
		})
	})

	When("having a file with three tunes in it, where two tunes are the same", func() {
		BeforeEach(func() {
			fileData, err = os.ReadFile("./testfiles/first_tune_no_title.bww")
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should return the tuneData for one no name tune", func() {
			Expect(tuneData.TuneTitles()).
				To(Equal([]string{
					"no name",
				}))
		})
	})

	When("having a file with a tune in it that has no title", func() {

	})
})
