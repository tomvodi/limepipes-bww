package interfaces

import (
	"github.com/tomvodi/limepipes-plugin-bww/internal/filestructure"
)

type StructureParser interface {
	ParseDocumentStructure(data []byte) (*filestructure.BwwFile, error)
}
