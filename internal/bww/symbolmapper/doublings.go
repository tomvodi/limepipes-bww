package symbolmapper

import (
	"fmt"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols/embellishment"
)

func init() {
	for _, p := range lowPitchesLgToHA {
		symbolsMap[fmt.Sprintf("db%s", p)] =
			newEmbellishment(embellishment.Type_Doubling)

		for _, p := range lowPitchesLgToF {
			symbolsMap[fmt.Sprintf("tdb%s", p)] =
				newEmbellishment(embellishment.Type_Doubling, embellishment.Variant_Thumb)
			symbolsMap[fmt.Sprintf("hdb%s", p)] =
				newEmbellishment(embellishment.Type_Doubling, embellishment.Variant_Half)
		}
	}
}
