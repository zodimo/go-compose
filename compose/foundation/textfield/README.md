# Basic Text Field
```kotlin
/**
 * Basic text composable that provides an interactive box that accepts text input through software
 * or hardware keyboard, but provides no decorations like hint or placeholder.
 *
 * All the editing state of this composable is hoisted through [state]. Whenever the contents of
 * this composable change via user input or semantics, [TextFieldState.text] gets updated.
 * Similarly, all the programmatic updates made to [state] also reflect on this composable.
 *
 * If you want to add decorations to your text field, such as icon or similar, and increase the hit
 * target area, use the decorator.
 *
 * In order to filter (e.g. only allow digits, limit the number of characters), or change (e.g.
 * convert every character to uppercase) the input received from the user, use an
 * [InputTransformation].
 *
 * Limiting the height of the [BasicTextField] in terms of line count and choosing a scroll
 * direction can be achieved by using [TextFieldLineLimits].
 *
 * Scroll state of the composable is also hoisted to enable observation and manipulation of the
 * scroll behavior by the developer, e.g. bringing a searched keyword into view by scrolling to its
 * position without focusing, or changing selection.
 *
 * It's also possible to internally wrap around an existing TextFieldState and expose a more
 * lightweight state hoisting mechanism through a value that dictates the content of the TextField
 * and an onValueChange callback that communicates the changes to this value.
 *
 * @param state [TextFieldState] object that holds the internal editing state of [BasicTextField].
 * @param modifier optional [Modifier] for this text field.
 * @param enabled controls the enabled state of the [BasicTextField]. When `false`, the text field
 *   will be neither editable nor focusable, the input of the text field will not be selectable.
 * @param readOnly controls the editable state of the [BasicTextField]. When `true`, the text field
 *   can not be modified, however, a user can focus it and copy text from it. Read-only text fields
 *   are usually used to display pre-filled forms that user can not edit.
 * @param inputTransformation Optional [InputTransformation] that will be used to transform changes
 *   to the [TextFieldState] made by the user. The transformation will be applied to changes made by
 *   hardware and software keyboard events, pasting or dropping text, accessibility services, and
 *   tests. The transformation will _not_ be applied when changing the [state] programmatically, or
 *   when the transformation is changed. If the transformation is changed on an existing text field,
 *   it will be applied to the next user edit. the transformation will not immediately affect the
 *   current [state].
 * @param textStyle Typographic and graphic style configuration for text content that's displayed in
 *   the editor.
 * @param keyboardOptions Software keyboard options that contain configurations such as
 *   [KeyboardType] and [ImeAction].
 * @param onKeyboardAction Called when the user presses the action button in the input method editor
 *   (IME), or by pressing the enter key on a hardware keyboard if the [lineLimits] is configured as
 *   [TextFieldLineLimits.SingleLine]. By default this parameter is null, and would execute the
 *   default behavior for a received IME Action e.g., [ImeAction.Done] would close the keyboard,
 *   [ImeAction.Next] would switch the focus to the next focusable item on the screen.
 * @param lineLimits Whether the text field should be [SingleLine], scroll horizontally, and ignore
 *   newlines; or [MultiLine] and grow and scroll vertically. If [SingleLine] is passed, all newline
 *   characters ('\n') within the text will be replaced with regular whitespace (' '), ensuring that
 *   the contents of the text field are presented in a single line.
 * @param onTextLayout Callback that is executed when the text layout becomes queryable. The
 *   callback receives a function that returns a [TextLayoutResult] if the layout can be calculated,
 *   or null if it cannot. The function reads the layout result from a snapshot state object, and
 *   will invalidate its caller when the layout result changes. A [TextLayoutResult] object contains
 *   paragraph information, size of the text, baselines and other details. The callback can be used
 *   to add additional decoration or functionality to the text. For example, to draw a cursor or
 *   selection around the text. [Density] scope is the one that was used while creating the given
 *   text layout.
 * @param interactionSource the [MutableInteractionSource] representing the stream of [Interaction]s
 *   for this TextField. You can create and pass in your own remembered [MutableInteractionSource]
 *   if you want to observe [Interaction]s and customize the appearance / behavior of this TextField
 *   for different [Interaction]s.
 * @param cursorBrush [Brush] to paint cursor with. If [SolidColor] with [Color.Unspecified]
 *   provided, then no cursor will be drawn.
 * @param outputTransformation An [OutputTransformation] that transforms how the contents of the
 *   text field are presented.
 * @param decorator Allows to add decorations around text field, such as icon, placeholder, helper
 *   messages or similar, and automatically increase the hit target area of the text field.
 * @param scrollState Scroll state that manages either horizontal or vertical scroll of TextField.
 *   If [lineLimits] is [SingleLine], this text field is treated as single line with horizontal
 *   scroll behavior. In other cases the text field becomes vertically scrollable.
 * @sample androidx.compose.foundation.samples.BasicTextFieldDecoratorSample
 * @sample androidx.compose.foundation.samples.BasicTextFieldCustomInputTransformationSample
 * @sample androidx.compose.foundation.samples.BasicTextFieldWithValueOnValueChangeSample
 */
// This takes a composable lambda, but it is not primarily a container.
@Suppress("ComposableLambdaParameterPosition")
@Composable
fun BasicTextField(
    state: TextFieldState,
    modifier: Modifier = Modifier,
    enabled: Boolean = true,
    readOnly: Boolean = false,
    inputTransformation: InputTransformation? = null,
    textStyle: TextStyle = TextStyle.Default,
    keyboardOptions: KeyboardOptions = KeyboardOptions.Default,
    onKeyboardAction: KeyboardActionHandler? = null,
    lineLimits: TextFieldLineLimits = TextFieldLineLimits.Default,
    onTextLayout: (Density.(getResult: () -> TextLayoutResult?) -> Unit)? = null,
    interactionSource: MutableInteractionSource? = null,
    cursorBrush: Brush = BasicTextFieldDefaults.CursorBrush,
    outputTransformation: OutputTransformation? = null,
    decorator: TextFieldDecorator? = null,
    scrollState: ScrollState = rememberScrollState(),
    // Last parameter must not be a function unless it's intended to be commonly used as a trailing
    // lambda.
) {
    BasicTextField(
        state = state,
        modifier = modifier,
        enabled = enabled,
        readOnly = readOnly,
        inputTransformation = inputTransformation,
        textStyle = textStyle,
        keyboardOptions = keyboardOptions,
        onKeyboardAction = onKeyboardAction,
        lineLimits = lineLimits,
        onTextLayout = onTextLayout,
        interactionSource = interactionSource,
        cursorBrush = cursorBrush,
        codepointTransformation = null,
        outputTransformation = outputTransformation,
        decorator = decorator,
        scrollState = scrollState,
    )
}
```