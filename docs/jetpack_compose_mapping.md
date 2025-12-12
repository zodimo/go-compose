# Jetpack Compose to Go-Compose Mapping

This document provides a mapping between Jetpack Compose (Kotlin) concepts and their equivalents in `go-compose`. It is intended to help developers familiar with Jetpack Compose understand and use this library.

## Core Concepts

### Composable Functions

In Jetpack Compose, a composable function is annotated with `@Composable`. In `go-compose`, a composable is a function that returns a `Composable` type (which is `func(Composer) Composer`).

| Jetpack Compose (Kotlin) | go-compose (Go) |
| :--- | :--- |
| `@Composable fun MyWidget() { ... }` | `func MyWidget() Composable { return func(c Composer) Composer { ... } }` |

**Note**: In Go, you often return a closure that accepts the `Composer`.

### Modifiers

Modifiers are used to decorate or add behavior to UI elements. They are chained in both systems.

| Jetpack Compose | go-compose |
| :--- | :--- |
| `Modifier` | `Modifier` (interface) |
| `Modifier.padding(10.dp)` | (Depends on implementation, typically passed via options) |
| `Modifier.then(other)` | `modifier.Then(other)` |
| `Modifier` (empty) | `EmptyModifier` |

### State Management

State management is crucial for declarative UIs.

| Jetpack Compose | go-compose | Notes |
| :--- | :--- | :--- |
| `remember { ... }` | `c.Remember(key, func() any)` | Requires a unique key string. |
| `mutableStateOf(...)` | `c.State(key, func() any)` | Returns `MutableValue` with `Get()` and `Set()`. |

**Example (State):**

```go
// go-compose
count := c.State("my-counter", func() any { return 0 })
val := count.Get().(int)
count.Set(val + 1)
```

## Foundation & Layouts

Layouts structure your UI.

| Jetpack Compose | go-compose | Notes |
| :--- | :--- | :--- |
| `Column { ... }` | `layout.Column(content, ...)` | Content passed as `Composable`. |
| `Row { ... }` | `layout.Row(content, ...)` | |
| `Box { ... }` | `layout.Box(content, ...)` | |

**Usage Difference**:
In Kotlin, children are passed as a trailing lambda. In Go, the content `Composable` is often the first argument or passed explicitly.

```go
// go-compose
column.Column(
    compose.Sequence(
        button.Text(func() { ... }, "Click Me"),
        button.Outlined(func() { ... }, "Or Me"),
    ),
    column.WithSpacing(10),
)
```

## Material 3 Components

`go-compose` aims to implement Material 3 components.

| Jetpack Compose | go-compose | Notes |
| :--- | :--- | :--- |
| `Button(onClick = { ... }) { Text("Label") }` | `button.Filled(onClick, "Label")` | Content is often simplified to a label string in current impl. |
| `OutlinedButton(...)` | `button.Outlined(onClick, "Label")` | |
| `TextButton(...)` | `button.Text(onClick, "Label")` | |
| `Card(...)` | `card.Card(...)` | |
| `Checkbox(...)` | `checkbox.Checkbox(...)` | |

**Options Pattern**:
Instead of named parameters with default values (which Go lacks), `go-compose` often uses the Functional Options pattern.

```go
button.Filled(
    onClick,
    "Label",
    button.WithModifier(myModifier),
)
```

## Structure Mapping

The directory structure of `go-compose` closely mirrors Jetpack Compose.

*   `compose/runtime` -> `androidx.compose.runtime`
*   `compose/ui` -> `androidx.compose.ui`
*   `compose/foundation` -> `androidx.compose.foundation`
*   `compose/foundation/layout` -> `androidx.compose.foundation.layout`
*   `compose/foundation/material3` -> `androidx.compose.material3`

## Key Differences Summary

1.  **Explicit Composer**: The `Composer` is explicitly passed and returned in the underlying implementation of Composables in Go.
2.  **Keys**: `Remember` and `State` calls currently require an explicit key string, whereas Jetpack Compose compiler generates these automatically in many cases.
3.  **Options Pattern**: Used for optional parameters instead of Kotlin's named arguments.
4.  **Type Safety**: Go's type system is less expressive for generics than Kotlin's in some ways, so you might see `any` (interface{}) used in state, requiring casting.
