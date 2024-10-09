package symbolmapper

import "github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols/embellishment"

func init() {
	str := newEmbellishment(embellishment.Type_DoubleStrike)
	for _, s := range lowPitchesLgToHA {
		symbolsMap["st2"+s] = str
	}
	symbolsMap["lst2d"] = newEmbellishment(
		embellishment.Type_DoubleStrike,
		embellishment.Weight_Light,
	)

	str = newEmbellishment(embellishment.Type_DoubleStrike, embellishment.Variant_G)
	for _, s := range lowPitchesLaToF {
		symbolsMap["gst2"+s] = str
	}
	symbolsMap["lgst2d"] = newEmbellishment(
		embellishment.Type_DoubleStrike,
		embellishment.Weight_Light,
		embellishment.Variant_G,
	)

	str = newEmbellishment(embellishment.Type_DoubleStrike, embellishment.Variant_Thumb)
	for _, s := range lowPitchesLaToHG {
		symbolsMap["tst2"+s] = str
	}
	symbolsMap["ltst2d"] = newEmbellishment(
		embellishment.Type_DoubleStrike,
		embellishment.Weight_Light,
		embellishment.Variant_Thumb,
	)

	str = newEmbellishment(embellishment.Type_DoubleStrike, embellishment.Variant_Half)
	for _, s := range lowPitchesLaToHA {
		symbolsMap["hst2"+s] = str
	}
	symbolsMap["lhst2d"] = newEmbellishment(
		embellishment.Type_DoubleStrike,
		embellishment.Weight_Light,
		embellishment.Variant_Half,
	)
}
