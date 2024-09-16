package symbolmapper

import (
	"github.com/tomvodi/limepipes-plugin-api/musicmodel/v1/measure"
	"golang.org/x/exp/maps"
)

func newTimeSignatureMap() map[string]*measure.TimeSignature {
	return map[string]*measure.TimeSignature{
		"2_2": {
			Beats:    2,
			BeatType: 2,
		},
		"3_2": {
			Beats:    3,
			BeatType: 2,
		},
		"2_4": {
			Beats:    2,
			BeatType: 4,
		},
		"3_4": {
			Beats:    3,
			BeatType: 4,
		},
		"4_4": {
			Beats:    4,
			BeatType: 4,
		},
		"5_4": {
			Beats:    5,
			BeatType: 4,
		},
		"6_4": {
			Beats:    6,
			BeatType: 4,
		},
		"7_4": {
			Beats:    7,
			BeatType: 4,
		},
		"C_": {
			Beats:    2,
			BeatType: 2,
		},
		"C": {
			Beats:    4,
			BeatType: 4,
		},
		"2_8": {
			Beats:    2,
			BeatType: 8,
		},
		"3_8": {
			Beats:    3,
			BeatType: 8,
		},
		"4_8": {
			Beats:    4,
			BeatType: 8,
		},
		"5_8": {
			Beats:    5,
			BeatType: 8,
		},
		"6_8": {
			Beats:    6,
			BeatType: 8,
		},
		"7_8": {
			Beats:    7,
			BeatType: 8,
		},
		"8_8": {
			Beats:    8,
			BeatType: 8,
		},
		"9_8": {
			Beats:    9,
			BeatType: 8,
		},
		"10_8": {
			Beats:    10,
			BeatType: 8,
		},
		"11_8": {
			Beats:    11,
			BeatType: 8,
		},
		"12_8": {
			Beats:    12,
			BeatType: 8,
		},
		"15_8": {
			Beats:    15,
			BeatType: 8,
		},
		"18_8": {
			Beats:    18,
			BeatType: 8,
		},
		"21_8": {
			Beats:    21,
			BeatType: 8,
		},
		"2_16": {
			Beats:    2,
			BeatType: 16,
		},
		"3_16": {
			Beats:    3,
			BeatType: 16,
		},
		"4_16": {
			Beats:    4,
			BeatType: 16,
		},
		"5_16": {
			Beats:    5,
			BeatType: 16,
		},
		"6_16": {
			Beats:    6,
			BeatType: 16,
		},
		"7_16": {
			Beats:    7,
			BeatType: 16,
		},
		"8_16": {
			Beats:    8,
			BeatType: 16,
		},
		"9_16": {
			Beats:    9,
			BeatType: 16,
		},
		"10_16": {
			Beats:    10,
			BeatType: 16,
		},
		"11_16": {
			Beats:    11,
			BeatType: 16,
		},
		"12_16": {
			Beats:    12,
			BeatType: 16,
		},
	}
}

func init() {
	maps.Copy(timeSignatureMap, newTimeSignatureMap())
	timeSignatureKeys = make([]string, len(timeSignatureMap))
	for i, k := range maps.Keys(timeSignatureMap) {
		timeSignatureKeys[i] = k
	}
}
