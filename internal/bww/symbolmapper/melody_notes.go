package symbolmapper

import (
	"fmt"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/length"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/pitch"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols"
	"maps"
)

var pitches = []string{"LG", "LA", "B", "C", "D", "E", "F", "HG", "HA"}
var lengthesAll = []uint8{1, 2, 4, 8, 16, 32}
var lengthesFlag = []uint8{8, 16, 32}
var flags = []string{"l", "r"}
var pitchMap = map[string]pitch.Pitch{
	"LG": pitch.Pitch_LowG,
	"LA": pitch.Pitch_LowA,
	"B":  pitch.Pitch_B,
	"C":  pitch.Pitch_C,
	"D":  pitch.Pitch_D,
	"E":  pitch.Pitch_E,
	"F":  pitch.Pitch_F,
	"HG": pitch.Pitch_HighG,
	"HA": pitch.Pitch_HighA,
}

var lengthMap = map[uint8]length.Length{
	1:  length.Length_Whole,
	2:  length.Length_Half,
	4:  length.Length_Quarter,
	8:  length.Length_Eighth,
	16: length.Length_Sixteenth,
	32: length.Length_Thirtysecond,
}

func newMelodyNotesMap() map[string]*symbols.Symbol {
	m := make(map[string]*symbols.Symbol, 108)
	const noFlagType = "%s_%d"
	const flagType = "%s%s_%d"
	for _, p := range pitches {
		for _, l := range lengthesAll {

			m[fmt.Sprintf(noFlagType, p, l)] = &symbols.Symbol{
				Note: &symbols.Note{
					Pitch:  pitchMap[p],
					Length: lengthMap[l],
				},
			}
		}
		for _, l := range lengthesFlag {
			for _, f := range flags {
				m[fmt.Sprintf(flagType, p, f, l)] = &symbols.Symbol{
					Note: &symbols.Note{
						Pitch:  pitchMap[p],
						Length: lengthMap[l],
					},
				}
			}
		}
	}

	return m
}

func init() {
	maps.Copy(symbolsMap, newMelodyNotesMap())
}
