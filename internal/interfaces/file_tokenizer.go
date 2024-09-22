package interfaces

import "github.com/tomvodi/limepipes-plugin-bww/internal/common"

type FileTokenizer interface {
	Tokenize(data []byte) ([]*common.Token, error)
}
