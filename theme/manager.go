package theme

import (
	"sync"

	"gioui.org/layout"
	"git.sr.ht/~schnwalter/gio-mw/token"
	"git.sr.ht/~schnwalter/gio-mw/wdk"
	"github.com/zodimo/go-compose/compose/ui/graphics"
)

var themeManagerSingleton ThemeManager

var ColorHelper ThemeColorHelper = nil

func init() {
	ColorHelper = newTheColorHelper()
}

type ThemeColorHelper interface {
	ColorSelector() *ColorRoleDescriptors
	SpecificColor(color graphics.Color) ColorDescriptor
	UnspecifiedColor() ColorDescriptor
}

var _ ThemeColorHelper = (*themeColorHelper)(nil)

type themeColorHelper struct {
	roleDescriptors ColorRoleDescriptors
}

func (tch themeColorHelper) ColorSelector() *ColorRoleDescriptors {
	return &tch.roleDescriptors
}
func (tch themeColorHelper) SpecificColor(color graphics.Color) ColorDescriptor {
	return SpecificColor(color)
}
func (tch themeColorHelper) UnspecifiedColor() ColorDescriptor {
	return ColorUnspecified
}
func newTheColorHelper() ThemeColorHelper {
	return themeColorHelper{
		roleDescriptors: NewColorRoleDescriptors(),
	}
}

// Runtime
type ThemeManager interface {
	materialTheme() *BasicTheme
	setMaterialTheme(theme *BasicTheme)

	Material3ThemeInit(gtx any) any
	setMaterial3Theme(gtx layout.Context, theme *token.Theme)
	getMaterial3Theme() *token.Theme

	ThemeColorResolver
}

var _ ThemeManager = (*themeManager)(nil)

type themeManager struct {
	mu                   sync.RWMutex
	basicTheme           *BasicTheme
	tokenTheme           *token.Theme
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

func (tm *themeManager) materialTheme() *BasicTheme {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	return tm.basicTheme
}

func (tm *themeManager) setMaterialTheme(theme *BasicTheme) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.basicTheme = theme
}

func (tm *themeManager) Material3ThemeInit(gtx any) any {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	g, ok := gtx.(layout.Context)
	if !ok {
		panic("theme.Material3ThemeInit: gtx must be gioui.org/layout.Context")
	}

	if tm.tokenTheme == nil {
		tm.tokenTheme = defaultMaterial3Theme(g)
	}

	g.Values = make(map[string]any)
	wdk.InitMaterialThemeInContext(g, tm.tokenTheme)
	return g
}

func (tm *themeManager) setMaterial3Theme(gtx layout.Context, theme *token.Theme) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.tokenTheme = theme
}

func (tm *themeManager) getMaterial3Theme() *token.Theme {
	tm.mu.RLock()
	defer tm.mu.RUnlock()
	if tm.tokenTheme == nil {
		panic("material3Theme is nil")
	}
	return tm.tokenTheme
}

func (tm *themeManager) ResolveColorDescriptor(desc ColorDescriptor) ThemeColor {
	return tm.themeColorResolver.ResolveColorDescriptor(desc)
}

func (tm *themeManager) ColorRoleDescriptors() ColorRoleDescriptors {
	return tm.colorRoleDescriptors
}

func (tm *themeManager) ColorDescriptor(color graphics.Color) ColorDescriptor {
	return SpecificColor(color)
}

// Deprecated, use local providers
func GetThemeManager() ThemeManager {
	return themeManagerSingleton
}

func init() {
	themeManagerSingleton = newThemeManager(defaultMaterialTheme())
}
