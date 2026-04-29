package pointer

import (
	"github.com/zodimo/go-compose/compose/ui/geometry"
	"github.com/zodimo/go-compose/compose/ui/platform"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

type DragGestureDetectorScope interface {
	DetectDragGestures(
		OnDragStart func(geometry.Offset),
		OnDragEnd func(),
		OnDragCancel func(),
		OnDrag func(change PointerInputChange, dragAmount geometry.Offset),
	)
}

type DragGestureDetectorScopeHandler struct {
	onDragStart  func(geometry.Offset)
	onDragEnd    func()
	onDragCancel func()
	onDrag       func(change PointerInputChange, dragAmount geometry.Offset)
}

var _ DragGestureDetectorScope = (*DragGestureDetectorScopeHandler)(nil)

func NewDragGestureDetectorScopeAndHandler() *DragGestureDetectorScopeHandler {
	return &DragGestureDetectorScopeHandler{}
}

func (s *DragGestureDetectorScopeHandler) DetectDragGestures(
	onDragStart func(geometry.Offset),
	onDragEnd func(),
	onDragCancel func(),
	onDrag func(change PointerInputChange, dragAmount geometry.Offset),
) {
	s.onDragStart = onDragStart
	s.onDragEnd = onDragEnd
	s.onDragCancel = onDragCancel
	s.onDrag = onDrag
}

func (s *DragGestureDetectorScopeHandler) DragStart(p geometry.Offset) {
	if s.onDragStart != nil {
		s.onDragStart(p)
	}
}

func (s *DragGestureDetectorScopeHandler) DragEnd() {
	if s.onDragEnd != nil {
		s.onDragEnd()
	}
}

func (s *DragGestureDetectorScopeHandler) DragCancel() {
	if s.onDragCancel != nil {
		s.onDragCancel()
	}
}

func (s *DragGestureDetectorScopeHandler) Drag(change PointerInputChange, dragAmount geometry.Offset) {
	if s.onDrag != nil {
		s.onDrag(change, dragAmount)
	}
}

/**
 * Receiver scope for [Modifier.pointerInput] that permits
 * [handling pointer input][awaitPointerEventScope].
 */
// Design note: this interface does _not_ implement CoroutineScope, even though doing so
// would more easily permit the use of launch {} inside Modifier.pointerInput {} blocks without
// requiring an additional coroutineScope {} layer of nesting. As it is encouraged to define
// gesture detectors as suspending extensions with a PointerInputScope receiver, also making this
// interface implement CoroutineScope would be an invitation to break structured concurrency in
// these extensions, leaving other launched coroutines running in the calling scope.
type PointerInputScope interface {
	unit.DensityScope
	/**
	 * The measured size of the pointer input region. Input events will be reported with a
	 * coordinate space of (0, 0) to (size.width, size,height) as the input region, with (0, 0)
	 * indicating the upper left corner.
	 */
	Size() unit.IntSize

	/**
	 * The additional space applied to each side of the layout area when the layout is smaller than
	 * [ViewConfiguration.minimumTouchTargetSize].
	 */
	ExtendedTouchPadding() geometry.Size

	/** The [ViewConfiguration] used to tune gesture detectors. */
	ViewConfiguration() platform.ViewConfiguration

	/**
	 * Suspend and install a pointer input [block] that can await input events and respond to them
	 * immediately. A call to [awaitPointerEventScope] will resume with [block]'s result after it
	 * completes.
	 *
	 * More than one [awaitPointerEventScope] can run concurrently in the same [PointerInputScope]
	 * by using [kotlinx.coroutines.launch]. [block]s are dispatched to in the order in which they
	 * were installed.
	 */
	// suspend fun <R> awaitPointerEventScope(block: suspend AwaitPointerEventScope.() -> R): R

	DragGestureDetectorScope
}

var _ PointerInputScope = (*pointerInputScope)(nil)

type pointerInputScope struct {
	unit.DensityScope
	size                 unit.IntSize
	extendedTouchPadding geometry.Size
	viewConfiguration    platform.ViewConfiguration

	DragGestureDetectorScope
}

func NewPointerInputScope(
	densityScope unit.DensityScope,
	size unit.IntSize,
	extendedTouchPadding geometry.Size,
	viewConfiguration platform.ViewConfiguration,
	dragGestureDetectorScope DragGestureDetectorScope,
) PointerInputScope {
	return &pointerInputScope{
		DensityScope:             densityScope,
		size:                     size,
		extendedTouchPadding:     extendedTouchPadding,
		viewConfiguration:        viewConfiguration,
		DragGestureDetectorScope: dragGestureDetectorScope,
	}
}

func (s *pointerInputScope) Size() unit.IntSize {
	return s.size
}

func (s *pointerInputScope) ExtendedTouchPadding() geometry.Size {
	return s.extendedTouchPadding
}

func (s *pointerInputScope) ViewConfiguration() platform.ViewConfiguration {
	return s.viewConfiguration
}
