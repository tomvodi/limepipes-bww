package interfaces

import (
	"github.com/tomvodi/limepipes-plugin-bww/internal/common"
	"github.com/tomvodi/limepipes-plugin-bww/internal/filestructure"
)

type TokenStructureConverter interface {
	Convert(tokens []*common.Token) (*filestructure.BwwFile, error)
}
