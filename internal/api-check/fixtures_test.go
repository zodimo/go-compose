package apicheck

import (
	"strings"
	"testing"
)



// fixture is a known gioui leak site captured before the cleanup.
type fixture struct {
	pkg    string // package import path
	symbol string // exported symbol name (Type, Func, Method as Type.Method, const, var)
}

// knownLeaks is the regression-fixture inventory: every pre-cleanup gioui
// leak site in the public package trees (spec: api-purity-analyzer — "Known
// leaks are regression fixtures"). Derived from the exploration inventory of
// the "abstract-render-backend" change.
//
// The fixtures assert RED state: each site must be flagged by the analyzer.
// After Phase 2 cleanup flips them to negative assertions (none flagged).
var knownLeaks = []fixture{
	// compose/ui/unit — gio_helpers.go + TextUnit.AsGioSp.
	{pkg: "github.com/zodimo/go-compose/compose/ui/unit", symbol: "DpToGioUnit"},
	{pkg: "github.com/zodimo/go-compose/compose/ui/unit", symbol: "DpToGioUnitUnsafe"},
	{pkg: "github.com/zodimo/go-compose/compose/ui/unit", symbol: "TextUnitToGioSp"},
	{pkg: "github.com/zodimo/go-compose/compose/ui/unit", symbol: "TextUnitToGioDp"},
	{pkg: "github.com/zodimo/go-compose/compose/ui/unit", symbol: "TextUnitToGioSpUnsafe"},
	{pkg: "github.com/zodimo/go-compose/compose/ui/unit", symbol: "DensityFromLayoutContext"},
	{pkg: "github.com/zodimo/go-compose/compose/ui/unit", symbol: "TextUnit.AsGioSp"},

	// compose/ui/text/style — enum underlying types + conversions.
	{pkg: "github.com/zodimo/go-compose/compose/ui/text/style", symbol: "TextAlign"},
	{pkg: "github.com/zodimo/go-compose/compose/ui/text/style", symbol: "LineBreak"},
	{pkg: "github.com/zodimo/go-compose/compose/ui/text/style", symbol: "FromGioTextAlign"},
	{pkg: "github.com/zodimo/go-compose/compose/ui/text/style", symbol: "TextAlignToGioTextAlignment"},
	{pkg: "github.com/zodimo/go-compose/compose/ui/text/style", symbol: "LineBreakToGioWrapPolicy"},
	{pkg: "github.com/zodimo/go-compose/compose/ui/text/style", symbol: "GioWrapPolicyToLineBreak"},

	// compose/ui/text/font — font converters.
	{pkg: "github.com/zodimo/go-compose/compose/ui/text/font", symbol: "ToGioFont"},
	{pkg: "github.com/zodimo/go-compose/compose/ui/text/font", symbol: "FromGioFont"},
	{pkg: "github.com/zodimo/go-compose/compose/ui/text/font", symbol: "ToGioWeight"},
	{pkg: "github.com/zodimo/go-compose/compose/ui/text/font", symbol: "FromGioWeight"},
	{pkg: "github.com/zodimo/go-compose/compose/ui/text/font", symbol: "ToGioStyle"},
	{pkg: "github.com/zodimo/go-compose/compose/ui/text/font", symbol: "FromGioStyle"},
	{pkg: "github.com/zodimo/go-compose/compose/ui/text/font", symbol: "FromGioTypeface"},

	// compose/ui/text — text helpers + TextShaper.
	{pkg: "github.com/zodimo/go-compose/compose/ui/text", symbol: "TextStyleFromGioFont"},
	{pkg: "github.com/zodimo/go-compose/compose/ui/text", symbol: "ToGioFont"},
	{pkg: "github.com/zodimo/go-compose/compose/ui/text", symbol: "TextShaper"},

	// compose/ui/next/text + font — mirror implementations.
	{pkg: "github.com/zodimo/go-compose/compose/ui/next/text", symbol: "TextShaper"},
	{pkg: "github.com/zodimo/go-compose/compose/ui/next/text/font", symbol: "ToGioFont"},
	{pkg: "github.com/zodimo/go-compose/compose/ui/next/text/font", symbol: "ToGioWeight"},
	{pkg: "github.com/zodimo/go-compose/compose/ui/next/text/font", symbol: "ToGioStyle"},

	// compose/ui/text/intl + next/text/intl — Locale alias.
	{pkg: "github.com/zodimo/go-compose/compose/ui/text/intl", symbol: "Locale"},
	{pkg: "github.com/zodimo/go-compose/compose/ui/next/text/intl", symbol: "Locale"},

	// compose/ui/platform — LocalWindow CompositionLocal of *app.Window.
	{pkg: "github.com/zodimo/go-compose/compose/ui/platform", symbol: "LocalWindow"},

	// compose/ui/graphics — ImageResource field + shape outline.
	{pkg: "github.com/zodimo/go-compose/compose/ui/graphics", symbol: "ImageResource"},
	{pkg: "github.com/zodimo/go-compose/compose/ui/graphics/shape", symbol: "Shape.CreateOutline"},
	{pkg: "github.com/zodimo/go-compose/compose/ui/graphics/shape", symbol: "Outline.Push"},
	{pkg: "github.com/zodimo/go-compose/compose/ui/graphics/shape", symbol: "Outline.Op"},
	{pkg: "github.com/zodimo/go-compose/compose/ui/graphics/shape", symbol: "Outline.Path"},

	// theme — BasicTheme alias + manager.
	{pkg: "github.com/zodimo/go-compose/theme", symbol: "BasicTheme"},

	// compose/material — ThemeInterface + LocalGioMaterialTheme.
	{pkg: "github.com/zodimo/go-compose/compose/material", symbol: "ThemeInterface.GioMaterialTheme"},
	{pkg: "github.com/zodimo/go-compose/compose/material", symbol: "LocalGioMaterialTheme"},

	// compose/foundation/layout — box/column/row/overlay aliases.
	{pkg: "github.com/zodimo/go-compose/runtime", symbol: "Runtime.Run"},
	{pkg: "github.com/zodimo/go-compose/compose/foundation/layout/box", symbol: "Stack"},
	{pkg: "github.com/zodimo/go-compose/compose/foundation/layout/box", symbol: "StackChild"},
	{pkg: "github.com/zodimo/go-compose/compose/foundation/layout/box", symbol: "LayoutContext"},
	{pkg: "github.com/zodimo/go-compose/compose/foundation/layout/box", symbol: "LayoutDimensions"},
	{pkg: "github.com/zodimo/go-compose/compose/foundation/layout/column", symbol: "Spacing"},
	{pkg: "github.com/zodimo/go-compose/compose/foundation/layout/column", symbol: "Alignment"},
	{pkg: "github.com/zodimo/go-compose/compose/foundation/layout/row", symbol: "Spacing"},
	{pkg: "github.com/zodimo/go-compose/compose/foundation/layout/row", symbol: "Alignment"},
	{pkg: "github.com/zodimo/go-compose/compose/foundation/layout/overlay", symbol: "LayoutContext"},
	{pkg: "github.com/zodimo/go-compose/compose/foundation/layout/overlay", symbol: "LayoutDimensions"},

	// compose/foundation/text + next/text — enum aliases.
	{pkg: "github.com/zodimo/go-compose/compose/foundation/text", symbol: "Alignment"},
	{pkg: "github.com/zodimo/go-compose/compose/foundation/text", symbol: "WrapPolicy"},
	{pkg: "github.com/zodimo/go-compose/compose/foundation/next/text", symbol: "Alignment"},
	{pkg: "github.com/zodimo/go-compose/compose/foundation/next/text", symbol: "WrapPolicy"},

	// compose/foundation/lazy — state structs + C/D aliases.
	{pkg: "github.com/zodimo/go-compose/compose/foundation/lazy", symbol: "LazyListState"},
	{pkg: "github.com/zodimo/go-compose/compose/foundation/lazy", symbol: "LazyGridState"},
	{pkg: "github.com/zodimo/go-compose/compose/foundation/lazy", symbol: "C"},
	{pkg: "github.com/zodimo/go-compose/compose/foundation/lazy", symbol: "D"},

	// compose/material3/next/textfield — models + C/D aliases.
	{pkg: "github.com/zodimo/go-compose/compose/material3/next/textfield", symbol: "TextFieldWidget"},
	{pkg: "github.com/zodimo/go-compose/compose/material3/next/textfield", symbol: "TextFieldComponent"},
	{pkg: "github.com/zodimo/go-compose/compose/material3/next/textfield", symbol: "C"},
	{pkg: "github.com/zodimo/go-compose/compose/material3/next/textfield", symbol: "D"},

	// compose/material3/textfield — widget structs embedding editor/click.
	{pkg: "github.com/zodimo/go-compose/compose/material3/textfield", symbol: "TextFieldWidget"},
	{pkg: "github.com/zodimo/go-compose/compose/material3/textfield", symbol: "FilledTextFieldWidget"},
	{pkg: "github.com/zodimo/go-compose/compose/material3/textfield", symbol: "OutlinedTextFieldWidget"},

	// modifiers — padding.RTL const, clickable alias.
	{pkg: "github.com/zodimo/go-compose/modifiers/padding", symbol: "RTL"},
	{pkg: "github.com/zodimo/go-compose/modifiers/clickable", symbol: "GioClickable"},

	// compose/material3/card — GioImage alias.
	{pkg: "github.com/zodimo/go-compose/compose/material3/card", symbol: "GioImage"},

	// compose/material3/snackbar — SnackbarData clickable fields.
	{pkg: "github.com/zodimo/go-compose/compose/material3/snackbar", symbol: "SnackbarData"},

	// runtime — Runtime.Run returns op.CallOp.
	{pkg: "github.com/zodimo/go-compose/runtime", symbol: "Runtime.Run"},

	// pkg/x/fileexplorer — RememberExplorer callback param.
	{pkg: "github.com/zodimo/go-compose/pkg/x/fileexplorer", symbol: "RememberExplorer"},
}

// TestKnownLeaksAreFlagged asserts the analyzer is in RED state: every known
// leak site captured as a fixture is flagged. This is the "fixtures first"
// step of the change — it passes while the leaks exist and must be flipped to
// a negative assertion (TestKnownLeaksAreClean) once Phase 2 removes them.
func TestKnownLeaksAreFlagged(t *testing.T) {
	violations, err := Check(Options{})
	if err != nil {
		t.Fatalf("Check: %v", err)
	}

	flagged := map[fixture]bool{}
	for _, v := range violations {
		if v.Kind == KindSignature {
			flagged[fixture{pkg: v.Pkg, symbol: v.Symbol}] = true
		}
	}

	var missing []string
	for _, f := range knownLeaks {
		if !flagged[f] {
			missing = append(missing, f.pkg+" "+f.symbol)
		}
	}
	if len(missing) > 0 {
		t.Errorf("analyzer failed to flag %d known leak fixture(s):\n  %s",
			len(missing), strings.Join(missing, "\n  "))
	}
	t.Logf("analyzer flagged %d signature violation(s) total", len(flagged))
}
