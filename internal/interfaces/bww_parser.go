package interfaces

import (
	"github.com/tomvodi/limepipes-plugin-bww/internal/common/music_model"
)

type BwwParser interface {
	ParseBwwData(data []byte) (music_model.MusicModel, error)
}
