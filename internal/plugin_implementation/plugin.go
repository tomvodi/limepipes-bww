package plugin_implementation

import (
	"fmt"
	"github.com/rs/zerolog/log"
	plugininterfaces "github.com/tomvodi/limepipes-plugin-api/plugin/v1/interfaces"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/messages"
	"github.com/tomvodi/limepipes-plugin-bww/internal/common/music_model"
	"github.com/tomvodi/limepipes-plugin-bww/internal/interfaces"
	"os"
)

type plug struct {
	parser       interfaces.BwwParser
	tuneFixer    interfaces.TuneFixer
	fileSplitter interfaces.BwwFileByTuneSplitter
}

func (p *plug) ImportLocalFile(filePath string) (*messages.ImportFileResponse, error) {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	log.Info().Msgf("importing file %s", filePath)

	msg, err := p.importTunesFromData(fileData)
	if err != nil {
		return nil, fmt.Errorf("failed importing tune file %s: %v", filePath, err)
	}

	return msg, nil
}

func (p *plug) Import(fileData []byte) (*messages.ImportFileResponse, error) {
	msg, err := p.importTunesFromData(fileData)
	if err != nil {
		return nil, fmt.Errorf("failed importing file data: %v", err)
	}

	return msg, nil
}

func (p *plug) importTunesFromData(tunesData []byte) (*messages.ImportFileResponse, error) {
	var muModel music_model.MusicModel
	muModel, err := p.parser.ParseBwwData(tunesData)
	if err != nil {
		return nil, fmt.Errorf("failed parsing file: %v", err)
	}

	log.Trace().Msgf("successfully parsed %d tunes",
		len(muModel),
	)

	p.tuneFixer.Fix(muModel)

	bwwFileTuneData, err := p.fileSplitter.SplitFileData(tunesData)
	if err != nil {
		msg := fmt.Sprintf("failed splitting data by tunes: %s", err.Error())
		return nil, fmt.Errorf(msg)
	}

	if len(bwwFileTuneData.TuneTitles()) != len(muModel) {
		log.Error().Msgf("splited bww file and music model don't have the same amount of tunes."+
			" Music model: %d, BWW file: %d", len(muModel), len(bwwFileTuneData.TuneTitles()))
	}

	parsedTunes := make([]*messages.ImportedTune, len(muModel))
	for i, tune := range muModel {
		parsedTunes[i] = &messages.ImportedTune{
			Tune:         tune,
			TuneFileData: bwwFileTuneData.Data(i),
		}
	}

	return &messages.ImportFileResponse{
		ImportedTunes: []*messages.ImportedTune{
			{
				Tune:         nil,
				TuneFileData: nil,
			},
		},
	}, nil
}

func (p *plug) PluginInfo() (*messages.PluginInfoResponse, error) {
	return &messages.PluginInfoResponse{
		Name:        "BWW Plugin",
		Description: "Import Bagpipe Music Writer and Bagpipe Player files.",
		Type:        messages.PluginType_IN,
		FileTypes:   []string{".bww", ".bmw"},
	}, nil
}

func NewPluginImplementation(
	parser interfaces.BwwParser,
	tuneFixer interfaces.TuneFixer,
	fileSplitter interfaces.BwwFileByTuneSplitter,
) plugininterfaces.LimePipesPlugin {
	return &plug{
		parser:       parser,
		tuneFixer:    tuneFixer,
		fileSplitter: fileSplitter,
	}
}
