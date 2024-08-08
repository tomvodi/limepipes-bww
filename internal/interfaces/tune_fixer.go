package interfaces

import (
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/musicmodel"
)

// TuneFixer fixes a tune's meta data like type composer and arranger
type TuneFixer interface {
	Fix(muMo musicmodel.MusicModel)
}
