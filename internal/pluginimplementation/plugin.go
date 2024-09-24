package pluginimplementation

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/tune"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/fileformat"
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/messages"
	"github.com/tomvodi/limepipes-plugin-bww/internal/interfaces"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
)

type Plugin struct {
	parser    interfaces.BwwParser
	tuneFixer interfaces.TuneFixer
}

func (p *Plugin) PluginInfo() (*messages.PluginInfoResponse, error) {
	return &messages.PluginInfoResponse{
		Name:           "BWW Plugin",
		Description:    "Import Bagpipe Music Writer and Bagpipe Player files.",
		FileFormat:     fileformat.Format_BWW,
		Type:           messages.PluginType_IN,
		FileExtensions: []string{".bww", ".bmw"},
	}, nil
}

func (p *Plugin) ExportToFile(
	[]*tune.Tune,
	string,
) error {
	return status.Error(codes.Unimplemented, "ExportToFile not implemented")
}

func (p *Plugin) Export(
	[]*tune.Tune,
) ([]byte, error) {
	return nil, status.Error(codes.Unimplemented, "Export not implemented")
}

func (p *Plugin) ParseFromFile(
	filePath string,
) ([]*messages.ParsedTune, error) {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	log.Info().Msgf("importing file %s", filePath)

	msg, err := p.parseTunesFromData(fileData)
	if err != nil {
		return nil, fmt.Errorf("failed importing tune file %s: %v", filePath, err)
	}

	return msg, nil
}

func (p *Plugin) Parse(
	data []byte,
) ([]*messages.ParsedTune, error) {
	return p.parseTunesFromData(data)
}

func (p *Plugin) parseTunesFromData(tunesData []byte) ([]*messages.ParsedTune, error) {
	parsedTunes, err := p.parser.ParseBwwData(tunesData)
	if err != nil {
		return nil, fmt.Errorf("failed parsing t data: %v", err)
	}

	p.tuneFixer.Fix(parsedTunes)

	return parsedTunes, nil
}

func NewPluginImplementation(
	parser interfaces.BwwParser,
	tuneFixer interfaces.TuneFixer,
) *Plugin {
	return &Plugin{
		parser:    parser,
		tuneFixer: tuneFixer,
	}
}
