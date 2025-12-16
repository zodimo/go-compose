package theme

import (
	"image/color"
)

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
	Compare(other ColorDescriptor) bool
	Updates() []ColorUpdate
}

var _ ColorDescriptor = (*colorDescriptor)(nil)

type colorDescriptor struct {
	color     color.NRGBA
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
