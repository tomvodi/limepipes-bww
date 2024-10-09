package symbolmapper

import "github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols"

var oneDot = []string{"'lg", "'la", "'b", "'c", "'d", "'e", "'f", "'hg", "'ha"}
var twoDots = []string{"''lg", "''la", "''b", "''c", "''d", "''e", "''f", "''hg", "''ha"}

func newDotsMap() map[string]*symbols.Note {
	m := make(map[string]*symbols.Note, len(oneDot)+len(twoDots))
	for _, p := range oneDot {
		m[p] = &symbols.Note{
			Dots: 1,
		}
	}
	for _, p := range twoDots {
		m[p] = &symbols.Note{
			Dots: 2,
		}
	}

	return m
}

func init() {
	for k, e := range newDotsMap() {
		symbolsMap[k] = &symbols.Symbol{
			Note: e,
		}
	}
}
