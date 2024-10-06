package pluginimplementation

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/helper"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/measure"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/tune"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/messages"
	"github.com/tomvodi/limepipes-plugin-bww/internal/interfaces/mocks"
)

var _ = Describe("Import tunes", func() {
	var lpPlug *Plugin
	var parser *mocks.BwwParser
	var tuneFixer *mocks.TuneFixer
	var testParsedTunes []*messages.ParsedTune
	var tuneData []byte
	var err error
	var parsedTunes []*messages.ParsedTune

	BeforeEach(func() {
		parser = mocks.NewBwwParser(GinkgoT())
		tuneFixer = mocks.NewTuneFixer(GinkgoT())
		lpPlug = &Plugin{
			parser:    parser,
			tuneFixer: tuneFixer,
		}
		tuneData = []byte("tune data")
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
			testParsedTunes = []*messages.ParsedTune{
				{
					Tune: &tune.Tune{
						Title: "test tune",
						Measures: []*measure.Measure{
							{
								Comments: []string{"comment"},
							},
						},
					},
					TuneFileData: []byte("tune file data"),
				},
			}
			parser.EXPECT().ParseBwwData(mock.Anything).
				Return(testParsedTunes, nil)
			tuneFixer.EXPECT().Fix(testParsedTunes)
		})

		When("successfully parsed tune data", func() {
			It("should succeed", func() {
				Expect(err).ShouldNot(HaveOccurred())
				Expect(parsedTunes).Should(BeComparableTo(
					testParsedTunes,
					helper.MusicModelCompareOptions))
			})
		})
	})
})
