package symbolmapper

func init() {
	for _, p := range lowPitchesLgToC {
		piobSymbols = append(piobSymbols, "ptrip"+p)
		piobSymbols = append(piobSymbols, "pttrip"+p)
		piobSymbols = append(piobSymbols, "phtrip"+p)
	}
}
