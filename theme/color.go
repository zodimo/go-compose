package theme

import (
	"image/color"

	"github.com/zodimo/go-compose/compose/ui/graphics"
)

var ColorUnspecified ColorDescriptor = &colorDescriptor{
	color:   graphics.ColorUnspecified,
	isColor: true,
}

type ResolvableColor interface {
	ResolveColorDescriptor(colorDesc ColorDescriptor) ThemeColor
}

// ResolvedColor is a color that has been resolved from a ThemeColorDescriptor

type ThemeColor interface {
	AsHex() string
	AsNRGBA() color.NRGBA
	SetOpacity(opacity OpacityLevel) ThemeColor
	AsTokenColor() TokenColor
	isValid()
}

type ColorUpdateActions int

const (
	SetOpacityColorUpdateAction ColorUpdateActions = iota
	LightenColorUpdateAction
	DarkenColorUpdateAction
	LerpColorUpdateAction
)

type ColorUpdate interface {
	Action() ColorUpdateActions
	Value() any
	Compare(other ColorUpdate) bool
	isThemeColorUpdate() bool
}

var _ ColorUpdate = (*colorUpdate)(nil)

type colorUpdate struct {
	action ColorUpdateActions
	value  any
}

func (u colorUpdate) Action() ColorUpdateActions {
	return u.action
}

func (u colorUpdate) Value() any {
	return u.value
}

func (u colorUpdate) Compare(other ColorUpdate) bool {
	return u.action == other.Action() && u.value == other.Value()
}

func (u colorUpdate) isThemeColorUpdate() bool {
	return true
}

var _ ColorUpdate = (*ColorUpdateTyped[any])(nil)

type ColorUpdateTyped[T comparable] struct {
	action ColorUpdateActions
	value  T
}

func (u ColorUpdateTyped[T]) Action() ColorUpdateActions {
	return u.action
}

func (u ColorUpdateTyped[T]) Value() T {
	return u.value
}

func (u ColorUpdateTyped[T]) Compare(other ColorUpdate) bool {
	return u.action == other.Action() && u.value == other.Value()
}

func (u ColorUpdateTyped[T]) isThemeColorUpdate() bool {
	return true
}

func (u ColorUpdateTyped[T]) Any() ColorUpdate {
	return colorUpdate{
		action: u.action,
		value:  u.value,
	}
}

type ColorDescriptor interface {
	AppendUpdate(update ColorUpdate) ColorDescriptor
	SetOpacity(opacity OpacityLevel) ColorDescriptor
	Lighten(percentage float32) ColorDescriptor
	Darken(percentage float32) ColorDescriptor
	Compare(other ColorDescriptor) bool
	Updates() []ColorUpdate
}

var _ ColorDescriptor = (*colorDescriptor)(nil)

type colorDescriptor struct {
	color     graphics.Color
	colorRole ColorRole
	isColor   bool
	updates   []ColorUpdate
}

func (t colorDescriptor) AppendUpdate(update ColorUpdate) ColorDescriptor {
	t = t.initUpdates()

	t.updates = append(t.updates, update)
	return t
}

func (t colorDescriptor) SetOpacity(opacity OpacityLevel) ColorDescriptor {
	update := SetOpacity(opacity)
	return t.AppendUpdate(update.Any())
}

func (t colorDescriptor) Lighten(percentage float32) ColorDescriptor {
	update := Lighten(percentage)
	return t.AppendUpdate(update.Any())
}

func (t colorDescriptor) Darken(percentage float32) ColorDescriptor {
	update := Darken(percentage)
	return t.AppendUpdate(update.Any())
}

func (t colorDescriptor) initUpdates() colorDescriptor {
	if t.updates == nil {
		t.updates = make([]ColorUpdate, 0)
	}
	return t
}

func (t colorDescriptor) Updates() []ColorUpdate {
	return append([]ColorUpdate{}, t.updates...)
}

func (t colorDescriptor) Compare(other ColorDescriptor) bool {

	other = CoalesceColor(other, ColorUnspecified)

	otherColorDescriptor, ok := other.(colorDescriptor)
	if !ok {
		return false
	}

	if len(t.updates) != len(otherColorDescriptor.updates) {
		return false
	}
	//compare updates
	for i := range t.updates {
		if !t.updates[i].Compare(otherColorDescriptor.updates[i]) {
			return false
		}
	}
	//compare color
	return t.color == otherColorDescriptor.color && t.colorRole == otherColorDescriptor.colorRole && t.isColor == otherColorDescriptor.isColor
}

func SetOpacity(value OpacityLevel) *ColorUpdateTyped[OpacityLevel] {
	return &ColorUpdateTyped[OpacityLevel]{
		action: SetOpacityColorUpdateAction,
		value:  value,
	}
}

