package symbolmerger

import "github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols"

type accidentalMerger struct {
}

func (x *accidentalMerger) MergeSymbols(
	left *symbols.Symbol,
	right *symbols.Symbol,
) bool {
	if left.IsOnlyAccidental() && right.IsValidNote() {
		right.Note.Accidental = left.Note.Accidental
		left.Note = right.Note

		return true
	}

	return false
}

func init() {
	mergers = append(mergers, &accidentalMerger{})
}
