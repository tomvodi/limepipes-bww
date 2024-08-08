package interfaces

import (
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/musicmodel"
)

type BwwParser interface {
	ParseBwwData(data []byte) (musicmodel.MusicModel, error)
}
