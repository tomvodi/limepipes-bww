package interfaces

import (
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/tune"
	"github.com/tomvodi/limepipes-plugin-bww/internal/filestructure"
)

type StructureToModelConverter interface {
	Convert(t *filestructure.Tune) (*tune.Tune, error)
}
