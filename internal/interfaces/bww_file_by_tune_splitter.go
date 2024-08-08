package interfaces

import (
	"github.com/tomvodi/limepipes-plugin-bww/internal/common"
)

type BwwFileByTuneSplitter interface {
	SplitFileData(data []byte) (*common.BwwFileTuneData, error)
}
