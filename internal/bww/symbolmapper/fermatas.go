package symbolmapper

import "github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols"

var fermatas = []string{
	"fermatlg",
	"fermatla",
	"fermatb",
	"fermatc",
	"fermatd",
	"fermate",
	"fermatf",
	"fermathg",
	"fermatha",
}

func newFermataMap() map[string]*symbols.Note {
	m := make(map[string]*symbols.Note, len(fermatas))
	for _, f := range fermatas {
		m[f] = &symbols.Note{
			Fermata: true,
		}
	}

	return m
}

func init() {
	for k, e := range newFermataMap() {
		symbolsMap[k] = &symbols.Symbol{
			Note: e,
		}
	}
}