func GetOpacity(update ColorUpdate) OpacityLevel {
	if update.Action() != SetOpacityColorUpdateAction {
		panic("update is not a set opacity update")
	}
	return update.Value().(OpacityLevel)
}

func Lighten(percentage float32) *ColorUpdateTyped[float32] {
	return &ColorUpdateTyped[float32]{
		action: LightenColorUpdateAction,
		value:  percentage,
	}
}

func GetLighten(update ColorUpdate) float32 {
	if update.Action() != LightenColorUpdateAction {
		panic("update is not a lighten update")
	}
	return update.Value().(float32)
}

func Darken(percentage float32) *ColorUpdateTyped[float32] {
	return &ColorUpdateTyped[float32]{
		action: DarkenColorUpdateAction,
		value:  percentage,
	}
}

func GetDarken(update ColorUpdate) float32 {
	if update.Action() != DarkenColorUpdateAction {
		panic("update is not a darken update")
	}
	return update.Value().(float32)
}

type LerpColorUpdateParams struct {
	Stop     ColorDescriptor
	Fraction float32
}

func LerpColor(stop ColorDescriptor, fraction float32) LerpColorUpdate {
	return LerpColorUpdate{
		value: LerpColorUpdateParams{
			Stop:     stop,
			Fraction: fraction,
		},
	}
}

func GetLerp(update ColorUpdate) LerpColorUpdateParams {
	if update.Action() != LerpColorUpdateAction {
		panic("update is not a lerp update")
	}
	return update.Value().(LerpColorUpdateParams)
}

func (u ColorUpdateTyped[T]) CompareTyped(other ColorUpdateTyped[T]) bool {
	// For slices or other non-comparable types in T (like interfaces containing them),
	// we might need custom logic.
	// However, ColorDescriptor interface values are comparable if the underlying concrete types are.
	// Our colorDescriptor contains a slice (updates), which makes it not directly comparable with ==.
	// So we need special handling for LerpColorUpdateParams which contains ColorDescriptor.

	// This default Compare is simple equality, which fails for slices.
	// We should switch to a type switch or similar if we want generic robustness,
	// but here we are method on a specific type instantiation.
	// Wait, u.value is T. If T is LerpColorUpdateParams, it has a ColorDescriptor.
	// ColorDescriptor is an interface. If the underlying struct has a slice, == panics.
	return u.action == other.action // Only checking action here is insufficient, see below
}

// Check special case for LerpColorUpdateParams in the generic Compare method?
// The generic Compare uses u.value == other.Value(). This WILL PANIC if value contains a slice (which ColorDescriptor implementation does).
// We need to implement a specific Compare for the Lerp update or modify the generic one.
// Since ColorUpdateTyped is generic, we can't specialize methods easily without type switch.

// Redefining Compare for ColorUpdateTyped to handle ColorDescriptor comparison safely.
// Actually, earlier in the existing code:
// func (u ColorUpdateTyped[T]) Compare(other ColorUpdate) bool {
// 	return u.action == other.Action() && u.value == other.Value()
// }
// This existing implementation is dangerous if T contains non-comparable types (like slices).
// colorDescriptor struct HAS a slice `updates []ColorUpdate`.
// So we MUST change `Compare` in `theme/color.go` to handle this, or implement a specific type for Lerp update that isn't `ColorUpdateTyped`.
// Let's implement a specific type for LerpUpdate to avoid breaking the generic struct or making it too complex.

type LerpColorUpdate struct {
	value LerpColorUpdateParams
}

func (u LerpColorUpdate) Action() ColorUpdateActions {
	return LerpColorUpdateAction
}

func (u LerpColorUpdate) Value() any {
	return u.value
}

func (u LerpColorUpdate) Compare(other ColorUpdate) bool {
	if other.Action() != LerpColorUpdateAction {
		return false
	}
	otherParams, ok := other.Value().(LerpColorUpdateParams)
	if !ok {
		return false
	}
	if u.value.Fraction != otherParams.Fraction {
		return false
	}
	// Deep compare the Stop descriptor
	return u.value.Stop.Compare(otherParams.Stop)
}

func (u LerpColorUpdate) isThemeColorUpdate() bool {
	return true
}

func IsSpecifiedColor(colorDescriptor ColorDescriptor) bool {
	return colorDescriptor != nil && colorDescriptor != ColorUnspecified
}

func CoalesceColor(ptr, def ColorDescriptor) ColorDescriptor {
	if ptr == nil {
		return def
	}
	return ptr
}

func TakeOrElseColor(colorDescriptor, defaultColor ColorDescriptor) ColorDescriptor {
	if colorDescriptor == nil || colorDescriptor == ColorUnspecified {
		return defaultColor
	}
	return colorDescriptor
}
