package slider

import (
	"fmt"
	"image"

	"github.com/zodimo/go-compose/compose/material3"
	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/internal/layoutnode"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget"

	gioUnit "gioui.org/unit"
)

const SliderNodeID = "Material3Slider"

// Slider is a Material 3 slider component.
//
// value is the current value of the slider.
// onValueChange is called when the value changes.
// options provides optional configuration.
func Slider(value float32, onValueChange func(float32), options ...SliderOption) Composable {

	opts := DefaultSliderOptions()
	for _, opt := range options {
		opt(&opts)
	}

	return func(c Composer) Composer {

		opts.Colors = resolveSliderColors(c, opts.Colors)

		// Calculate mapped value [0, 1] for internal widget
		rangeDiff := opts.ValueRange.Max - opts.ValueRange.Min
		internalValue := float32(0)
		if rangeDiff > 0 {
			internalValue = (value - opts.ValueRange.Min) / rangeDiff
		}

		key := c.GenerateID()
		path := c.GetPath()
		floatStatePath := fmt.Sprintf("%d/%s/widget_float", key, path)

		// Create or retrieving the widget.Float state
		floatState := c.State(floatStatePath, func() any {
			return &widget.Float{Value: internalValue}
		})
		wFloat := floatState.Get().(*widget.Float)

		// Synchronize external value to internal state if not dragging
		// This prevents jitter during drag if the external update lags
		if !wFloat.Dragging() {
			wFloat.Value = internalValue
		}

		constructorArgs := sliderConstructorArgs{
			Value:                 value,
			OnValueChange:         onValueChange,
			OnValueChangeFinished: opts.OnValueChangeFinished,
			InternalFloat:         wFloat,
			ValueRange:            opts.ValueRange,
			Steps:                 opts.Steps,
			Colors:                opts.Colors,
			Enabled:               opts.Enabled,
		}

		c.StartBlock(SliderNodeID)
		c.Modifier(func(modifier ui.Modifier) ui.Modifier {
			return modifier.Then(opts.Modifier)
		})
		c.SetWidgetConstructor(sliderWidgetConstructor(constructorArgs))
		return c.EndBlock()
	}
}

func resolveSliderColors(c Composer, colors SliderColors) SliderColors {
	theme := material3.Theme(c)
	return SliderColors{
		ThumbColor:            colors.ThumbColor.TakeOrElse(theme.ColorScheme().Primary),                                            //, selector.PrimaryRoles.Primary),
		ActiveTrackColor:      colors.ActiveTrackColor.TakeOrElse(theme.ColorScheme().Primary),                                      //, selector.PrimaryRoles.Primary),
		ActiveTickColor:       colors.ActiveTickColor.TakeOrElse(graphics.SetOpacity(theme.ColorScheme().OnPrimary, 0.38)),          //, selector.PrimaryRoles.OnPrimary.SetOpacity(0.38)),
		InactiveTrackColor:    colors.InactiveTrackColor.TakeOrElse(theme.ColorScheme().SurfaceContainerHighest),                    //, selector.SurfaceRoles.ContainerHighest),
		InactiveTickColor:     colors.InactiveTickColor.TakeOrElse(graphics.SetOpacity(theme.ColorScheme().OnSurfaceVariant, 0.38)), //, selector.SurfaceRoles.OnVariant.SetOpacity(0.38)),
		DisabledThumbColor:    colors.DisabledThumbColor.TakeOrElse(graphics.SetOpacity(theme.ColorScheme().OnSurface, 0.38)),       //, selector.SurfaceRoles.OnSurface.SetOpacity(0.38)),
		DisabledActiveTrack:   colors.DisabledActiveTrack.TakeOrElse(graphics.SetOpacity(theme.ColorScheme().OnSurface, 0.38)),      //, selector.SurfaceRoles.OnSurface.SetOpacity(0.38)),
		DisabledActiveTick:    colors.DisabledActiveTick.TakeOrElse(graphics.SetOpacity(theme.ColorScheme().OnSurface, 0.38)),       //, selector.SurfaceRoles.OnSurface.SetOpacity(0.38)),
		DisabledInactiveTrack: colors.DisabledInactiveTrack.TakeOrElse(graphics.SetOpacity(theme.ColorScheme().OnSurface, 0.12)),    //, selector.SurfaceRoles.OnSurface.SetOpacity(0.12)),
		DisabledInactiveTick:  colors.DisabledInactiveTick.TakeOrElse(graphics.SetOpacity(theme.ColorScheme().OnSurface, 0.12)),     //, selector.SurfaceRoles.OnSurface.SetOpacity(0.12)),
	}
}

