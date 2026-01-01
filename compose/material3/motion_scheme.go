package material3

import (
	"math"
	"time"
)

// The default duration used in [VectorizedAnimationSpec]s and [AnimationSpec].
var DefaultDurationMillis time.Duration = 300 * time.Millisecond

// The value that is used when the animation time is not yet set.
var UnspecifiedTime time.Duration = math.MinInt64

// Easing is a way to adjust an animation's fraction. Easing allows transitioning
// elements to speed up and slow down, rather than moving at a constant rate.
//
// Fraction is a value between 0 and 1.0 indicating the current point in the
// animation where 0 represents the start and 1.0 represents the end.
type Easing interface {
	Transform(fraction float32) float32
}

// CubicBezierEasing is a cubic polynomial easing implementing third-order Bézier curves.
//
// This is equivalent to Android's PathInterpolator when a single cubic Bézier curve
// is specified.
//
// Parameters:
//   - A: The x coordinate of the first control point
//   - B: The y coordinate of the first control point
//   - C: The x coordinate of the second control point
//   - D: The y coordinate of the second control point
type CubicBezierEasing struct {
	A, B, C, D float32
}

// NewCubicBezierEasing creates a new CubicBezierEasing with the given control points.
func NewCubicBezierEasing(a, b, c, d float32) *CubicBezierEasing {
	return &CubicBezierEasing{A: a, B: b, C: c, D: d}
}

// Transform transforms the specified fraction (0..1) by this cubic Bézier curve.
func (e *CubicBezierEasing) Transform(fraction float32) float32 {
	if fraction <= 0 {
		return 0
	}
	if fraction >= 1 {
		return 1
	}

	// Find t for the given x (fraction) using Newton-Raphson method
	t := fraction
	for i := 0; i < 8; i++ {
		x := e.evaluateX(t) - fraction
		if abs32(x) < 1e-6 {
			break
		}
		dx := e.evaluateDX(t)
		if abs32(dx) < 1e-6 {
			break
		}
		t -= x / dx
	}

	// Clamp t to [0, 1]
	if t < 0 {
		t = 0
	} else if t > 1 {
		t = 1
	}

	return e.evaluateY(t)
}

// evaluateX evaluates the x coordinate of the Bézier curve at parameter t
func (e *CubicBezierEasing) evaluateX(t float32) float32 {
	// B(t) = (1-t)³*0 + 3*(1-t)²*t*a + 3*(1-t)*t²*c + t³*1
	oneMinusT := 1 - t
	return 3*oneMinusT*oneMinusT*t*e.A + 3*oneMinusT*t*t*e.C + t*t*t
}

// evaluateDX evaluates the derivative of x with respect to t
func (e *CubicBezierEasing) evaluateDX(t float32) float32 {
	// dB(t)/dt = 3*(1-t)²*a + 6*(1-t)*t*(c-a) + 3*t²*(1-c)
	oneMinusT := 1 - t
	return 3*oneMinusT*oneMinusT*e.A + 6*oneMinusT*t*(e.C-e.A) + 3*t*t*(1-e.C)
}

// evaluateY evaluates the y coordinate of the Bézier curve at parameter t
func (e *CubicBezierEasing) evaluateY(t float32) float32 {
	// B(t) = (1-t)³*0 + 3*(1-t)²*t*b + 3*(1-t)*t²*d + t³*1
	oneMinusT := 1 - t
	return 3*oneMinusT*oneMinusT*t*e.B + 3*oneMinusT*t*t*e.D + t*t*t
}

func abs32(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}

// ---- Material 3 Motion Easing Tokens ----

// EasingEmphasized is the default emphasized easing curve.
// Elements that begin and end at rest use this standard easing.
var EasingEmphasized = NewCubicBezierEasing(0.2, 0.0, 0.0, 1.0)

// EasingEmphasizedAccelerate is for elements that are exiting the screen.
var EasingEmphasizedAccelerate = NewCubicBezierEasing(0.3, 0.0, 0.8, 0.15)

// EasingEmphasizedDecelerate is for elements that are entering the screen.
var EasingEmphasizedDecelerate = NewCubicBezierEasing(0.05, 0.7, 0.1, 1.0)

// EasingLegacy is the legacy Material easing curve (FastOutSlowIn).
var EasingLegacy = NewCubicBezierEasing(0.4, 0.0, 0.2, 1.0)

