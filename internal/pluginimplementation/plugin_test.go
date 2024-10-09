package pluginimplementation

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/mock"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/helper"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/measure"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/tune"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/fileformat"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/messages"
	"github.com/tomvodi/limepipes-plugin-bww/internal/interfaces/mocks"
)

var _ = Describe("PluginInfo", func() {
	var lpPlug *Plugin
	var err error
	var pluginInfo *messages.PluginInfoResponse

	BeforeEach(func() {
		lpPlug = &Plugin{}
	})

	JustBeforeEach(func() {
		pluginInfo, err = lpPlug.PluginInfo()
	})

	It("should return plugin info", func() {
		Expect(err).ShouldNot(HaveOccurred())
		Expect(pluginInfo).Should(Equal(&messages.PluginInfoResponse{
			Name:           "BWW Plugin",
			Description:    "Import Bagpipe Music Writer and Bagpipe Player files.",
			FileFormat:     fileformat.Format_BWW,
			Type:           messages.PluginType_IN,
			FileExtensions: []string{".bww", ".bmw"},
		}))
	})
})

var _ = Describe("Parse", func() {
	var lpPlug *Plugin
	var parser *mocks.BwwParser
	var tuneFixer *mocks.TuneFixer
	var testParsedTunes []*messages.ParsedTune
	var tuneData []byte
	var afs afero.Fs
	var err error
	var parsedTunes []*messages.ParsedTune

	BeforeEach(func() {
		parser = mocks.NewBwwParser(GinkgoT())
		tuneFixer = mocks.NewTuneFixer(GinkgoT())
		afs = afero.NewMemMapFs()
		lpPlug = &Plugin{
			afs:       afs,
			parser:    parser,
			tuneFixer: tuneFixer,
		}
		tuneData = []byte("tune data")
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
	})

	Context("parsing from file", func() {
		var filePath string

		JustBeforeEach(func() {
			parsedTunes, err = lpPlug.ParseFromFile(filePath)
		})

		When("file does not exist", func() {
			BeforeEach(func() {
				filePath = "test.bww"
			})

			It("should return an error", func() {
				Expect(err).Should(HaveOccurred())
			})
		})

		When("file exists", func() {
			BeforeEach(func() {
				filePath = "test.bww"
				err = afero.WriteFile(afs, filePath, tuneData, 0644)
				Expect(err).ShouldNot(HaveOccurred())
			})

			Context("parser returns an error", func() {
				BeforeEach(func() {
					parser.EXPECT().ParseBwwData(tuneData).
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
	})

	Context("parsing from data", func() {
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
})

var _ = Describe("Export", func() {
	var err error
	var lpPlug *Plugin
	var exportTunes []*tune.Tune

	BeforeEach(func() {
		lpPlug = &Plugin{}
	})

	JustBeforeEach(func() {
		_, err = lpPlug.Export(exportTunes)
	})

	When("exporting a tune", func() {
		BeforeEach(func() {
			exportTunes = []*tune.Tune{
				{
					Title: "test tune",
					Measures: []*measure.Measure{
						{
							Comments: []string{"comment"},
						},
					},
				},
			}
		})

		It("should return an error", func() {
			Expect(err).Should(HaveOccurred())
		})
	})
})

var _ = Describe("ExportToFile", func() {
	var err error
	var lpPlug *Plugin
	var exportTunes []*tune.Tune
	var exportPath string

	BeforeEach(func() {
		lpPlug = &Plugin{}
	})

	JustBeforeEach(func() {
		err = lpPlug.ExportToFile(exportTunes, exportPath)
	})

	When("exporting a tune", func() {
		BeforeEach(func() {
			exportTunes = []*tune.Tune{
				{
					Title: "test tune",
					Measures: []*measure.Measure{
						{
							Comments: []string{"comment"},
						},
					},
				},
			}
			exportPath = "test.bww"
		})

		It("should return an error", func() {
			Expect(err).Should(HaveOccurred())
		})
	})
})
