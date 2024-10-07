package symbolmerger

import "github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols"

type fermataMerger struct {
}

func (x *fermataMerger) MergeSymbols(
	left *symbols.Symbol,
	right *symbols.Symbol,
) bool {
	if left.IsValidNote() && right.IsOnlyFermata() {
		left.Note.Fermata = right.Note.Fermata

		return true
	}

	return false
}

func init() {
	mergers = append(mergers, &fermataMerger{})
}
