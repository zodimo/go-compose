package textinput

import (
	gioMaterial "gioui.org/widget/material"
)

// Theme wraps gioui.org/widget/material.Theme to hide gioui types from public API.
type Theme struct {
	T *gioMaterial.Theme
}

// NewTheme wraps a gio material.Theme.
func NewTheme(t *gioMaterial.Theme) *Theme {
	return &Theme{T: t}
}
