# M3 Button

- [Ref](https://developer.android.com/reference/kotlin/androidx/compose/material3/package-summary#Button(kotlin.Function0,androidx.compose.material3.ButtonShapes,androidx.compose.ui.Modifier,kotlin.Boolean,androidx.compose.material3.ButtonColors,androidx.compose.material3.ButtonElevation,androidx.compose.foundation.BorderStroke,androidx.compose.foundation.layout.PaddingValues,androidx.compose.foundation.interaction.MutableInteractionSource,kotlin.Function1))
```kotlin
@Composable
@ExperimentalMaterial3ExpressiveApi
fun Button(
    onClick: () -> Unit,
    shapes: ButtonShapes,
    modifier: Modifier = Modifier,
    enabled: Boolean = true,
    colors: ButtonColors = ButtonDefaults.buttonColors(),
    elevation: ButtonElevation? = ButtonDefaults.buttonElevation(),
    border: BorderStroke? = null,
    contentPadding: PaddingValues = ButtonDefaults.contentPaddingFor(ButtonDefaults.MinHeight),
    interactionSource: MutableInteractionSource? = null,
    content: @Composable RowScope.() -> Unit
): Unit
```


# ButtonColors
- [ref](https://developer.android.com/reference/kotlin/androidx/compose/material3/ButtonColors)

```kotlin
ButtonColors(
    containerColor: Color,
    contentColor: Color,
    disabledContainerColor: Color,
    disabledContentColor: Color
)
```

# ButtonElevation
- [ref](https://developer.android.com/reference/kotlin/androidx/compose/material3/ButtonElevation)


```kotlin
@Composable
fun buttonElevation(
    defaultElevation: Dp = FilledButtonTokens.ContainerElevation,
    pressedElevation: Dp = FilledButtonTokens.PressedContainerElevation,
    focusedElevation: Dp = FilledButtonTokens.FocusedContainerElevation,
    hoveredElevation: Dp = FilledButtonTokens.HoveredContainerElevation,
    disabledElevation: Dp = FilledButtonTokens.DisabledContainerElevation
): ButtonElevation
```