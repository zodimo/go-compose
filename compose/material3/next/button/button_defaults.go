package button

import (
	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/foundation"
	"github.com/zodimo/go-compose/compose/foundation/layout"
	"github.com/zodimo/go-compose/compose/material3"
	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/graphics/shape"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

func DefaultButtonOptions() ButtonOptions {
	return ButtonOptions{
		Modifier:  ui.EmptyModifier,
		Colors:    ButtonColorsUnspecified,
		Enabled:   true,
		Shape:     shape.ShapeCompanion.Unspecified(),
		Elevation: ButtonElevationUnspecified,
		// Border:    foundation.BorderStrokeUnspecified,

		// Border         foundation.BorderStroke
		// ContentPadding layout.PaddingValues
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

	// Button Shape Tokens (M3 spec from ButtonSmallTokens)
	// ContainerShapeRound = ShapeKeyTokens.CornerFull -> CircleShape
	buttonContainerShapeRound = shape.CircleShape
	// PressedContainerShape = ShapeKeyTokens.CornerSmall -> RoundedCornerShape(8.dp)
	buttonPressedContainerShape = &shape.RoundedCornerShape{Radius: unit.Dp(8)}
)

// OutlinedButtonColors returns the default colors for an OutlinedButton
func (buttonDefaults) OutlinedButtonColors(c compose.Composer) *ButtonColors {
	colorScheme := material3.Theme(c).ColorScheme()
	return &ButtonColors{
		ContainerColor:         outlinedButtonContainerColor,
		ContentColor:           colorScheme.Primary,
		DisabledContainerColor: outlinedButtonDisabledContainerColor,
		DisabledContentColor:   colorScheme.OnSurface.Copy(graphics.CopyWithAlpha(0.38)),
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
func (buttonDefaults) OutlinedButtonBorder(c compose.Composer, enabled bool) *foundation.BorderStroke {
	colorScheme := material3.Theme(c).ColorScheme()
	var borderColor graphics.Color
	if enabled {
		borderColor = colorScheme.Outline
	} else {
		borderColor = colorScheme.Outline.Copy(
			graphics.CopyWithAlpha(0.12),
		)
	}
	return foundation.NewBorderStrokeWithColor(outlinedButtonOutlineWidth, borderColor)
}

// OutlinedButtonShape returns the default shape for an OutlinedButton.
// M3 token: ButtonSmallTokens.ContainerShapeRound = ShapeKeyTokens.CornerFull
func (buttonDefaults) OutlinedButtonShape(c compose.Composer) shape.Shape {
	return buttonContainerShapeRound
}

// OutlinedButtonShapes returns the default ButtonShapes for an OutlinedButton,
// including both default and pressed shapes for expressive shape morphing.
// M3 tokens:
//   - Shape: ButtonSmallTokens.ContainerShapeRound = ShapeKeyTokens.CornerFull (CircleShape)
//   - PressedShape: ButtonSmallTokens.PressedContainerShape = ShapeKeyTokens.CornerSmall (8dp rounded)
func (buttonDefaults) OutlinedButtonShapes(c compose.Composer) *ButtonShapes {
	return &ButtonShapes{
		Shape:        buttonContainerShapeRound,
		PressedShape: buttonPressedContainerShape,
	}
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
