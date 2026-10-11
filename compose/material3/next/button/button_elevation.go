package button

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/ui/unit"
)

// ButtonElevationUnspecified is the sentinel value for unspecified ButtonElevation
var ButtonElevationUnspecified = &ButtonElevation{
	DefaultElevation:  unit.DpCompanion.Unspecified(),
	PressedElevation:  unit.DpCompanion.Unspecified(),
	FocusedElevation:  unit.DpCompanion.Unspecified(),
	HoveredElevation:  unit.DpCompanion.Unspecified(),
	DisabledElevation: unit.DpCompanion.Unspecified(),
}

// ButtonElevation holds elevation configuration for button states
type ButtonElevation struct {
	DefaultElevation  unit.Dp
	PressedElevation  unit.Dp
	FocusedElevation  unit.Dp
	HoveredElevation  unit.Dp
	DisabledElevation unit.Dp
}

// IsSpecifiedButtonElevation returns true if e is specified (not nil and not the sentinel)
func IsSpecifiedButtonElevation(e *ButtonElevation) bool {
	return e != nil && e != ButtonElevationUnspecified
}

// TakeOrElseButtonElevation returns e if specified, otherwise returns defaultElevation
func TakeOrElseButtonElevation(e, defaultElevation *ButtonElevation) *ButtonElevation {
	if e == nil || e == ButtonElevationUnspecified {
		return defaultElevation
	}
	return e
}

// MergeButtonElevation merges two ButtonElevation, preferring b's specified values over a's
func MergeButtonElevation(a, b *ButtonElevation) *ButtonElevation {
	a = CoalesceButtonElevation(a, ButtonElevationUnspecified)
	b = CoalesceButtonElevation(b, ButtonElevationUnspecified)

	if a == ButtonElevationUnspecified {
		return b
	}
	if b == ButtonElevationUnspecified {
		return a
	}

	// Both are custom: allocate new merged style
	return &ButtonElevation{
		DefaultElevation:  unit.DpCompanion.TakeOrElse(b.DefaultElevation, a.DefaultElevation),
		PressedElevation:  unit.DpCompanion.TakeOrElse(b.PressedElevation, a.PressedElevation),
		FocusedElevation:  unit.DpCompanion.TakeOrElse(b.FocusedElevation, a.FocusedElevation),
		HoveredElevation:  unit.DpCompanion.TakeOrElse(b.HoveredElevation, a.HoveredElevation),
		DisabledElevation: unit.DpCompanion.TakeOrElse(b.DisabledElevation, a.DisabledElevation),
	}
}

// StringButtonElevation returns a string representation of ButtonElevation
func StringButtonElevation(e *ButtonElevation) string {
	if !IsSpecifiedButtonElevation(e) {
		return "ButtonElevation{Unspecified}"
	}

	return fmt.Sprintf(
		"ButtonElevation{DefaultElevation: %s, PressedElevation: %s, FocusedElevation: %s, HoveredElevation: %s, DisabledElevation: %s}",
		e.DefaultElevation,
		e.PressedElevation,
		e.FocusedElevation,
		e.HoveredElevation,
		e.DisabledElevation,
	)
}

// CoalesceButtonElevation returns ptr if not nil, otherwise returns def
func CoalesceButtonElevation(ptr, def *ButtonElevation) *ButtonElevation {
	if ptr == nil {
		return def
	}
	return ptr
}

// SameButtonElevation returns true if a and b are the same pointer or both unspecified
func SameButtonElevation(a, b *ButtonElevation) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return b == ButtonElevationUnspecified
	}
	if b == nil {
		return a == ButtonElevationUnspecified
	}
	return a == b
}

// SemanticEqualButtonElevation checks field-by-field equality
func SemanticEqualButtonElevation(a, b *ButtonElevation) bool {
	a = CoalesceButtonElevation(a, ButtonElevationUnspecified)
	b = CoalesceButtonElevation(b, ButtonElevationUnspecified)

	return a.DefaultElevation == b.DefaultElevation &&
		a.PressedElevation == b.PressedElevation &&
		a.FocusedElevation == b.FocusedElevation &&
		a.HoveredElevation == b.HoveredElevation &&
		a.DisabledElevation == b.DisabledElevation
}

// EqualButtonElevation returns true if a and b are semantically equal
func EqualButtonElevation(a, b *ButtonElevation) bool {
	if SameButtonElevation(a, b) {
		return true
	}
	return SemanticEqualButtonElevation(a, b)
}

// ButtonElevationOption is a functional option for CopyButtonElevation
type ButtonElevationOption func(*ButtonElevation)

// WithDefaultElevation sets the default elevation option
func WithDefaultElevation(e unit.Dp) ButtonElevationOption {
	return func(o *ButtonElevation) {
		o.DefaultElevation = e
	}
}

// WithPressedElevation sets the pressed elevation option
func WithPressedElevation(e unit.Dp) ButtonElevationOption {
	return func(o *ButtonElevation) {
		o.PressedElevation = e
	}
}

// WithFocusedElevation sets the focused elevation option
func WithFocusedElevation(e unit.Dp) ButtonElevationOption {
	return func(o *ButtonElevation) {
		o.FocusedElevation = e
	}
}

// WithHoveredElevation sets the hovered elevation option
func WithHoveredElevation(e unit.Dp) ButtonElevationOption {
	return func(o *ButtonElevation) {
		o.HoveredElevation = e
	}
}

// WithDisabledElevation sets the disabled elevation option
func WithDisabledElevation(e unit.Dp) ButtonElevationOption {
	return func(o *ButtonElevation) {
		o.DisabledElevation = e
	}
}

// CopyButtonElevation creates a copy with optional modifications
func CopyButtonElevation(e *ButtonElevation, options ...ButtonElevationOption) *ButtonElevation {
	opt := *ButtonElevationUnspecified

	for _, option := range options {
		option(&opt)
	}

	e = CoalesceButtonElevation(e, ButtonElevationUnspecified)

	return &ButtonElevation{
		DefaultElevation:  unit.DpCompanion.TakeOrElse(opt.DefaultElevation, e.DefaultElevation),
		PressedElevation:  unit.DpCompanion.TakeOrElse(opt.PressedElevation, e.PressedElevation),
		FocusedElevation:  unit.DpCompanion.TakeOrElse(opt.FocusedElevation, e.FocusedElevation),
		HoveredElevation:  unit.DpCompanion.TakeOrElse(opt.HoveredElevation, e.HoveredElevation),
		DisabledElevation: unit.DpCompanion.TakeOrElse(opt.DisabledElevation, e.DisabledElevation),
	}
}
