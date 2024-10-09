package symbolmapper

import (
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/boundary"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols/timeline"
)

func init() {
	symbolsMap["'1"] = newTimeLine(boundary.Boundary_Start, timeline.Type_First)
	symbolsMap["'si"] = newTimeLine(boundary.Boundary_Start, timeline.Type_Singling)
	symbolsMap["'2"] = newTimeLine(boundary.Boundary_Start, timeline.Type_Second)
	symbolsMap["'do"] = newTimeLine(boundary.Boundary_Start, timeline.Type_Doubling)
	symbolsMap["'bis"] = newTimeLine(boundary.Boundary_Start, timeline.Type_Bis)
	symbolsMap["'22"] = newTimeLine(boundary.Boundary_Start, timeline.Type_SecondOf2)
	symbolsMap["'23"] = newTimeLine(boundary.Boundary_Start, timeline.Type_SecondOf3)
	symbolsMap["'24"] = newTimeLine(boundary.Boundary_Start, timeline.Type_SecondOf4)
	symbolsMap["'224"] = newTimeLine(boundary.Boundary_Start, timeline.Type_SecondOf2And4)
	symbolsMap["'intro"] = newTimeLine(boundary.Boundary_Start, timeline.Type_Intro)
	symbolsMap["'25"] = newTimeLine(boundary.Boundary_Start, timeline.Type_SecondOf5)
	symbolsMap["'26"] = newTimeLine(boundary.Boundary_Start, timeline.Type_SecondOf6)
	symbolsMap["'27"] = newTimeLine(boundary.Boundary_Start, timeline.Type_SecondOf7)
	symbolsMap["'28"] = newTimeLine(boundary.Boundary_Start, timeline.Type_SecondOf8)
	symbolsMap["_'"] = newTimeLine(boundary.Boundary_End, timeline.Type_NoType)
	symbolsMap["bis_'"] = newTimeLine(boundary.Boundary_End, timeline.Type_Bis)
}

func newTimeLine(
	boundType boundary.Boundary,
	timelineType timeline.Type,
) *symbols.Symbol {
	sym := &symbols.Symbol{
		Timeline: &timeline.TimeLine{
			BoundaryType: boundType,
			Type:         timelineType,
		},
	}
	return sym
}
