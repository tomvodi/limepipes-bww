package symbolmapper

import "github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols/embellishment"

func init() {
	str := newEmbellishment(embellishment.Type_Strike)
	for _, s := range lowPitchesLgToHG {
		symbolsMap["str"+s] = str
	}

	str = newEmbellishment(embellishment.Type_Strike, embellishment.Variant_G)
	for _, s := range lowPitchesLaToF {
		symbolsMap["gst"+s] = str
	}
	symbolsMap["lgstd"] = newEmbellishment(
		embellishment.Type_Strike,
		embellishment.Variant_G,
		embellishment.Weight_Light,
	)

	str = newEmbellishment(embellishment.Type_Strike, embellishment.Variant_Thumb)
	for _, s := range lowPitchesLaToHG {
		symbolsMap["tst"+s] = str
	}
	symbolsMap["ltstd"] = newEmbellishment(
		embellishment.Type_Strike,
		embellishment.Variant_Thumb,
		embellishment.Weight_Light,
	)

	str = newEmbellishment(embellishment.Type_Strike, embellishment.Variant_Half)
	for _, s := range lowPitchesLaToHG {
		symbolsMap["hst"+s] = str
	}
	symbolsMap["lhstd"] = newEmbellishment(
		embellishment.Type_Strike,
		embellishment.Variant_Half,
		embellishment.Weight_Light,
	)
}
