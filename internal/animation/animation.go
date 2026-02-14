package animation

import (
	"time"

	"gioui.org/layout"
	"gioui.org/op"
)

// Animation holds state for an Animation between two states that
// is not invertible.
type Animation struct {
	duration  time.Duration
	startTime time.Time
}

func NewAnimation(duration time.Duration) Animation {
	return Animation{
		duration: duration,
	}
}

// Progress returns the current progress through the animation
// as a value in the range [0,1]
func (n *Animation) Progress(gtx layout.Context) float32 {
	if n.duration == time.Duration(0) {
		return 0
	}
	progressDur := gtx.Now.Sub(n.startTime)
	if progressDur > n.duration {
		return 1
	}
	gtx.Execute(op.InvalidateCmd{})
	progress := float32(progressDur.Milliseconds()) / float32(n.duration.Milliseconds())
	return progress
}

func (n *Animation) Start(now time.Time) {
	n.startTime = now
}

func (n *Animation) SetDuration(d time.Duration) {
	n.duration = d
}

func (n *Animation) Animating(gtx layout.Context) bool {
	if n.duration == 0 {
		return false
	}
	if gtx.Now.After(n.startTime.Add(n.duration)) {
		return false
	}
	return true
}

func (n *Animation) StartTime() time.Time {
	return n.startTime
}

func (n *Animation) Duration() time.Duration {
	return n.duration
}
