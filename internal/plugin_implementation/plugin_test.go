package plugin_implementation

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/helper"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/measure"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/musicmodel"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/messages"
	"github.com/tomvodi/limepipes-plugin-bww/internal/common"
	"github.com/tomvodi/limepipes-plugin-bww/internal/interfaces/mocks"
)

var _ = Describe("Import tunes", func() {
	var lpPlug *plug
	var parser *mocks.BwwParser
	var tuneFixer *mocks.TuneFixer
	var fileSplitter *mocks.BwwFileByTuneSplitter
	var testMusicModel musicmodel.MusicModel
	var tune1FileData []byte
	var tuneData []byte
	var err error
	var parsedTunes []*messages.ParsedTune

	BeforeEach(func() {
		parser = mocks.NewBwwParser(GinkgoT())
		tuneFixer = mocks.NewTuneFixer(GinkgoT())
		fileSplitter = mocks.NewBwwFileByTuneSplitter(GinkgoT())
		lpPlug = &plug{
			parser:       parser,
			tuneFixer:    tuneFixer,
			fileSplitter: fileSplitter,
		}
		tuneData = []byte("tune data")
		tune1FileData = []byte("tune 1 data")
	})

	JustBeforeEach(func() {
		parsedTunes, err = lpPlug.Parse(tuneData)
	})

	Context("parser returns an error", func() {
		BeforeEach(func() {
			parser.EXPECT().ParseBwwData(mock.Anything).
				Return(nil, fmt.Errorf("failed parsing"))
		})

		When("importing a tune data", func() {
			It("should return an error", func() {
				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Context("Having a tune returned by the parser", func() {
		BeforeEach(func() {
			testMusicModel = musicmodel.MusicModel{
				{
					Title: "test tune",
					Measures: []*measure.Measure{
						{
							Comments: []string{"comment"},
						},
					},
				},
			}
			parser.EXPECT().ParseBwwData(mock.Anything).
				Return(testMusicModel, nil)
			tuneFixer.EXPECT().Fix(testMusicModel)
		})

		When("there is an error when splitting the file by tunes", func() {
			BeforeEach(func() {
				fileSplitter.EXPECT().SplitFileData(tuneData).
					Return(nil, fmt.Errorf("failed splitting"))
			})

			It("should return an error", func() {
				Expect(err).Should(HaveOccurred())
			})
		})

		When("there is a difference between parsed tunes and found tunes in file", func() {
			BeforeEach(func() {
				tuneFileData := &common.BwwFileTuneData{}
				tuneFileData.AddTuneData("tune 1", []byte("tune 1 data"))
				tuneFileData.AddTuneData("tune 2", []byte("tune 2 data"))
				fileSplitter.EXPECT().SplitFileData(tuneData).
					Return(tuneFileData, nil)
			})

			It("should return an error", func() {
				Expect(err).Should(HaveOccurred())
			})
		})

		When("successfully parsed tune data", func() {
			BeforeEach(func() {
				tuneFileData := &common.BwwFileTuneData{}
				tuneFileData.AddTuneData("tune 1", tune1FileData)
				fileSplitter.EXPECT().SplitFileData(tuneData).
					Return(tuneFileData, nil)
			})

			It("should succeed", func() {
				Expect(err).ShouldNot(HaveOccurred())
				Expect(parsedTunes).Should(BeComparableTo(
					[]*messages.ParsedTune{
						{
							Tune:         testMusicModel[0],
							TuneFileData: tune1FileData,
						},
					},
					helper.MusicModelCompareOptions))
			})
		})
	})
})
