package button

import (
	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/foundation"
	"github.com/zodimo/go-compose/compose/foundation/layout"
	"github.com/zodimo/go-compose/compose/material3"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/graphics/shape"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

type ButtonColors struct {
	ContainerColor         graphics.Color
	ContentColor           graphics.Color
	DisabledContainerColor graphics.Color
	DisabledContentColor   graphics.Color
}

type ButtonElevation struct {
	DefaultElevation  unit.Dp
	PressedElevation  unit.Dp
	FocusedElevation  unit.Dp
	HoveredElevation  unit.Dp
	DisabledElevation unit.Dp
}

type ButtonShapes struct {
	Shape        shape.Shape
	PressedShape shape.Shape
}

func DefaultButtonOptions() ButtonOptions {
	return ButtonOptions{
		Modifier: EmptyModifier,
	}
}

var ButtonDefaults = buttonDefaults{}

type buttonDefaults struct{}

var (
	// Outlined Button Tokens
	outlinedButtonContainerColor           = graphics.ColorTransparent
	outlinedButtonDisabledContainerColor   = graphics.ColorTransparent
	outlinedButtonOutlineWidth             = unit.Dp(1)
	outlinedButtonDisabledContainerOpacity = 0.12
)

// OutlinedButtonColors returns the default colors for an OutlinedButton
func (buttonDefaults) OutlinedButtonColors(c compose.Composer) ButtonColors {
	colorScheme := material3.Theme(c).ColorScheme()
	return ButtonColors{
		ContainerColor:         outlinedButtonContainerColor,
		ContentColor:           colorScheme.Primary.Color,
		DisabledContainerColor: outlinedButtonDisabledContainerColor,
		DisabledContentColor:   colorScheme.Surface.OnColor.Copy(graphics.CopyWithAlpha(0.38)),
	}
}

// OutlinedButtonElevation returns the default elevation for an OutlinedButton
func (buttonDefaults) OutlinedButtonElevation(c compose.Composer) ButtonElevation {
	return ButtonElevation{
		DefaultElevation:  0,
		PressedElevation:  0,
		FocusedElevation:  0,
		HoveredElevation:  0,
		DisabledElevation: 0,
	}
}

// OutlinedButtonBorder returns the default border for an OutlinedButton
func (buttonDefaults) OutlinedButtonBorder(c compose.Composer, enabled bool) foundation.BorderStroke {
	colorScheme := material3.Theme(c).ColorScheme()
	var borderColor graphics.Color
	if enabled {
		borderColor = colorScheme.Outline
	} else {
		borderColor = colorScheme.Outline.Copy(
			graphics.CopyWithAlpha(0.12),
		)
	}
	return foundation.BorderStroke{
		Width: outlinedButtonOutlineWidth,
		Brush: graphics.SolidColor{Value: borderColor},
	}
}

// OutlinedButtonShape returns the default shape for an OutlinedButton
func (buttonDefaults) OutlinedButtonShape(c compose.Composer) shape.Shape {
	// M3 Default is specific rounded corner, here using a placeholder or theme shape if available
	// M3 token: ShapeScale.CornerFull (Circle) or CornerSmall?
	// Kotlin defaults use ButtonSmallTokens.ContainerShapeRound which is ShapeKeyTokens.CornerFull -> Circle
	return shape.CircleShape
}

var (
	buttonVerticalPadding   = unit.Dp(8)
	buttonHorizontalPadding = unit.Dp(24)
)

// ContentPadding returns the default content padding for a Button
func (buttonDefaults) ContentPadding() layout.PaddingValues {
	return layout.PaddingValues{
		Start:  buttonHorizontalPadding,
		Top:    buttonVerticalPadding,
		End:    buttonHorizontalPadding,
		Bottom: buttonVerticalPadding,
	}
}
