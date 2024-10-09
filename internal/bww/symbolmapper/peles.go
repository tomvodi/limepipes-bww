package symbolmapper

import (
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols/embellishment"
)

func init() {
	pel := newEmbellishment(embellishment.Type_Pele)
	for _, s := range lowPitchesLaToF {
		symbolsMap["pel"+s] = pel
	}
	symbolsMap["lpeld"] = newEmbellishment(
		embellishment.Type_Pele,
		embellishment.Weight_Light,
	)

	pel = newEmbellishment(embellishment.Type_Pele, embellishment.Variant_Thumb)
	for _, s := range lowPitchesLaToHG {
		symbolsMap["tpel"+s] = pel
	}
	symbolsMap["ltpeld"] = newEmbellishment(
		embellishment.Type_Pele,
		embellishment.Variant_Thumb,
		embellishment.Weight_Light,
	)

	pel = newEmbellishment(embellishment.Type_Pele, embellishment.Variant_Half)
	for _, s := range lowPitchesLaToHG {
		symbolsMap["hpel"+s] = pel
	}
	symbolsMap["lhpeld"] = newEmbellishment(
		embellishment.Type_Pele,
		embellishment.Variant_Half,
		embellishment.Weight_Light,
	)
}
