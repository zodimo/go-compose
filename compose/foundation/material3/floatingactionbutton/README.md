# FAB
```kotlin
@Composable
fun FloatingActionButton(
    onClick: () -> Unit,
    modifier: Modifier = Modifier,
    shape: Shape = FloatingActionButtonDefaults.shape,
    containerColor: Color = FloatingActionButtonDefaults.containerColor,
    contentColor: Color = contentColorFor(containerColor),
    elevation: FloatingActionButtonElevation = FloatingActionButtonDefaults.elevation(),
    interactionSource: MutableInteractionSource? = null,
    content: @Composable () -> Unit
): Unit
```


|Type|Original M3|M3 Expressive|
|---|---|---|
|FAB|Available|Available|
|Medium FAB|--|Available|
|Large FAB|Available|Available|
|Small FAB|Available|Deprecated. Use a larger size.|
