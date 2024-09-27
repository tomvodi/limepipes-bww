package common

import "fmt"

var ErrSymbolNotFound = fmt.Errorf("symbol not found")
var ErrSymbolSkip = fmt.Errorf("symbol should be skipped")
var ErrLineSkip = fmt.Errorf("line should be skipped")
