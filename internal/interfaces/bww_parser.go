package interfaces

import (
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/messages"
)

type BwwParser interface {
	ParseBwwData(data []byte) ([]*messages.ParsedTune, error)
}
