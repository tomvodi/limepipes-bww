package symbolmapper

import "fmt"

func init() {
	cads := []string{
		"cadged",
		"cadge",
		"caded",
		"cade",
		"cadaed",
		"cadae",
		"cadgf",
		"cadaf",
	}

	for _, s := range cads {
		cads = append(cads, fmt.Sprintf("f%s", s))
	}

	piobSymbols = append(piobSymbols, cads...)
}
