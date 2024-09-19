package symbolmapper

import (
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/pitch"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols/embellishment"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols/movement"
)

var lowPitchesLgToHA = []string{"lg", "la", "b", "c", "d", "e", "f", "hg", "ha"}
var lowPitchesLgToC = []string{"lg", "la", "b", "c"}
var lowPitchesLgToD = []string{"lg", "la", "b", "c", "d"}
var lowPitchesLgToE = []string{"lg", "la", "b", "c", "d", "e"}
var lowPitchesLgToF = []string{"lg", "la", "b", "c", "d", "e", "f"}
var lowPitchesLgToHG = []string{"lg", "la", "b", "c", "d", "e", "f", "hg"}
var lowPitchesLaToHG = []string{"la", "b", "c", "d", "e", "f", "hg"}
var lowPitchesLaToHA = []string{"la", "b", "c", "d", "e", "f", "hg", "ha"}
var lowPitchesLaToF = []string{"la", "b", "c", "d", "e", "f"}

func newEmbellishment(
	eType embellishment.Type,
	args ...interface{},
) *symbols.Symbol {
	sym := &symbols.Symbol{
		Note: &symbols.Note{
			Embellishment: &embellishment.Embellishment{
				Type: eType,
			},
		},
	}

	for _, arg := range args {
		switch t := arg.(type) {
		case embellishment.Variant:
			sym.Note.Embellishment.Variant = t
		case pitch.Pitch:
			sym.Note.Embellishment.Pitch = t
		case embellishment.Weight:
			sym.Note.Embellishment.Weight = t
		default:
			panic("Unknown argument to newEmbellishment")
		}
	}

	return sym
}

func newMovement(
	mType movement.Type,
	args ...interface{},
) *symbols.Symbol {
	sym := &symbols.Symbol{
		Note: &symbols.Note{
			Movement: &movement.Movement{
				Type: mType,
			},
		},
	}

	for _, arg := range args {
		switch t := arg.(type) {
		case movement.Variant:
			sym.Note.Movement.Variant = t
		case bool:
			sym.Note.Movement.Abbreviate = t
		case pitch.Pitch:
			sym.Note.Movement.PitchHint = t
		default:
			panic("Unknown argument to newMovement")
		}
	}

	return sym
}
