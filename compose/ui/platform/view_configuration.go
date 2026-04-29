package platform

import (
	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

type ViewConfiguration interface {
	/** The duration before a press turns into a long press. */
	LongPressTimeoutMillis() int
	/**
	 * The duration between the first tap's up event and the second tap's down event for an
	 * interaction to be considered a double-tap.
	 */
	DoubleTapTimeoutMillis() int

	/**
	 * The minimum duration between the first tap's up event and the second tap's down event for an
	 * interaction to be considered a double-tap.
	 */
	DoubleTapMinTimeMillis() int

	/** Distance in pixels a touch can wander before we think the user is scrolling. */
	TouchSlop() float32

	/** Distance in pixels a stylus touch can wander before we think the user is handwriting. */
	HandwritingSlop() float32

	/**
	 * The minimum touch target size. If layout has reduced the pointer input bounds below this, the
	 * touch target will be expanded evenly around the layout to ensure that it is at least this
	 * big.
	 */
	MinimumTouchTargetSize() unit.DpSize

	/**
	 * The maximum velocity a fling have at any given time. This value should be in pixels/second.
	 */
	MaximumFlingVelocity() float32

	/** Minimum velocity to initiate a fling, as measured in pixels per second */
	MinimumFlingVelocity() float32

	/**
	 * Margin in pixels around text line bounds where stylus handwriting gestures should be
	 * supported.
	 */
	HandwritingGestureLineMargin() float32
}

type viewConfiguration struct {
	longPressTimeoutMillis       int
	doubleTapTimeoutMillis       int
	doubleTapMinTimeMillis       int
	touchSlop                    float32
	handwritingSlop              float32
	minimumTouchTargetSize       unit.DpSize
	maximumFlingVelocity         float32
	minimumFlingVelocity         float32
	handwritingGestureLineMargin float32
}

func (v *viewConfiguration) LongPressTimeoutMillis() int {
	return v.longPressTimeoutMillis
}

func (v *viewConfiguration) DoubleTapTimeoutMillis() int {
	return v.doubleTapTimeoutMillis
}

func (v *viewConfiguration) DoubleTapMinTimeMillis() int {
	return v.doubleTapMinTimeMillis
}

func (v *viewConfiguration) TouchSlop() float32 {
	return v.touchSlop
}

func (v *viewConfiguration) HandwritingSlop() float32 {
	return v.handwritingSlop
}

func (v *viewConfiguration) MinimumTouchTargetSize() unit.DpSize {
	return v.minimumTouchTargetSize
}

func (v *viewConfiguration) MaximumFlingVelocity() float32 {
	return v.maximumFlingVelocity
}

func (v *viewConfiguration) MinimumFlingVelocity() float32 {
	return v.minimumFlingVelocity
}

func (v *viewConfiguration) HandwritingGestureLineMargin() float32 {
	return v.handwritingGestureLineMargin
}

func NewDefaultViewConfiguration() ViewConfiguration {
	return &viewConfiguration{
		longPressTimeoutMillis:       500,
		doubleTapTimeoutMillis:       100,
		doubleTapMinTimeMillis:       100,
		touchSlop:                    0.117,
		handwritingSlop:              0.117,
		minimumTouchTargetSize:       unit.DpSize{Width: 48, Height: 48},
		maximumFlingVelocity:         1500,
		minimumFlingVelocity:         100,
		handwritingGestureLineMargin: 4,
	}
}

var LocalViewConfiguration = compose.CompositionLocalOf[ViewConfiguration](func() ViewConfiguration {
	panic("CompositionLocal LocalViewConfiguration not present")
})
