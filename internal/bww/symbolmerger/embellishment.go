package symbolmerger

import "github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols"

type embellishmentMerger struct {
}

func (x *embellishmentMerger) MergeSymbols(
	left *symbols.Symbol,
	right *symbols.Symbol,
) bool {
	if left.IsOnlyEmbellishment() && right.IsValidNote() {
		right.Note.Embellishment = left.Note.Embellishment
		left.Note = right.Note

		return true
	}

	return false
}

func init() {
	mergers = append(mergers, &embellishmentMerger{})
}
