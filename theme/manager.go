package theme

import (
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

	ColorRoleDescriptors() ColorRoleDescriptors

	ThemeColorResolver
}

var _ ThemeManager = (*themeManager)(nil)

type themeManager struct {
	mu                   sync.RWMutex
	basicTheme           *BasicTheme
	theme                *Theme
	themeColorResolver   ThemeColorResolver
	colorRoleDescriptors ColorRoleDescriptors
}

func newThemeManager(theme *BasicTheme) ThemeManager {
	tm := &themeManager{
		basicTheme:           theme,
		colorRoleDescriptors: NewColorRoleDescriptors(),
	}
	tm.themeColorResolver = newThemeColorResolver(tm)
	return tm
}

func (tm *themeManager) MaterialTheme() *material.Theme {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.basicTheme

}
func (tm *themeManager) SetMaterialTheme(theme *material.Theme) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.basicTheme = theme
}

func (tm *themeManager) Material3ThemeInit(gtx layout.Context) layout.Context {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	if tm.theme == nil {
		tm.theme = defaultMaterial3Theme(gtx)
	}

	gtx.Values = make(map[string]any)
	wdk.InitMaterialThemeInContext(gtx, tm.theme)
	return gtx

}
func (tm *themeManager) SetMaterial3Theme(gtx layout.Context, theme *token.Theme) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.theme = theme
}

func (tm *themeManager) GetMaterial3Theme() *token.Theme {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	if tm.theme == nil {
		panic("material3Theme is nil")
	}
	return tm.theme
}

func (tm *themeManager) ResolveColorDescriptor(desc ThemeColorDescriptor) ThemeColor {
	return tm.themeColorResolver.ResolveColorDescriptor(desc)
}

func (tm *themeManager) ColorRoleDescriptors() ColorRoleDescriptors {
	return tm.colorRoleDescriptors
}

func GetThemeManager() ThemeManager {
	return themeManagerSingleton
}

func init() {
	themeManagerSingleton = newThemeManager(defaultMaterialTheme())
}
