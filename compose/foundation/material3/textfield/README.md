# Textfield

```kotlin
@Composable
fun TextField(
    value: TextFieldValue,
    onValueChange: (TextFieldValue) -> Unit,
    modifier: Modifier = Modifier,
    enabled: Boolean = true,
    readOnly: Boolean = false,
    textStyle: TextStyle = LocalTextStyle.current,
    label: (@Composable () -> Unit)? = null,
    placeholder: (@Composable () -> Unit)? = null,
    leadingIcon: (@Composable () -> Unit)? = null,
    trailingIcon: (@Composable () -> Unit)? = null,
    prefix: (@Composable () -> Unit)? = null,
    suffix: (@Composable () -> Unit)? = null,
    supportingText: (@Composable () -> Unit)? = null,
    isError: Boolean = false,
    visualTransformation: VisualTransformation = VisualTransformation.None,
    keyboardOptions: KeyboardOptions = KeyboardOptions.Default,
    keyboardActions: KeyboardActions = KeyboardActions.Default,
    singleLine: Boolean = false,
    maxLines: Int = if (singleLine) 1 else Int.MAX_VALUE,
    minLines: Int = 1,
    interactionSource: MutableInteractionSource? = null,
    shape: Shape = TextFieldDefaults.shape,
    colors: TextFieldColors = TextFieldDefaults.colors()
): Unit
```

# Primary

```
value: TextFieldValue,
onValueChange: (TextFieldValue) -> Unit,
```

# Optional

```
modifier: Modifier = Modifier,
enabled: Boolean = true,
readOnly: Boolean = false,
textStyle: TextStyle = LocalTextStyle.current,
label: (@Composable () -> Unit)? = null,
placeholder: (@Composable () -> Unit)? = null,
leadingIcon: (@Composable () -> Unit)? = null,
trailingIcon: (@Composable () -> Unit)? = null,
prefix: (@Composable () -> Unit)? = null,
suffix: (@Composable () -> Unit)? = null,
supportingText: (@Composable () -> Unit)? = null,
isError: Boolean = false,
visualTransformation: VisualTransformation = VisualTransformation.None,
keyboardOptions: KeyboardOptions = KeyboardOptions.Default,
keyboardActions: KeyboardActions = KeyboardActions.Default,
singleLine: Boolean = false,
maxLines: Int = if (singleLine) 1 else Int.MAX_VALUE,
minLines: Int = 1,
interactionSource: MutableInteractionSource? = null,
shape: Shape = TextFieldDefaults.shape,
colors: TextFieldColors = TextFieldDefaults.colors()
```