package symbolmapper

func init() {
	dd := []string{
		"darodo",
		"darodo16",
		"hdarodo",
	}

	for _, s := range dd {
		dd = append(dd, "p"+s)
	}

	piobSymbols = append(piobSymbols, dd...)
}
