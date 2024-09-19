package symbolmapper

import (
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/pitch"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/symbols/movement"
)

func init() {
	symbolsMap["embari"] = newMovement(movement.Type_Embari)
	symbolsMap["endari"] = newMovement(movement.Type_Endari)
	symbolsMap["chedari"] = newMovement(movement.Type_Chedari)
	symbolsMap["hedari"] = newMovement(movement.Type_Hedari)
	symbolsMap["pembari"] = newMovement(movement.Type_Embari, true)
	symbolsMap["pendari"] = newMovement(movement.Type_Endari, true)
	symbolsMap["pchedari"] = newMovement(movement.Type_Chedari, true)
	symbolsMap["phedari"] = newMovement(movement.Type_Hedari, true)

	symbolsMap["dili"] = newMovement(movement.Type_Dili)
	symbolsMap["tra"] = newMovement(movement.Type_Tra)
	symbolsMap["htra"] = newMovement(movement.Type_Tra, movement.Variant_Half)
	symbolsMap["tra8"] = newMovement(movement.Type_Tra, movement.Variant_LongLowG)
	symbolsMap["pdili"] = newMovement(movement.Type_Dili, true)
	symbolsMap["ptra"] = newMovement(movement.Type_Tra, true)
	symbolsMap["phtra"] = newMovement(movement.Type_Tra, true, movement.Variant_Half)
	symbolsMap["ptra8"] = newMovement(movement.Type_Tra, true, movement.Variant_LongLowG)

	symbolsMap["edre"] = newMovement(movement.Type_Edre)
	symbolsMap["edreb"] = newMovement(movement.Type_Edre, pitch.Pitch_B)
	symbolsMap["edrec"] = newMovement(movement.Type_Edre, pitch.Pitch_C)
	symbolsMap["edred"] = newMovement(movement.Type_Edre, pitch.Pitch_D)
	symbolsMap["pedre"] = newMovement(movement.Type_Edre, true)
	symbolsMap["pedreb"] = newMovement(movement.Type_Edre, true, pitch.Pitch_B)
	symbolsMap["pedrec"] = newMovement(movement.Type_Edre, true, pitch.Pitch_C)
	symbolsMap["pedred"] = newMovement(movement.Type_Edre, true, pitch.Pitch_D)

	symbolsMap["dare"] = newMovement(movement.Type_Dare)
	symbolsMap["chedare"] = newMovement(movement.Type_Dare)
	symbolsMap["pdare"] = newMovement(movement.Type_Dare, true)
	symbolsMap["chechere"] = newMovement(movement.Type_CheCheRe)
	symbolsMap["pchechere"] = newMovement(movement.Type_CheCheRe, true)

	symbolsMap["gedre"] = newMovement(movement.Type_Edre, movement.Variant_G)
	symbolsMap["gdare"] = newMovement(movement.Type_Dare, movement.Variant_G)
	symbolsMap["tedre"] = newMovement(movement.Type_Edre, movement.Variant_Thumb)
	symbolsMap["tdare"] = newMovement(movement.Type_Dare, movement.Variant_Thumb)
	symbolsMap["tchechere"] = newMovement(movement.Type_CheCheRe, movement.Variant_Thumb)

	symbolsMap["dre"] = newMovement(movement.Type_Edre, movement.Variant_Half)
	symbolsMap["hedale"] = newMovement(movement.Type_Dare, movement.Variant_Half)
	symbolsMap["hchechere"] = newMovement(movement.Type_CheCheRe, movement.Variant_Half)
}
