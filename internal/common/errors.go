package common

import "fmt"

var ErrSymbolNotFound = fmt.Errorf("symbol not found")
var ErrSymbolSkip = fmt.Errorf("symbol should be skipped")