type sliderConstructorArgs struct {
	Value                 float32
	OnValueChange         func(float32)
	OnValueChangeFinished func()
	InternalFloat         *widget.Float
	ValueRange            struct{ Min, Max float32 }
	Steps                 int
	Colors                SliderColors
	Enabled               bool
}

func sliderWidgetConstructor(args sliderConstructorArgs) layoutnode.LayoutNodeWidgetConstructor {
	return layoutnode.NewLayoutNodeWidgetConstructor(func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {
		return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
			wFloat := args.InternalFloat

			// Check for updates from the widget
			if wFloat.Dragging() {
				// Map 0-1 back to value range
				rangeDiff := args.ValueRange.Max - args.ValueRange.Min
				newValue := args.ValueRange.Min + wFloat.Value*rangeDiff

				// Handle steps if configured
				if args.Steps > 0 {
					stepSize := rangeDiff / float32(args.Steps+1)
					steps := (newValue - args.ValueRange.Min) / stepSize
					roundedSteps := float32(int(steps + 0.5))
					newValue = args.ValueRange.Min + roundedSteps*stepSize

					// Update internal float to snap visually
					wFloat.Value = (newValue - args.ValueRange.Min) / rangeDiff
				}

				if newValue != args.Value {
					args.OnValueChange(newValue)
				}
			}
			// Layout Constants
			trackHeight := gtx.ToGio().Dp(gioUnit.Dp(TrackHeight))
			tSize := ThumbSize
			if wFloat.Dragging() {
				tSize = ActiveThumbSize
			}
			thumbSize := gtx.ToGio().Dp(gioUnit.Dp(tSize))

			// Main Axis: Max width, fixed height (thumb size or min touch size)
			// Cross Axis: Thumb size

			// Minimum touch target size (48dp usually)
			minTouchSize := gtx.ToGio().Dp(gioUnit.Dp(48))
			height := max(thumbSize, minTouchSize)

			size := image.Pt(gtx.ToGio().Constraints.Max.X, height)

			// Center the component vertically
			centerOffset := op.Offset(image.Pt(0, (height-thumbSize)/2)).Push(gtx.ToGio().Ops)

			// Draw Track
			trackY := (thumbSize - trackHeight) / 2
			trackWidth := size.X

			// Calculate Filled Fraction
			fraction := wFloat.Value
			if args.Steps > 0 {
				// Snap fraction for display if stepping
				// (Already handled in logic above for dragging, but ensure consistency)
				rangeDiff := args.ValueRange.Max - args.ValueRange.Min
				if rangeDiff > 0 {
					fraction = (args.Value - args.ValueRange.Min) / rangeDiff
				}
			}

			activeTrackWidth := int(float32(trackWidth) * fraction)

			// Inactive Track (Full width first)
			inactiveTrackRect := image.Rect(0, trackY, trackWidth, trackY+trackHeight)
			roundedCorners := trackHeight / 2

			trackColor := graphics.ColorToNRGBA(args.Colors.Track(args.Enabled, false))

			// Using Clip/Paint
			inactiveTrackClip := clip.RRect{
				Rect: inactiveTrackRect,
				SE:   roundedCorners, SW: roundedCorners, NW: roundedCorners, NE: roundedCorners,
			}.Push(gtx.ToGio().Ops)
			paint.ColorOp{Color: trackColor}.Add(gtx.ToGio().Ops)
			paint.PaintOp{}.Add(gtx.ToGio().Ops)
			inactiveTrackClip.Pop()

			// Active Track (Overlay)
			if activeTrackWidth > 0 {
				activeTrackRect := image.Rect(0, trackY, activeTrackWidth, trackY+trackHeight)
				activeColor := graphics.ColorToNRGBA(args.Colors.Track(args.Enabled, true))

				activeTrackClip := clip.RRect{
					Rect: activeTrackRect,
					SE:   roundedCorners, SW: roundedCorners, NW: roundedCorners, NE: roundedCorners,
				}.Push(gtx.ToGio().Ops)
				paint.ColorOp{Color: activeColor}.Add(gtx.ToGio().Ops)
				paint.PaintOp{}.Add(gtx.ToGio().Ops)
				activeTrackClip.Pop()
			}

			// Ticks
			if args.Steps > 0 {
				tickColor := graphics.ColorToNRGBA(args.Colors.Tick(args.Enabled, true))

				inactiveTickColor := graphics.ColorToNRGBA(args.Colors.Tick(args.Enabled, false))

				tickSizePx := gtx.ToGio().Dp(gioUnit.Dp(TickSize))
				stepSizePx := float32(trackWidth) / float32(args.Steps+1)

				for i := 0; i <= args.Steps+1; i++ {
					x := int(float32(i) * stepSizePx)
					// Skip start and end ticks if desired, usually M3 keeps them but thumb covers

					// Determine color based on active/inactive
					tColor := activeTrackWidth >= x
					c := inactiveTickColor
					if tColor {
						c = tickColor
					}

					tickRect := image.Rect(x-tickSizePx/2, trackY+(trackHeight-tickSizePx)/2, x+tickSizePx/2, trackY+(trackHeight+tickSizePx)/2+tickSizePx)
					// Actually simplified to small circles
					tickCircle := clip.Ellipse{
						Min: tickRect.Min,
						Max: tickRect.Max,
					}
					tickClip := tickCircle.Push(gtx.ToGio().Ops)
					paint.ColorOp{Color: c}.Add(gtx.ToGio().Ops)
					paint.PaintOp{}.Add(gtx.ToGio().Ops)
					tickClip.Pop()
				}
			}

			// Draw Thumb
			thumbX := int(float32(trackWidth) * fraction)
			thumbRect := image.Rect(thumbX-thumbSize/2, 0, thumbX+thumbSize/2, thumbSize)

			thumbColor := graphics.ColorToNRGBA(args.Colors.Thumb(args.Enabled))

			thumbClip := clip.Ellipse{
				Min: thumbRect.Min,
				Max: thumbRect.Max,
			}.Push(gtx.ToGio().Ops)
			paint.ColorOp{Color: thumbColor}.Add(gtx.ToGio().Ops)
			paint.PaintOp{}.Add(gtx.ToGio().Ops)
			thumbClip.Pop()

			centerOffset.Pop()

			// Input Handling
			// The original widget.Float expects layout to pass dimensions and handle input
			// We wrap it securely
			// Ensure we occupy the full space for hit testing
			// Logic: widget.Float.Layout processes events. We must call it.
			// It draws nothing by itself usually in generic context, or we just consume its logic.
			// Re-reading widget.Float: it uses gtx.Constraints.Min for size.

			// We use a separate macro to isolate the input area if needed, or just overlay
			// widget.Float.Layout returns dimensions.

			// Set constraints for input to match our laid out track area (with touch target padding)
			g := *gtx.ToGio()
			g.Constraints.Min = size

			// Use layout.Stack to ensure input covers everything?
			// Actually widget.Float assumes it fills the area provided.
			// We just call it to register input.

			if args.Enabled {
				// Offset input to top-left of context (0,0)
				// widget.Float.Layout(gtx, axis, margin)
				// We need to ensure it processes input over the full 'size'
				// The widget.Float implementation uses axis.Convert(size) for length calculation.
				wasDragging := wFloat.Dragging()
				wFloat.Layout(g, layout.Horizontal, gioUnit.Dp(0))
				if wasDragging && !wFloat.Dragging() && args.OnValueChangeFinished != nil {
					args.OnValueChangeFinished()
				}
			}

			return layoutnode.LayoutDimensions{Size: size}
		}
	})
}

// max is a helper for integer max
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
