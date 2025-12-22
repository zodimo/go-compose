package graphics

// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/ui/ui-graphics/src/commonMain/kotlin/androidx/compose/ui/graphics/Brush.kt;drc=182f9f08cf26aaa426d03f39b6dcef1c33e9f41b;l=35

// Brush is the interface for all brush types used for drawing.
type Brush interface {
	isBrush()
}

// SolidColor is a Brush that represents a single solid color.
// SolidColor brushes are stored as regular Colors in TextForegroundStyle.
type SolidColor struct {
	Value Color
}

func (s SolidColor) isBrush() {}

// NewSolidColor creates a new SolidColor brush from a Color.
func NewSolidColor(color Color) SolidColor {
	return SolidColor{Value: color}
}

// ShaderBrush is a Brush implementation that wraps a shader.
// The shader can be lazily created based on a given size.
type ShaderBrush interface {
	Brush
	isShaderBrush()
}

// shaderBrushImpl is an internal implementation of ShaderBrush.
type shaderBrushImpl struct{}

func (s shaderBrushImpl) isBrush()       {}
func (s shaderBrushImpl) isShaderBrush() {}

// ShaderBrushForTest is an exported ShaderBrush implementation for cross-package testing.
// Use NewShaderBrushForTest() to create instances in tests.
type ShaderBrushForTest struct{}

func (s ShaderBrushForTest) isBrush()       {}
func (s ShaderBrushForTest) isShaderBrush() {}

// NewShaderBrushForTest creates a ShaderBrush for testing purposes.
// This is intended for use in test files across packages.
func NewShaderBrushForTest() ShaderBrush {
	return ShaderBrushForTest{}
}

// IsSolidColor returns true if the brush is a SolidColor.
func IsSolidColor(b Brush) bool {
	_, ok := b.(SolidColor)
	return ok
}

// IsShaderBrush returns true if the brush is a ShaderBrush.
func IsShaderBrush(b Brush) bool {
	_, ok := b.(ShaderBrush)
	return ok
}

// AsSolidColor returns the brush as a SolidColor if it is one, otherwise nil.
func AsSolidColor(b Brush) *SolidColor {
	if sc, ok := b.(SolidColor); ok {
		return &sc
	}
	return nil
}

// AsShaderBrush returns the brush as a ShaderBrush if it is one, otherwise nil.
func AsShaderBrush(b Brush) ShaderBrush {
	if sb, ok := b.(ShaderBrush); ok {
		return sb
	}
	return nil
}
