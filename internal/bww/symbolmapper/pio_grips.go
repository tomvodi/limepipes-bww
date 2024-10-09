package symbolmapper

import "fmt"

func init() {
	grps := []string{
		"enbain",
		"otro",
		"odro",
		"adeda",
	}

	for _, s := range grps {
		grps = append(grps, fmt.Sprintf("p%s", s))
		grps = append(grps, fmt.Sprintf("g%s", s))
		grps = append(grps, fmt.Sprintf("t%s", s))
	}

	grps = append(grps, "pgrp", "deda")

	piobSymbols = append(piobSymbols, grps...)
}
