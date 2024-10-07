package symbolmerger

import "github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols"

type Merger struct {
}

func (x *Merger) MergeSymbols(
	left *symbols.Symbol,
	right *symbols.Symbol,
) bool {
	// tie start is merged with the next note
	if left.IsOnlyTieStart() && right.IsValidNote() {
		right.Note.Tie = left.Note.Tie
		left.Note = right.Note

		return true
	}

	// tie end is merged with the previous note
	if left.IsValidNote() && right.IsOnlyTieEnd() {
		left.Note.Tie = right.Note.Tie

		return true
	}

	return false
}

func init() {
	mergers = append(mergers, &Merger{})
}
