package symbolmapper

import (
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols/embellishment"
)

func init() {
	bubble := newEmbellishment(embellishment.Type_Bubbly)
	symbolsMap["bubly"] = bubble
	symbolsMap["hbubly"] = bubble
}
