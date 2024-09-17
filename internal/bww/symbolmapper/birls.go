package symbolmapper

import (
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols/embellishment"
)

func init() {
	symbolsMap["brl"] = newEmbellishment(embellishment.Type_Birl)
	symbolsMap["abr"] = newEmbellishment(embellishment.Type_ABirl)
	symbolsMap["gbr"] = newEmbellishment(embellishment.Type_GraceBirl)
	symbolsMap["tbr"] = newEmbellishment(embellishment.Type_GraceBirl)
}
