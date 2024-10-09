package interfaces

import (
	"github.com/tomvodi/limepipes-plugin-api/plugin/v1/messages"
)

// TuneFixer fixes a tune's meta data like type composer and arranger
type TuneFixer interface {
	Fix(parsedTunes []*messages.ParsedTune)
}
