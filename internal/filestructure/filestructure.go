package filestructure

type BagpipePlayerVersion string
type TimelineEnd string
type TuneComment string
type InlineComment string
type InlineText string
type TuneTitle string
type TuneType string
type TuneComposer string
type TuneFooter string
type TuneInline string
type Barline string
type StaffStart string
type StaffEnd string
type DalSegno string
type DacapoAlFine string
type StaffComment string
type StaffInline string
type TuneTempo uint32
type TempoChange uint32

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
	Tempo       TuneTempo
}

type Measure struct {
	StaffComments    []StaffComment // A comment that is placed directly above a staff
	StaffInlineTexts []StaffInline  // Text that is placed directly above a staff
	InlineTexts      []InlineText
	InlineComments   []InlineComment
	LeftBarline      Barline
	RightBarline     Barline
	Symbols          []*MusicSymbol
}

type MusicSymbol struct {
	Pos         Position
	Text        string
	InlineTexts []InlineText
	Comments    []InlineComment
	TempoChange TempoChange
}

func (m *MusicSymbol) IsTempoChange() bool {
	return m.TempoChange > 0
}

type Position struct {
	Line   int
	Column int
}
