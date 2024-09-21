package parserstate

//go:generate go run github.com/dmarkham/enumer -json -yaml -transform=lower -type=ParserState

type ParserState uint8

const (
	File ParserState = iota // Top level state
	Stave
)
