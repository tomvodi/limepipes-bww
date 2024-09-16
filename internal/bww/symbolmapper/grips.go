package symbolmapper

import (
	"fmt"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols/embellishment"
)

func init() {
	grp := &symbols.Symbol{
		Note: &symbols.Note{
			Embellishment: &embellishment.Embellishment{
				Type: embellishment.Type_Grip,
			},
		},
	}
	symbolsMap["grp"] = grp
	symbolsMap["hgrp"] = grp
	symbolsMap["grpb"] = grp

	ggrp := &symbols.Symbol{
		Note: &symbols.Note{
			Embellishment: &embellishment.Embellishment{
				Type:    embellishment.Type_Grip,
				Variant: embellishment.Variant_G,
			},
		},
	}

	for _, p := range lowPitchesLgToF {
		symbolsMap[fmt.Sprintf("ggrp%s", p)] = ggrp
	}
	symbolsMap["ggrpdb"] = ggrp

	tgrp := &symbols.Symbol{
		Note: &symbols.Note{
			Embellishment: &embellishment.Embellishment{
				Type:    embellishment.Type_Grip,
				Variant: embellishment.Variant_Thumb,
			},
		},
	}

	for _, p := range lowPitchesLgToHG {
		symbolsMap[fmt.Sprintf("tgrp%s", p)] = tgrp
	}
	symbolsMap["tgrpdb"] = tgrp

	hgrp := &symbols.Symbol{
		Note: &symbols.Note{
			Embellishment: &embellishment.Embellishment{
				Type:    embellishment.Type_Grip,
				Variant: embellishment.Variant_Half,
			},
		},
	}
	for _, p := range lowPitches {
		symbolsMap[fmt.Sprintf("hgrp%s", p)] = hgrp
	}
	symbolsMap["hgrpdb"] = hgrp
}
