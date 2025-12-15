package theme

import (
	"image/color"
)

type ResolvableColor interface {
	ResolveColorDescriptor(colorDesc ThemeColorDescriptor) ThemeColor
}

// ResolvedColor is a color that has been resolved from a ThemeColorDescriptor

type ThemeColor interface {
	AsHex() string
	AsNRGBA() color.NRGBA
	SetOpacity(opacity OpacityLevel) ThemeColor
	AsTokenColor() TokenColor
}

type ThemeColorUpdateActions int

const (
	ThemeColorUpdateActionsSetOpacity ThemeColorUpdateActions = iota
)

type ThemeColorUpdate interface {
	Action() ThemeColorUpdateActions
	Value() any
	Compare(other ThemeColorUpdate) bool
	isThemeColorUpdate() bool
}

var _ ThemeColorUpdate = (*themeColorUpdate)(nil)

type themeColorUpdate struct {
	action ThemeColorUpdateActions
	value  any
}

func (u themeColorUpdate) Action() ThemeColorUpdateActions {
	return u.action
}

func (u themeColorUpdate) Value() any {
	return u.value
}

func (u themeColorUpdate) Compare(other ThemeColorUpdate) bool {
	return u.action == other.Action() && u.value == other.Value()
}

func (u themeColorUpdate) isThemeColorUpdate() bool {
	return true
}

var _ ThemeColorUpdate = (*ThemeColorUpdateTyped[any])(nil)

type ThemeColorUpdateTyped[T comparable] struct {
	action ThemeColorUpdateActions
	value  T
}

func (u ThemeColorUpdateTyped[T]) Action() ThemeColorUpdateActions {
	return u.action
}

func (u ThemeColorUpdateTyped[T]) Value() T {
	return u.value
}

func (u ThemeColorUpdateTyped[T]) Compare(other ThemeColorUpdate) bool {
	return u.action == other.Action() && u.value == other.Value()
}

func (u ThemeColorUpdateTyped[T]) isThemeColorUpdate() bool {
	return true
}

func (u ThemeColorUpdateTyped[T]) Any() ThemeColorUpdate {
	return themeColorUpdate{
		action: u.action,
		value:  u.value,
	}
}

type ThemeColorDescriptor struct {
	color     color.NRGBA
	colorRole ColorRole
	isColor   bool
	updates   []ThemeColorUpdate
}

func (t ThemeColorDescriptor) AppendUpdate(update ThemeColorUpdate) ThemeColorDescriptor {
	t = t.initUpdates()

	t.updates = append(t.updates, update)
	return t
}

func (t ThemeColorDescriptor) SetOpacity(opacity OpacityLevel) ThemeColorDescriptor {

	update := SetOpacity(opacity)
	return t.AppendUpdate(update.Any())
}

func (t ThemeColorDescriptor) initUpdates() ThemeColorDescriptor {
	if t.updates == nil {
		t.updates = make([]ThemeColorUpdate, 0)
	}
	return t
}

func (t ThemeColorDescriptor) Compare(other ThemeColorDescriptor) bool {
	if len(t.updates) != len(other.updates) {
		return false
	}
	//compare updates
	for i := range t.updates {
		if !t.updates[i].Compare(other.updates[i]) {
			return false
		}
	}
	//compare color
	return t.color == other.color && t.colorRole == other.colorRole && t.isColor == other.isColor
}

func (t ThemeColorDescriptor) Resolve() ThemeColor {
	return GetThemeManager().ResolveColorDescriptor(t)
}

func SetOpacity(value OpacityLevel) *ThemeColorUpdateTyped[OpacityLevel] {
	return &ThemeColorUpdateTyped[OpacityLevel]{
		action: ThemeColorUpdateActionsSetOpacity,
		value:  value,
	}
}

func GetOpacity(update ThemeColorUpdate) OpacityLevel {
	if update.Action() != ThemeColorUpdateActionsSetOpacity {
		panic("update is not a set opacity update")
	}
	return update.Value().(OpacityLevel)
}
