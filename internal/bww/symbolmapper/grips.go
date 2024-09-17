package symbolmapper

import (
	"fmt"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols/embellishment"
)

func init() {
	grp := newEmbellishment(embellishment.Type_Grip)
	symbolsMap["grp"] = grp
	symbolsMap["hgrp"] = grp
	symbolsMap["grpb"] = grp

	ggrp := newEmbellishment(embellishment.Type_Grip, embellishment.Variant_G)

	for _, p := range lowPitchesLgToF {
		symbolsMap[fmt.Sprintf("ggrp%s", p)] = ggrp
	}
	symbolsMap["ggrpdb"] = ggrp

	tgrp := newEmbellishment(embellishment.Type_Grip, embellishment.Variant_Thumb)

	for _, p := range lowPitchesLgToHG {
		symbolsMap[fmt.Sprintf("tgrp%s", p)] = tgrp
	}
	symbolsMap["tgrpdb"] = tgrp

	hgrp := newEmbellishment(embellishment.Type_Grip, embellishment.Variant_Half)
	for _, p := range lowPitchesLgToHA {
		symbolsMap[fmt.Sprintf("hgrp%s", p)] = hgrp
	}
	symbolsMap["hgrpdb"] = hgrp
}
