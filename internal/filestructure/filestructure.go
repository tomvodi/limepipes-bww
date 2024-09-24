package filestructure

type BagpipePlayerVersion string
type TimelineEnd string
type TuneComment string
type InlineComment string
type InlineText string
type Tempo int
type TuneTitle string
type TuneType string
type TuneComposer string
type TuneFooter string
type TuneInline string
type Barline string
type StaffStart string
type StaffEnd string
type StaffComment string
type StaffInline string

type BwwFile struct {
	BagpipePlayerVersion BagpipePlayerVersion
	TuneDefs             []TuneDefinition
}

type TuneDefinition struct {
	Data []byte
	Tune *Tune
}

type Tune struct {
	Header   *TuneHeader
	Measures []*Measure
}

type TuneHeader struct {
	Title       TuneTitle
	Type        TuneType
	Composer    TuneComposer
	Footer      []TuneFooter
	InlineTexts []TuneInline
	Comments    []TuneComment
}

type Measure struct {
	StaffComments    []StaffComment
	StaffInlineTexts []StaffInline
	InlineTexts      []InlineText
	InlineComments   []InlineComment
	Symbols          []*MusicSymbol
}

type MusicSymbol struct {
	Pos         Position
	Text        string
	InlineTexts []InlineText
	Comments    []InlineComment
}

type Position struct {
	Line   int
	Column int
}
