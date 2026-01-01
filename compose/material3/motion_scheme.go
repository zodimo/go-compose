package material3

import (
	"math"
	"time"
)

// The default duration used in [VectorizedAnimationSpec]s and [AnimationSpec].
var DefaultDurationMillis time.Duration = 300 * time.Millisecond

// The value that is used when the animation time is not yet set.
var UnspecifiedTime time.Duration = math.MinInt64

var MotionSchemeUnspecified = &MotionScheme{
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

type MotionScheme struct {
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

var DefaultMotionScheme = &MotionScheme{
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
