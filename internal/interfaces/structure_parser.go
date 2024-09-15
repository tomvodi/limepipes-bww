package interfaces

import (
	"github.com/tomvodi/limepipes-plugin-bww/internal/common"
)

type StructureParser interface {
	ParseDocumentStructure(data []byte) (*common.BwwStructure, error)
}
