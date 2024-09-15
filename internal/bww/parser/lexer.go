package parser

import "github.com/alecthomas/participle/v2/lexer"

var Lexer = lexer.MustStateful(lexer.Rules{
	"Root": {
		{
			Name:    "BagpipeReader",
			Pattern: `Bagpipe Reader|Bagpipe Music Writer Gold|Bagpipe Musicworks Gold`,
			Action:  lexer.Push("BagpipeReader"),
		},
		{
			Name:    "TEMPO_DEF",
			Pattern: `TuneTempo`,
			Action:  lexer.Push("TuneTempo"),
		},
		{
			Name:    "PARAM_START",
			Pattern: `\(`,
			Action:  lexer.Push("ParamList"),
		},
		{
			Name:    "TIMELINE_END",
			Pattern: `_'|bis_'`,
		},
		{
			Name:    "DALSEGNO",
			Pattern: `dalsegno`,
		},
		{
			Name:    "DACAPOALFINE",
			Pattern: `dacapoalfine`,
		},
		{
			Name:    "PARAM_DEF",
			Pattern: `MIDINoteMappings|FrequencyMappings|InstrumentMappings|GracenoteDurations|FontSizes|TuneFormat`,
		},
		{
			Name:    "PARAM_SEP",
			Pattern: `,`,
		},
		{
			Name:    "STAFF_START",
			Pattern: `&`,
			Action:  lexer.Push("Staff"),
		},
		{
			Name:    "STRING",
			Pattern: `"[^"]*"`,
		},
		{
			Name:    "WHITESPACE",
			Pattern: `\s+`,
		},
	},
	"BagpipeReader": {
		{
			Name:    "VERSION_SEP",
			Pattern: `:`,
		},
		{
			Name:    "VersionNumber",
			Pattern: `\d+\.\d+`,
			Action:  lexer.Pop(),
		},
	},
	"TuneTempo": {
		{
			Name:    "PARAM_SEP",
			Pattern: `,`,
		},
		{
			Name:    "TEMPO_VALUE",
			Pattern: `\d+`,
			Action:  lexer.Pop(),
		},
	},
	"ParamList": {
		{
			Name:    "PARAM_END",
			Pattern: `\)`,
			Action:  lexer.Pop(),
		},
		{
			Name:    "PARAM",
			Pattern: `[^,\)]+`,
		},
		{
			Name:    "PARAM_SEP",
			Pattern: `,`,
		},
	},
	"Staff": {
		{
			Name:    "STAFF_END",
			Pattern: `''!I|!t|!I`,
			Action:  lexer.Pop(),
		},
		{
			Name:    "PART_START",
			Pattern: `I!''|I!`,
		},
		{
			Name:    "NEXT_STAFF_START",
			Pattern: `&`,
		},
		{
			Name:    "BARLINE",
			Pattern: `!`,
		},
		{
			Name:    "SPACE",
			Pattern: `space`,
		},
		{
			Name:    "STRING",
			Pattern: `"[^"]*"`,
		},
		{
			Name:    "TEMPO_DEF",
			Pattern: `TuneTempo`,
			Action:  lexer.Push("TuneTempo"),
		},
		{
			Name:    "PARAM_START",
			Pattern: `\(`,
			Action:  lexer.Push("ParamList"),
		},
		{
			Name:    "PARAM_SEP",
			Pattern: `,`,
		},
		{
			Name:    "MUSIC_SYMBOL",
			Pattern: `[A-Za-z0-9_'^]+`,
		},
		{
			Name:    "WHITESPACE",
			Pattern: `\s+`,
		},
	},
})
