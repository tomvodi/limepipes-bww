package symbolmapper

func init() {
	for _, p := range lowPitchesLgToHA {
		piobSymbols = append(piobSymbols, "echo"+p)
	}
}
