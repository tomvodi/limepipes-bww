package symbolmapper

import "github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols/embellishment"

func init() {
	str := newEmbellishment(embellishment.Type_TripleStrike)
	for _, s := range lowPitchesLgToHA {
		symbolsMap["st3"+s] = str
	}
	symbolsMap["lst3d"] = newEmbellishment(
		embellishment.Type_TripleStrike,
		embellishment.Weight_Light,
	)

	str = newEmbellishment(embellishment.Type_TripleStrike, embellishment.Variant_G)
	for _, s := range lowPitchesLaToF {
		symbolsMap["gst3"+s] = str
	}
	symbolsMap["lgst3d"] = newEmbellishment(
		embellishment.Type_TripleStrike,
		embellishment.Weight_Light,
		embellishment.Variant_G,
	)

	str = newEmbellishment(embellishment.Type_TripleStrike, embellishment.Variant_Thumb)
	for _, s := range lowPitchesLaToHG {
		symbolsMap["tst3"+s] = str
	}
	symbolsMap["ltst3d"] = newEmbellishment(
		embellishment.Type_TripleStrike,
		embellishment.Weight_Light,
		embellishment.Variant_Thumb,
	)

	str = newEmbellishment(embellishment.Type_TripleStrike, embellishment.Variant_Half)
	for _, s := range lowPitchesLaToHA {
		symbolsMap["hst3"+s] = str
	}
	symbolsMap["lhst3d"] = newEmbellishment(
		embellishment.Type_TripleStrike,
		embellishment.Weight_Light,
		embellishment.Variant_Half,
	)
}