// EasingLegacyAccelerate is the legacy accelerate curve (FastOutLinearIn).
var EasingLegacyAccelerate = NewCubicBezierEasing(0.4, 0.0, 1.0, 1.0)

// EasingLegacyDecelerate is the legacy decelerate curve (LinearOutSlowIn).
var EasingLegacyDecelerate = NewCubicBezierEasing(0.0, 0.0, 0.2, 1.0)

// EasingLinear is a linear easing with no acceleration or deceleration.
var EasingLinear = NewCubicBezierEasing(0.0, 0.0, 1.0, 1.0)

// EasingStandard is the standard Material 3 easing curve.
var EasingStandard = NewCubicBezierEasing(0.2, 0.0, 0.0, 1.0)

// EasingStandardAccelerate is the standard accelerate curve.
var EasingStandardAccelerate = NewCubicBezierEasing(0.3, 0.0, 1.0, 1.0)

// EasingStandardDecelerate is the standard decelerate curve.
var EasingStandardDecelerate = NewCubicBezierEasing(0.0, 0.0, 0.0, 1.0)

// ---- Common Easing Aliases (matching animation-core conventions) ----

// FastOutSlowInEasing is equivalent to EasingLegacy.
var FastOutSlowInEasing = EasingLegacy

// LinearOutSlowInEasing is equivalent to EasingLegacyDecelerate.
var LinearOutSlowInEasing = EasingLegacyDecelerate

// FastOutLinearInEasing is equivalent to EasingLegacyAccelerate.
var FastOutLinearInEasing = EasingLegacyAccelerate

// LinearEasing returns the fraction unmodified.
type linearEasing struct{}

func (linearEasing) Transform(fraction float32) float32 { return fraction }

var MotionTokensUnspecified = &MotionTokens{
	DurationShort1:     UnspecifiedTime,
	DurationShort2:     UnspecifiedTime,
	DurationShort3:     UnspecifiedTime,
	DurationShort4:     UnspecifiedTime,
	DurationMedium1:    UnspecifiedTime,
	DurationMedium2:    UnspecifiedTime,
	DurationMedium3:    UnspecifiedTime,
	DurationMedium4:    UnspecifiedTime,
	DurationLong1:      UnspecifiedTime,
	DurationLong2:      UnspecifiedTime,
	DurationLong3:      UnspecifiedTime,
	DurationLong4:      UnspecifiedTime,
	DurationExtraLong1: UnspecifiedTime,
	DurationExtraLong2: UnspecifiedTime,
	DurationExtraLong3: UnspecifiedTime,
	DurationExtraLong4: UnspecifiedTime,
}

type MotionTokens struct {
	DurationShort1 time.Duration
	DurationShort2 time.Duration
	DurationShort3 time.Duration
	DurationShort4 time.Duration

	DurationMedium1 time.Duration
	DurationMedium2 time.Duration
	DurationMedium3 time.Duration
	DurationMedium4 time.Duration

	DurationLong1 time.Duration
	DurationLong2 time.Duration
	DurationLong3 time.Duration
	DurationLong4 time.Duration

	DurationExtraLong1 time.Duration
	DurationExtraLong2 time.Duration
	DurationExtraLong3 time.Duration
	DurationExtraLong4 time.Duration
}

var DefaultMotionTokens = &MotionTokens{
	DurationShort1:     50 * time.Millisecond,
	DurationShort2:     100 * time.Millisecond,
	DurationShort3:     150 * time.Millisecond,
	DurationShort4:     200 * time.Millisecond,
	DurationMedium1:    250 * time.Millisecond,
	DurationMedium2:    300 * time.Millisecond,
	DurationMedium3:    350 * time.Millisecond,
	DurationMedium4:    400 * time.Millisecond,
	DurationLong1:      450 * time.Millisecond,
	DurationLong2:      500 * time.Millisecond,
	DurationLong3:      550 * time.Millisecond,
	DurationLong4:      600 * time.Millisecond,
	DurationExtraLong1: 700 * time.Millisecond,
	DurationExtraLong2: 800 * time.Millisecond,
	DurationExtraLong3: 900 * time.Millisecond,
	DurationExtraLong4: 1000 * time.Millisecond,
}

type MotionScheme struct{}

// DefaultMotionScheme is the default motion scheme with all standard easing curves.
var DefaultMotionScheme = &MotionScheme{}
