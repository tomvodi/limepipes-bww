package symbolmapper

import "github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/barline"

func init() {
	barlineMap["!"] = nil // a regular barline can be nil
	barlineMap["I!"] = &barline.Barline{
		Type: barline.Type_Heavy,
	}
	barlineMap["!I"] = &barline.Barline{
		Type: barline.Type_Heavy,
	}
	barlineMap["I!''"] = &barline.Barline{
		Type: barline.Type_Heavy,
		Time: barline.Time_Repeat,
	}
	barlineMap["''!I"] = &barline.Barline{
		Type: barline.Type_Heavy,
		Time: barline.Time_Repeat,
	}
}
