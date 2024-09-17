package symbolmapper

import (
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/boundary"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols/tuplet"
)

func init() {
	symbolsMap["^2s"] = newTuplet(boundary.Boundary_Start, 2, 3)
	symbolsMap["^2e"] = newTuplet(boundary.Boundary_End, 2, 3)
	symbolsMap["^3s"] = newTuplet(boundary.Boundary_Start, 3, 2)
	symbolsMap["^3e"] = newTuplet(boundary.Boundary_End, 3, 2)
	symbolsMap["^43s"] = newTuplet(boundary.Boundary_Start, 4, 3)
	symbolsMap["^43e"] = newTuplet(boundary.Boundary_End, 4, 3)
	symbolsMap["^46s"] = newTuplet(boundary.Boundary_Start, 4, 6)
	symbolsMap["^46e"] = newTuplet(boundary.Boundary_End, 4, 6)
	symbolsMap["^46s"] = newTuplet(boundary.Boundary_Start, 4, 6)
	symbolsMap["^46e"] = newTuplet(boundary.Boundary_End, 4, 6)
	symbolsMap["^53s"] = newTuplet(boundary.Boundary_Start, 5, 3)
	symbolsMap["^53e"] = newTuplet(boundary.Boundary_End, 5, 3)
	symbolsMap["^54s"] = newTuplet(boundary.Boundary_Start, 5, 4)
	symbolsMap["^54e"] = newTuplet(boundary.Boundary_End, 5, 4)
	symbolsMap["^64s"] = newTuplet(boundary.Boundary_Start, 6, 4)
	symbolsMap["^64e"] = newTuplet(boundary.Boundary_End, 6, 4)
	symbolsMap["^74s"] = newTuplet(boundary.Boundary_Start, 7, 4)
	symbolsMap["^74e"] = newTuplet(boundary.Boundary_End, 7, 4)
	symbolsMap["^76s"] = newTuplet(boundary.Boundary_Start, 7, 6)
	symbolsMap["^76e"] = newTuplet(boundary.Boundary_End, 7, 6)
}

func newTuplet(
	boundaryType boundary.Boundary,
	visibleNotes uint32,
	playedNotes uint32,
) *symbols.Symbol {
	sym := &symbols.Symbol{
		Tuplet: &tuplet.Tuplet{
			BoundaryType: boundaryType,
			VisibleNotes: visibleNotes,
			PlayedNotes:  playedNotes,
		},
	}

	return sym
}
