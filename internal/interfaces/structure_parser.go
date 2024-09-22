package interfaces

import (
	"github.com/tomvodi/limepipes-plugin-bww/internal/structure"
)

type StructureParser interface {
	ParseDocumentStructure(data []byte) (*structure.BwwFile, error)
}
