# Canvas in go-compose

# components
- compose.foundation.canvas
- compose.ui.graphics.canvas



# Kotlin Constructors
```kotlin


package androidx.compose.foundation

import androidx.compose.foundation.layout.ColumnScope
import androidx.compose.foundation.layout.Spacer
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.drawBehind
import androidx.compose.ui.graphics.drawscope.DrawScope
import androidx.compose.ui.semantics.contentDescription
import androidx.compose.ui.semantics.semantics

/**
 * Component that allow you to specify an area on the screen and perform canvas drawing on this
 * area. You MUST specify size with modifier, whether with exact sizes via [Modifier.size] modifier,
 * or relative to parent, via [Modifier.fillMaxSize], [ColumnScope.weight], etc. If parent wraps
 * this child, only exact sizes must be specified.
 *
 * @sample androidx.compose.foundation.samples.CanvasSample
 * @param modifier mandatory modifier to specify size strategy for this composable
 * @param onDraw lambda that will be called to perform drawing. Note that this lambda will be called
 *   during draw stage, you have no access to composition scope, meaning that [Composable] function
 *   invocation inside it will result to runtime exception
 */
@Composable
fun Canvas(modifier: Modifier, onDraw: DrawScope.() -> Unit) = Spacer(modifier.drawBehind(onDraw))

/**
 * Component that allow you to specify an area on the screen and perform canvas drawing on this
 * area. You MUST specify size with modifier, whether with exact sizes via [Modifier.size] modifier,
 * or relative to parent, via [Modifier.fillMaxSize], [ColumnScope.weight], etc. If parent wraps
 * this child, only exact sizes must be specified.
 *
 * @sample androidx.compose.foundation.samples.CanvasPieChartSample
 * @param modifier mandatory modifier to specify size strategy for this composable
 * @param contentDescription text used by accessibility services to describe what this canvas
 *   represents. This should be provided unless the canvas is used for decorative purposes or as
 *   part of a larger entity already described in some other way. This text should be localized,
 *   such as by using [androidx.compose.ui.res.stringResource]
 * @param onDraw lambda that will be called to perform drawing. Note that this lambda will be called
 *   during draw stage, you have no access to composition scope, meaning that [Composable] function
 *   invocation inside it will result to runtime exception
 */
@Composable
fun Canvas(modifier: Modifier, contentDescription: String, onDraw: DrawScope.() -> Unit) =
    Spacer(modifier.drawBehind(onDraw).semantics { this.contentDescription = contentDescription })


```

# ui.graphics.canvas

compose/ui/ui-graphics/src/commonMain/kotlin/androidx/compose/ui/graphics/Canvas.kt
```kotlin
package androidx.compose.ui.graphics

import androidx.compose.ui.geometry.Offset
import androidx.compose.ui.geometry.Rect
import androidx.compose.ui.graphics.internal.JvmDefaultWithCompatibility
import androidx.compose.ui.unit.IntOffset
import androidx.compose.ui.unit.IntSize

/** Create a new Canvas instance that targets its drawing commands to the provided [ImageBitmap] */
fun Canvas(image: ImageBitmap): Canvas = ActualCanvas(image)

internal expect fun ActualCanvas(image: ImageBitmap): Canvas

expect class NativeCanvas


```

compose/ui/ui-graphics/src/skikoMain/kotlin/androidx/compose/ui/graphics/SkiaBackedCanvas.skiko.kt

```kotlin
package androidx.compose.ui.graphics

import androidx.compose.runtime.InternalComposeApi
import org.jetbrains.skia.ClipMode as SkClipMode
import org.jetbrains.skia.RRect as SkRRect
import org.jetbrains.skia.Rect as SkRect
import androidx.compose.ui.geometry.Offset
import androidx.compose.ui.geometry.Rect
import androidx.compose.ui.geometry.Size
import androidx.compose.ui.unit.IntOffset
import androidx.compose.ui.unit.IntSize
import androidx.compose.ui.util.fastForEach
import org.jetbrains.skia.CubicResampler
import org.jetbrains.skia.FilterMipmap
import org.jetbrains.skia.FilterMode
import org.jetbrains.skia.Image
import org.jetbrains.skia.Matrix44
import org.jetbrains.skia.MipmapMode
import org.jetbrains.skia.SamplingMode
import org.jetbrains.skia.impl.use

actual typealias NativeCanvas = org.jetbrains.skia.Canvas

internal actual fun ActualCanvas(image: ImageBitmap): Canvas {
    val skiaBitmap = image.asSkiaBitmap()
    require(!skiaBitmap.isImmutable) {
        "Cannot draw on immutable ImageBitmap"
    }
    return SkiaBackedCanvas(org.jetbrains.skia.Canvas(skiaBitmap))
}

/**
 * Convert the [org.jetbrains.skia.Canvas] instance into a Compose-compatible Canvas
 */
fun org.jetbrains.skia.Canvas.asComposeCanvas(): Canvas = SkiaBackedCanvas(this)

actual val Canvas.nativeCanvas: NativeCanvas get() = (this as SkiaBackedCanvas).skia

```

# DrawScope

compose-compact/ui/graphics/drawscope/DrawScope.kt
```kotlin
import androidx.compose.ui.graphics.drawscope

```



<!-- actual typealias NativeCanvas = org.jetbrains.skia.Canvas -->



# modifier.drawBehind
compose-multiplatform-core/compose/ui/ui/src/commonMain/kotlin/androidx/compose/ui/draw/DrawModifier.kt
```kotlin
/** Draw into a [Canvas] behind the modified content. */
fun Modifier.drawBehind(onDraw: DrawScope.() -> Unit) = this then DrawBehindElement(onDraw)

private class DrawBehindElement(val onDraw: DrawScope.() -> Unit) :
    ModifierNodeElement<DrawBackgroundModifier>() {
    override fun create() = DrawBackgroundModifier(onDraw)

    override fun update(node: DrawBackgroundModifier) {
        node.onDraw = onDraw
    }

    override fun InspectorInfo.inspectableProperties() {
        name = "drawBehind"
        properties["onDraw"] = onDraw
    }

    override fun equals(other: Any?): Boolean {
        if (this === other) return true
        if (other !is DrawBehindElement) return false

        if (onDraw !== other.onDraw) return false

        return true
    }

    override fun hashCode(): Int {
        return onDraw.hashCode()
    }
}

internal class DrawBackgroundModifier(var onDraw: DrawScope.() -> Unit) :  Modifier.Node(), DrawModifierNode {

    override fun ContentDrawScope.draw() {
        onDraw()
        drawContent()
    }
}

```