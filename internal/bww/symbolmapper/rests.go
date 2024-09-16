package symbolmapper

import (
	"fmt"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/length"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols"
)

func init() {
	rm := map[uint8]length.Length{
		1:  length.Length_Whole,
		2:  length.Length_Half,
		4:  length.Length_Quarter,
		8:  length.Length_Eighth,
		16: length.Length_Sixteenth,
		32: length.Length_Thirtysecond,
	}

	for k, v := range rm {
		symbolsMap[fmt.Sprintf("REST_%d", k)] = &symbols.Symbol{
			Rest: &symbols.Rest{
				Length: v,
			},
		}
	}
}
