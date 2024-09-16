package symbolmapper

import (
	"fmt"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols/embellishment"
)

func init() {
	for _, p := range lowPitches {
		symbolsMap[fmt.Sprintf("db%s", p)] = &symbols.Symbol{
			Note: &symbols.Note{
				Embellishment: &embellishment.Embellishment{
					Type: embellishment.Type_Doubling,
				},
			},
		}
	}

	for _, p := range lowPitchesLgToF {
		symbolsMap[fmt.Sprintf("tdb%s", p)] = &symbols.Symbol{
			Note: &symbols.Note{
				Embellishment: &embellishment.Embellishment{
					Type:    embellishment.Type_Doubling,
					Variant: embellishment.Variant_Thumb,
				},
			},
		}
		symbolsMap[fmt.Sprintf("hdb%s", p)] = &symbols.Symbol{
			Note: &symbols.Note{
				Embellishment: &embellishment.Embellishment{
					Type:    embellishment.Type_Doubling,
					Variant: embellishment.Variant_Half,
				},
			},
		}
	}
}
