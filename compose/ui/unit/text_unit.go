package unit

import (
	textunit "github.com/zodimo/go-compose/compose/ui/text/unit"
)

// TextUnit packs unit type (Sp/Em) and float value into 64 bits.
// Struct wrapper prevents `TextUnit(24)` - must use `Sp(24)` or `Em(1.5)`.
type TextUnit = textunit.TextUnit

// TextUnitCompanion re-exports the text-unit companion; the unspecified value is
// reachable only through TextUnitCompanion.Unspecified().
var TextUnitCompanion = textunit.TextUnitCompanion

var Sp = textunit.Sp
var Em = textunit.Em
