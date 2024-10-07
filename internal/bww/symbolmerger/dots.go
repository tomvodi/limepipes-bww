package symbolmerger

import "github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols"

type dotsMerger struct {
}

func (x *dotsMerger) MergeSymbols(
	left *symbols.Symbol,
	right *symbols.Symbol,
) bool {
	if left.IsValidNote() && right.IsOnlyDots() {
		left.Note.Dots = right.Note.Dots

		return true
	}

	return false
}

func init() {
	mergers = append(mergers, &dotsMerger{})
}
