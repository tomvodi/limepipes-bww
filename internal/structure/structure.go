package structure

type TimelineEnd string
type Comment string
type InlineText string
type Tempo int
type TuneTitle string
type TuneType string
type TuneComposer string
type TuneFooter string
type TuneInline string
type Barline string

type BwwFile struct {
	BagpipePlayerVersion string
	TuneDefs             []TuneDefinition
}

type TuneDefinition struct {
	Data []byte
	Tune Tune
}

type Tune struct {
	Header TuneHeader
	Staffs []Staff
}

type TuneHeader struct {
	Title      string
	Type       string
	Composer   string
	Footer     string
	InlineText string
}

type Staff struct {
	Measures []Measure
}

type Measure struct {
	Components []any
}

type MusicSymbol struct {
	Pos  Position
	Text string
}

type Position struct {
	Line   int
	Column int
}
