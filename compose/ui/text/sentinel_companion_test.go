package text

import (
	"testing"

	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/text/font"
	"github.com/zodimo/go-compose/compose/ui/text/style"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

// === Options must override, and must not be swallowed ===

// The pre-rewrite CopyTextStyle seeded its working value from the sentinel and then
// merged the sentinel as a base, so an option like WithFontSize(24) disappeared.
func TestCopyTextStyleWithOptions_OptionsOverride(t *testing.T) {
	t.Run("single_option_is_applied", func(t *testing.T) {
		ts := TextStyleFromOptions(WithFontSize(unit.Sp(24)))
		if ts.FontSize() != unit.Sp(24) {
			t.Errorf("FontSize() = %v, want Sp(24) — option was swallowed", ts.FontSize())
		}
	})

	t.Run("one_option_does_not_clobber_siblings", func(t *testing.T) {
		ts := TextStyleFromOptions(
			WithFontSize(unit.Sp(20)),
			WithColor(graphics.ColorRed),
			WithFontWeight(font.FontWeightBold),
			WithTextAlign(style.TextAlignMiddle),
			WithLineHeight(unit.Sp(30)),
		)
		if ts.FontSize() != unit.Sp(20) {
			t.Errorf("FontSize() = %v, want Sp(20)", ts.FontSize())
		}
		if ts.Color() != graphics.ColorRed {
			t.Errorf("Color() = %v, want Red", ts.Color())
		}
		if ts.FontWeight() != font.FontWeightBold {
			t.Errorf("FontWeight() = %v, want Bold", ts.FontWeight())
		}
		if ts.TextAlign() != style.TextAlignMiddle {
			t.Errorf("TextAlign() = %v, want Middle", ts.TextAlign())
		}
		if ts.LineHeight() != unit.Sp(30) {
			t.Errorf("LineHeight() = %v, want Sp(30)", ts.LineHeight())
		}
	})

	t.Run("option_overrides_source_and_inherits_the_rest", func(t *testing.T) {
		base := TextStyleFromOptions(
			WithColor(graphics.ColorGreen),
			WithFontSize(unit.Sp(16)),
			WithFontWeight(font.FontWeightBold),
		)
		cp := CopyTextStyleWithOptions(base,
			WithColor(graphics.ColorBlue),
			WithFontSize(unit.Sp(22)),
		)

		if cp.Color() != graphics.ColorBlue {
			t.Errorf("Color() = %v, want Blue — option failed to override", cp.Color())
		}
		if cp.FontSize() != unit.Sp(22) {
			t.Errorf("FontSize() = %v, want Sp(22) — option failed to override", cp.FontSize())
		}
		if cp.FontWeight() != font.FontWeightBold {
			t.Errorf("FontWeight() = %v, want Bold — untouched field not inherited", cp.FontWeight())
		}
	})

	// A documented limitation of the sentinel pattern rather than a bug: passing the
	// sentinel as an option VALUE means "I have no opinion", because TakeOrElse keeps
	// the source wherever the override is unspecified. There is deliberately no in-band
	// way to say "explicitly set this field to unspecified" — the sentinel cannot mean
	// both "unset" and "set to unset". Clearing a field means starting from a fresh
	// style that simply never set it.
	t.Run("passing_the_sentinel_defers_to_the_source", func(t *testing.T) {
		base := TextStyleFromOptions(WithColor(graphics.ColorRed))
		cp := CopyTextStyleWithOptions(base, WithColor(graphics.ColorUnspecified))
		if cp.Color() != graphics.ColorRed {
			t.Errorf("Color() = %v, want Red — an unspecified option value must defer to the source", cp.Color())
		}
	})

	t.Run("clearing_a_field_means_not_setting_it", func(t *testing.T) {
		base := TextStyleFromOptions(WithColor(graphics.ColorRed), WithFontSize(unit.Sp(16)))
		// Rebuild from the sentinel and set only what should remain — the way to
		// express "no color" is to never apply a color option.
		cleared := CopyTextStyleWithOptions(
			TextStyleCompanion.Unspecified(),
			WithFontSize(unit.Sp(16)),
		)
		if cleared.Color() != graphics.ColorUnspecified {
			t.Errorf("Color() = %v, want Unspecified", cleared.Color())
		}
		if cleared.FontSize() != unit.Sp(16) {
			t.Errorf("FontSize() = %v, want Sp(16)", cleared.FontSize())
		}
		if base.Color() != graphics.ColorRed {
			t.Error("base must be unchanged")
		}
	})
}

// === Companion.Unspecified is the only door to the sentinel ===

func TestCompanions_Unspecified(t *testing.T) {
	if ts := TextStyleCompanion.Unspecified(); isSpecifiedTextStyle(ts) {
		t.Error("TextStyleCompanion.Unspecified() must be unspecified")
	}
	if s := SpanStyleCompanion.Unspecified(); isSpecifiedSpanStyle(s) {
		t.Error("SpanStyleCompanion.Unspecified() must be unspecified")
	}
	if p := ParagraphStyleCompanion.Unspecified(); isSpecifiedParagraphStyle(p) {
		t.Error("ParagraphStyleCompanion.Unspecified() must be unspecified")
	}

	// The sentinel is a single shared value, so identity comparison is meaningful.
	if TextStyleCompanion.Unspecified() != textStyleUnspecified {
		t.Error("TextStyleCompanion.Unspecified() must return the package sentinel")
	}

	// A derived style is never the sentinel, even when built with no options.
	if TextStyleFromOptions() == TextStyleCompanion.Unspecified() {
		t.Error("TextStyleFromOptions() must return a fresh specified style, not the sentinel")
	}
}

func TestCompanions_IsSpecifiedAndString(t *testing.T) {
	unspecified := TextStyleCompanion.Unspecified()
	specified := TextStyleFromOptions(WithColor(graphics.ColorRed))

	if TextStyleCompanion.IsSpecified(unspecified) {
		t.Error("IsSpecified(sentinel) = true, want false")
	}
	if !TextStyleCompanion.IsSpecified(specified) {
		t.Error("IsSpecified(derived) = false, want true")
	}
	if !TextStyleCompanion.IsUnspecified(unspecified) {
		t.Error("IsUnspecified(sentinel) = false, want true")
	}
	if got := TextStyleCompanion.String(unspecified); got != "TextStyle{Unspecified}" {
		t.Errorf("String(sentinel) = %q, want TextStyle{Unspecified}", got)
	}
	if TextStyleCompanion.String(specified) == "TextStyle{Unspecified}" {
		t.Error("String(derived) must not read as unspecified")
	}
}

// === nil-safety ===

func TestNilSafety(t *testing.T) {
	t.Run("copy_from_nil_does_not_panic", func(t *testing.T) {
		var nilTS *TextStyle
		ts := CopyTextStyleWithOptions(nilTS, WithColor(graphics.ColorRed))
		if ts == nil {
			t.Fatal("CopyTextStyleWithOptions(nil) returned nil")
		}
		if ts.Color() != graphics.ColorRed {
			t.Errorf("Color() = %v, want Red", ts.Color())
		}
	})

	t.Run("span_and_paragraph_copy_from_nil", func(t *testing.T) {
		var nilSpan *SpanStyle
		if s := CopySpanStyleWithOptions(nilSpan, SpanStyleWithFontWeight(font.FontWeightBold)); s == nil {
			t.Error("CopySpanStyleWithOptions(nil) returned nil")
		} else if s.FontWeight() != font.FontWeightBold {
			t.Errorf("FontWeight() = %v, want Bold", s.FontWeight())
		}

		var nilPara *ParagraphStyle
		if p := CopyParagraphStyleWithOptions(nilPara, ParagraphStyleWithTextAlign(style.TextAlignMiddle)); p == nil {
			t.Error("CopyParagraphStyleWithOptions(nil) returned nil")
		} else if p.TextAlign() != style.TextAlignMiddle {
			t.Errorf("TextAlign() = %v, want Middle", p.TextAlign())
		}
	})

	t.Run("resolve_defaults_from_nil", func(t *testing.T) {
		var nilTS *TextStyle
		resolved := ResolveTextStyleDefaults(nilTS, unit.LayoutDirectionLtr)
		if resolved == nil {
			t.Fatal("ResolveTextStyleDefaults(nil) returned nil")
		}
		if resolved.FontSize() != DefaultFontSize {
			t.Errorf("FontSize() = %v, want %v", resolved.FontSize(), DefaultFontSize)
		}
	})
}

// === Merge direction ===

// mergeX(a, b) must take b where b is specified and fall back to a elsewhere.
func TestMergeDirection(t *testing.T) {
	t.Run("b_overrides_a_and_a_survives_elsewhere", func(t *testing.T) {
		a := TextStyleFromOptions(
			WithColor(graphics.ColorRed),
			WithFontSize(unit.Sp(16)),
		)
		b := TextStyleFromOptions(
			WithColor(graphics.ColorBlue),
			WithFontWeight(font.FontWeightBold),
		)

		merged := TextStyleCompanion.Merge(a, b)

		if merged.Color() != graphics.ColorBlue {
			t.Errorf("Color() = %v, want Blue — b must win where specified", merged.Color())
		}
		if merged.FontWeight() != font.FontWeightBold {
			t.Errorf("FontWeight() = %v, want Bold — b's value lost", merged.FontWeight())
		}
		if merged.FontSize() != unit.Sp(16) {
			t.Errorf("FontSize() = %v, want Sp(16) — a's value should survive", merged.FontSize())
		}
	})

	t.Run("unspecified_b_leaves_a_untouched", func(t *testing.T) {
		a := TextStyleFromOptions(WithColor(graphics.ColorRed))
		merged := TextStyleCompanion.Merge(a, TextStyleCompanion.Unspecified())
		if merged != a {
			t.Error("Merge(a, unspecified) should return a")
		}
	})

	t.Run("unspecified_a_yields_b", func(t *testing.T) {
		b := TextStyleFromOptions(WithColor(graphics.ColorRed))
		merged := TextStyleCompanion.Merge(TextStyleCompanion.Unspecified(), b)
		if merged != b {
			t.Error("Merge(unspecified, b) should return b")
		}
	})

	t.Run("nil_layers_do_not_panic", func(t *testing.T) {
		a := TextStyleFromOptions(WithColor(graphics.ColorRed))
		if got := TextStyleCompanion.Merge(a, nil); got != a {
			t.Error("Merge(a, nil) should return a")
		}
		if got := TextStyleCompanion.Merge(nil, a); got != a {
			t.Error("Merge(nil, a) should return a")
		}
		if got := TextStyleCompanion.Merge(nil, nil); got != TextStyleCompanion.Unspecified() {
			t.Error("Merge(nil, nil) should return the sentinel")
		}
	})
}

// === unspecified == unspecified ===

// Two unspecified values are both "unset" and must compare equal; comparing NaN or
// pointer identity naively would report them as different.
func TestUnspecifiedEqualsUnspecified(t *testing.T) {
	u1 := TextStyleCompanion.Unspecified()
	u2 := TextStyleCompanion.Unspecified()

	if !TextStyleCompanion.SemanticEqual(u1, u2) {
		t.Error("SemanticEqual(unspecified, unspecified) = false, want true")
	}
	if !TextStyleCompanion.Equal(u1, u2) {
		t.Error("Equal(unspecified, unspecified) = false, want true")
	}
	if !TextStyleCompanion.Same(nil, nil) {
		t.Error("Same(nil, nil) = false, want true")
	}

	// nil and the sentinel are both "unspecified" under this package's chosen semantics.
	if !TextStyleCompanion.SemanticEqual(nil, u1) {
		t.Error("SemanticEqual(nil, unspecified) = false, want true")
	}

	// Round-tripping through Merge must not accidentally make an unspecified style
	// compare unequal to another unspecified style.
	roundTripped := TextStyleCompanion.Merge(u1, u2)
	if !TextStyleCompanion.SemanticEqual(roundTripped, u2) {
		t.Error("merge of two unspecified styles should still be unspecified-equal")
	}

	// An unspecified style and a derived style must NOT be equal.
	derived := TextStyleFromOptions(WithColor(graphics.ColorRed))
	if TextStyleCompanion.SemanticEqual(u1, derived) {
		t.Error("SemanticEqual(unspecified, derived) = true, want false")
	}
}

// === TakeOrElse ===

func TestCompanion_TakeOrElse(t *testing.T) {
	def := TextStyleFromOptions(WithColor(graphics.ColorBlue))
	custom := TextStyleFromOptions(WithColor(graphics.ColorRed))

	if got := TextStyleCompanion.TakeOrElse(TextStyleCompanion.Unspecified(), def); got != def {
		t.Error("TakeOrElse(unspecified, def) should return def")
	}
	if got := TextStyleCompanion.TakeOrElse(nil, def); got != def {
		t.Error("TakeOrElse(nil, def) should return def")
	}
	if got := TextStyleCompanion.TakeOrElse(custom, def); got != custom {
		t.Error("TakeOrElse(custom, def) should return custom")
	}
}

// === Source is never mutated ===

func TestCopyDoesNotMutateSource(t *testing.T) {
	// A derived style must survive being used as a copy source.
	source := TextStyleFromOptions(
		WithColor(graphics.ColorRed),
		WithFontSize(unit.Sp(18)),
	)
	before := TextStyleCompanion.String(source)

	_ = CopyTextStyleWithOptions(source, WithColor(graphics.ColorBlue))
	_ = CopyTextStyleWithOptions(source, WithFontSize(unit.Sp(40)))

	if after := TextStyleCompanion.String(source); after != before {
		t.Errorf("source mutated by copy:\n before: %s\n after:  %s", before, after)
	}

	// Options must not be applied in place to a non-sentinel source either.
	cp := CopyTextStyleWithOptions(source, WithColor(graphics.ColorBlue))
	if source.Color() != graphics.ColorRed {
		t.Error("source Color changed after copy")
	}
	if cp == source {
		t.Error("copy should return a distinct value")
	}

	// Copying with no options still yields a value equal to the source.
	same := TextStyleCompanion.Copy(source)
	if !TextStyleCompanion.SemanticEqual(same, source) {
		t.Error("Copy with no options should be semantically equal to the source")
	}
}
