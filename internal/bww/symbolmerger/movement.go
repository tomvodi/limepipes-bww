package symbolmerger

import "github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols"

type movementMerger struct {
}

func (x *movementMerger) MergeSymbols(
	left *symbols.Symbol,
	right *symbols.Symbol,
) bool {
	if left.IsOnlyMovement() && right.IsValidNote() {
		right.Note.Movement = left.Note.Movement
		left.Note = right.Note

		return true
	}

	return false
}

func init() {
	mergers = append(mergers, &movementMerger{})
}
