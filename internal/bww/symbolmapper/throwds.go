package symbolmapper

import (
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols/embellishment"
)

func init() {
	thr := newEmbellishment(embellishment.Type_ThrowD)
	symbolsMap["hvthrd"] = thr
	symbolsMap["hhvthrd"] = thr
	thr = newEmbellishment(embellishment.Type_ThrowD, embellishment.Weight_Light)
	symbolsMap["thrd"] = thr
	symbolsMap["hthrd"] = thr
}
