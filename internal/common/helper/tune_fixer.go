package helper

import (
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/musicmodel"
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/tune"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"regexp"
	"strings"
)

type TuneFixer struct {
}

func (tf *TuneFixer) Fix(muMo musicmodel.MusicModel) {
	for _, t := range muMo {
		fixComposerArranger(t)
		fixComposerTrad(t)
		removeTimeSigFromTuneType(t)
		removeSpecialCharsFromTuneType(t)
		trimSpaces(t)
		fixTitle(t)
		capitalizeTuneType(t)
	}
}

func fixComposerArranger(t *tune.Tune) {
	if t.Composer == "" {
		return
	}

	//parts := strings.Split(t.Composer, ",-([")
	regX := regexp.MustCompile(`(?i)arranged|arrangement|arr[.:/ ]+`)
	parts := regX.Split(t.Composer, -1)
	if len(parts) == 2 {
		arranger := fixComposerArrangerField(parts[1])
		sep := regexp.MustCompile(`,|-`)
		arrangerSplit := sep.Split(arranger, -1)

		// tuneTitle after arranger
		if len(arrangerSplit) > 1 {
			t.Arranger = fixComposerArrangerField(arrangerSplit[0])
			t.Composer = fixComposerArrangerField(arrangerSplit[1])
			return
		}

		// tuneTitle before arranger
		if strings.TrimSpace(parts[0]) != "" {
			t.Composer = fixComposerArrangerField(parts[0])
			t.Arranger = arranger
		} else {
			// only arranger in tuneTitle field
			t.Composer = ""
			t.Arranger = arranger
		}
	}
}

func fixComposerTrad(t *tune.Tune) {
	regX := regexp.MustCompile(`(?i)^trad\.?$`)
	trimmedComposer := strings.TrimSpace(t.Composer)
	if regX.MatchString(trimmedComposer) {
		t.Composer = "Traditional"
	}
}

func removeTimeSigFromTuneType(t *tune.Tune) {
	trimmedType := strings.TrimSpace(t.Type)
	regX := regexp.MustCompile(`\d+/\d+`)
	typeWithoutTimesig := regX.ReplaceAllString(trimmedType, "")
	typeWithoutTimesig = strings.TrimSpace(typeWithoutTimesig)
	t.Type = typeWithoutTimesig
}

func removeSpecialCharsFromTuneType(t *tune.Tune) {
	trimmedType := strings.TrimSpace(t.Type)
	trimmedType = strings.Trim(trimmedType, ".:/-|")
	t.Type = trimmedType
}

func trimSpaces(t *tune.Tune) {
	t.Title = strings.TrimSpace(t.Title)
	t.Composer = strings.TrimSpace(t.Composer)
	t.Arranger = strings.TrimSpace(t.Arranger)
	t.Type = strings.TrimSpace(t.Type)
}

func capitalizeTuneType(t *tune.Tune) {
	trimmedType := strings.TrimSpace(t.Type)
	caser := cases.Title(language.English)
	t.Type = caser.String(trimmedType)
}

func fixComposerArrangerField(arr string) string {
	arranger := strings.TrimSpace(arr)
	arranger = strings.Trim(arranger, ".:/-[](),")
	arranger = strings.Replace(arranger, "by", "", -1)
	arranger = strings.TrimSpace(arranger)

	return arranger
}

func fixTitle(t *tune.Tune) {
	t.Title = strings.Replace(t.Title, "_", " ", -1)

	c := cases.Title(language.English)
	t.Title = c.String(t.Title)
}

func NewTuneFixer() *TuneFixer {
	return &TuneFixer{}
}
