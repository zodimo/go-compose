package theme

import (
	"image/color"
	"sync"

	"gioui.org/layout"
	"gioui.org/widget/material"
	"git.sr.ht/~schnwalter/gio-mw/token"
	"git.sr.ht/~schnwalter/gio-mw/wdk"
)

var themeManagerSingleton ThemeManager

type ThemeManager interface {
	MaterialTheme() *material.Theme
	SetMaterialTheme(theme *material.Theme)

	Material3ThemeInit(gtx layout.Context) layout.Context
	SetMaterial3Theme(gtx layout.Context, theme *token.Theme)
	GetMaterial3Theme() *token.Theme

	ThemeColorResolver() ThemeColorResolver
}

var _ ThemeManager = (*themeManager)(nil)

type themeManager struct {
	mu             sync.RWMutex
	materialTheme  *material.Theme
	material3Theme *token.Theme
}

func newThemeManager(materialTheme *material.Theme) ThemeManager {
	return &themeManager{
		materialTheme: materialTheme,
	}
}

func (tm *themeManager) MaterialTheme() *material.Theme {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.materialTheme

}
func (tm *themeManager) SetMaterialTheme(theme *material.Theme) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.materialTheme = theme
}

func (tm *themeManager) Material3ThemeInit(gtx layout.Context) layout.Context {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	if tm.material3Theme == nil {
		tm.material3Theme = defaultMaterial3Theme(gtx)
	}

	gtx.Values = make(map[string]any)
	wdk.InitMaterialThemeInContext(gtx, tm.material3Theme)
	return gtx

}
func (tm *themeManager) SetMaterial3Theme(gtx layout.Context, theme *token.Theme) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.material3Theme = theme
}

func (tm *themeManager) GetMaterial3Theme() *token.Theme {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	if tm.material3Theme == nil {
		panic("material3Theme is nil")
	}
	return tm.material3Theme
}

func (tm *themeManager) ThemeColorResolver() ThemeColorResolver {
	return newThemeColor(tm)
}

func GetThemeManager() ThemeManager {
	return themeManagerSingleton
}

func init() {
	themeManagerSingleton = newThemeManager(defaultMaterialTheme())

}

type ThemeColorResolver interface {
	Material3(func(theme *token.Theme) color.Color) color.Color
	Material(func(theme *material.Theme) color.Color) color.Color
}

type themeColorResolver struct {
	tm ThemeManager
}

func (tc *themeColorResolver) Material3(reader func(theme *token.Theme) color.Color) color.Color {
	return reader(tc.tm.GetMaterial3Theme())
}

func (tc *themeColorResolver) Material(reader func(theme *material.Theme) color.Color) color.Color {
	return reader(tc.tm.MaterialTheme())
}
func newThemeColor(tm ThemeManager) ThemeColorResolver {
	return &themeColorResolver{tm: tm}
}
