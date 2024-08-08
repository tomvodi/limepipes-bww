package interfaces

import "github.com/tomvodi/limepipes-plugin-bww/internal/common/music_model"

// TuneFixer fixes a tune's meta data like type composer and arranger
type TuneFixer interface {
	Fix(muMo music_model.MusicModel)
}
