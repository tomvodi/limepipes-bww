package symbolmapper

import (
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols/embellishment"
)

func init() {
	tar := &symbols.Symbol{
		Note: &symbols.Note{
			Embellishment: &embellishment.Embellishment{
				Type: embellishment.Type_Taorluath,
			},
		},
	}
	symbolsMap["tar"] = tar
	symbolsMap["tarb"] = tar
	symbolsMap["htar"] = tar
}
