package pointer

import (
	"github.com/zodimo/go-compose/compose/ui/geometry"
	"github.com/zodimo/go-compose/compose/ui/platform"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

/**
 * Receiver scope for awaiting pointer events in a call to
 * [PointerInputScope.awaitPointerEventScope].
 *
 * This is a restricted suspension scope. Code in this scope is always called un-dispatched and may
 * only suspend for calls to [awaitPointerEvent]. These functions resume synchronously and the
 * caller may mutate the result **before** the next await call to affect the next stage of the input
 * processing pipeline.
 */

type AwaitPointerEventScope interface {
	unit.DensityScope
	/**
	 * The measured size of the pointer input region. Input events will be reported with a
	 * coordinate space of (0, 0) to (size.width, size,height) as the input region, with (0, 0)
	 * indicating the upper left corner.
	 */
	Size() unit.IntSize

	/*
	 * The additional space applied to each side of the layout area. This can be
	 * non-[zero][Size.Zero] when `minimumTouchTargetSize` is set in [pointerInput].
	 */
	ExtendedTouchPadding() geometry.Size

	/** The [PointerEvent] from the most recent touch event. */
	CurrentEvent() PointerEvent

	/** The [ViewConfiguration] used to tune gesture detectors. */
	ViewConfiguration() platform.ViewConfiguration

	/**
	 * Suspend until a [PointerEvent] is reported to the specified input [pass]. [pass] defaults to
	 * [PointerEventPass.Main].
	 *
	 * [awaitPointerEvent] resumes **synchronously** in the restricted suspension scope. This means
	 * that callers can react immediately to input after [awaitPointerEvent] returns and affect both
	 * the current frame and the next handler or phase of the input processing pipeline. Callers
	 * should mutate the returned [PointerEvent] before awaiting another event to consume aspects of
	 * the event before the next stage of input processing runs.
	 */
	// AwaitPointerEvent(pass: PointerEventPass = PointerEventPass.Main): PointerEvent

	/**
	 * Runs [block] and returns the result of [block] or `null` if [timeMillis] has passed before
	 * [timeMillis].
	 */
	// suspend fun <T> withTimeoutOrNull(
	//     timeMillis: Long,
	//     block: suspend AwaitPointerEventScope.() -> T,
	// ): T? = block()

	/**
	 * Runs [block] and returns its results. An [PointerEventTimeoutCancellationException] is thrown
	 * if [timeMillis] has passed before [block] completes.
	 */
	// suspend fun <T> withTimeout(
	//     timeMillis: Long,
	//     block: suspend AwaitPointerEventScope.() -> T,
	// ): T = block()
}
