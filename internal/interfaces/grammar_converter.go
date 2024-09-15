package interfaces

import (
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/musicmodel"
	"github.com/tomvodi/limepipes-plugin-bww/internal/common"
)

type GrammarConverter interface {
	Convert(grammar *common.BwwStructure) (musicmodel.MusicModel, error)
}
