package interfaces

import (
	"github.com/tomvodi/limepipes-plugin-bww/internal/common"
	"github.com/tomvodi/limepipes-plugin-bww/internal/structure"
)

type TokenStructureConverter interface {
	Convert(tokens []*common.Token) (*structure.BwwFile, error)
}
