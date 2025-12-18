package progress

import (
	"fmt"
	"image"
	"math"
	"time"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"

	"git.sr.ht/~schnwalter/gio-mw/wdk"
	"git.sr.ht/~schnwalter/gio-mw/widget/indicator"
	"github.com/zodimo/go-compose/internal/layoutnode"
)

// LoadingIndicator displays a circular indeterminate progress indicator.
func LoadingIndicator(options ...IndicatorOption) Composable {
	return func(c Composer) Composer {
		opts := DefaultIndicatorOptions()
		for _, option := range options {
			if option == nil {
				continue
			}
			option(&opts)
		}

		// key := c.GenerateID() // unused
		path := c.GetPath()
		statePath := fmt.Sprintf("%v/loading_anim_state", path)

		// State for animation
		animState := c.State(statePath, func() any {
			return &loadingState{}
		}).Get().(*loadingState)

		c.StartBlock("LoadingIndicator")
		c.Modifier(func(modifier Modifier) Modifier {
			return modifier.Then(opts.Modifier)
		})

		c.SetWidgetConstructor(layoutnode.NewLayoutNodeWidgetConstructor(func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {
			return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
				return drawLoadingIndicator(gtx, animState, opts)
			}
		}))

		return c.EndBlock()
	}
}

type loadingState struct {
	startTime time.Time
}

func drawLoadingIndicator(gtx layout.Context, state *loadingState, opts IndicatorOptions) layout.Dimensions {
	// Ensure continuous animation
	gtx.Execute(op.InvalidateCmd{})

	if state.startTime.IsZero() {
		state.startTime = gtx.Now
	}

	elapsed := gtx.Now.Sub(state.startTime)

	// Default size matching gio-mw
	diameter := unit.Dp(48)

	// Get theme defaults using gio-mw's BuildTheme if possible, or fallback manually
	th := indicator.BuildTheme(gtx)

	col := th.EnabledActiveIndicatorColor.AsNRGBA()

	gtx.Constraints.Max.Y = gtx.Dp(diameter)
	actualDiameter := float32(gtx.Dp(diameter))

	// Animation parameters (Material Design specs-ish)
	// Cycle duration: ~1333ms is standard for Material
	const cycleDuration = 1333 * time.Millisecond

	cyclePos := float64(elapsed%cycleDuration) / float64(cycleDuration)

	// 1. Rotation of the whole indicator
	// Rotate linearly
	rotation := float32(elapsed.Seconds() * 1.5 * 2 * math.Pi)

	// Oscillate sweep angle (arc length)
	// range from ~10 degrees to ~300 degrees
	minSweep := float32(10)
	maxSweep := float32(300)

	// Sine wave based sweep
	sweepFactor := math.Sin(cyclePos * 2 * math.Pi) // -1 to 1
	sweepFactor = (sweepFactor + 1) / 2             // 0 to 1
	sweepAngle := minSweep + float32(sweepFactor)*(maxSweep-minSweep)

	defer op.Affine(f32.Affine2D{}.Rotate(f32.Pt(actualDiameter/2, actualDiameter/2), rotation)).Push(gtx.Ops).Pop()

	shape := wdk.Arc{
		Center:     f32.Pt(actualDiameter/2, actualDiameter/2),
		Diameter:   actualDiameter,
		StartAngle: 0, // Rotation handles the movement
		SweepAngle: sweepAngle,
	}

	stroke := clip.Stroke{
		Path:  shape.Path(gtx),
		Width: float32(gtx.Dp(th.EnabledActiveIndicatorThickness)),
	}.Op().Push(gtx.Ops)

	paint.Fill(gtx.Ops, col)
	stroke.Pop()

	return layout.Dimensions{
		Size: image.Point{
			X: gtx.Dp(diameter),
			Y: gtx.Dp(diameter),
		},
	}
}
