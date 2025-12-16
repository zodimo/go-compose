# State Management Guide

This guide covers best practices for managing state in go-compose applications.

## Core Concepts

### MutableValue

`MutableValue` is the primary state container in go-compose:

```go
type MutableValue interface {
    Get() any
    Set(value any)
}
```

Create state in a composable using `c.State()`:

```go
func UI(c api.Composer) api.Composer {
    counterValue := c.State("counter", func() any { return 0 })
    count := counterValue.Get().(int)
    
    // Use count in your UI...
    
    button.Filled(func() {
        counterValue.Set(count + 1)
    }, "Increment")(c)
    
    return c
}
```

## Immutable State Pattern

For complex state objects, use the **immutable pattern** where state mutation methods return new instances:

```go
type AppState struct {
    Items  []Item
    Filter Filter
}

// ✅ Good: Returns new state (immutable)
func (s *AppState) SetFilter(f Filter) *AppState {
    return &AppState{
        Items:  s.Items,
        Filter: f,
    }
}

// ❌ Bad: Mutates in place (mutable)
func (s *AppState) SetFilter(f Filter) {
    s.Filter = f  // This won't trigger recomposition!
}
```

### Usage Pattern

```go
func onFilterChange(filter Filter) {
    // Get current state, create new state, set it
    newState := GetAppState(stateValue).SetFilter(filter)
    stateValue.Set(newState)
}
```

## Component State Patterns

### Passing State to Sub-Components

**Option 1: Pass MutableValue** (recommended for components that update state)

```go
func TodoFooter(stateValue state.MutableValue) api.Composable {
    return func(c api.Composer) api.Composer {
        state := stateValue.Get().(*TodoState)
        
        button.Filled(func() {
            // Get fresh state at click time
            newState := stateValue.Get().(*TodoState).ClearCompleted()
            stateValue.Set(newState)
        }, "Clear")(c)
        
        return c
    }
}
```

**Option 2: Pass callbacks** (for simpler cases)

```go
func TodoItem(
    todo Todo,
    onToggle func(),
    onDelete func(),
) api.Composable
```

### Type-Safe State Helper

Use a typed getter to avoid repetitive type assertions:

```go
func GetAppState(mv api.MutableValue) *AppState {
    return mv.Get().(*AppState)
}
```

## Conditional Component Rendering

### The Caching Problem

go-compose caches widget instances by their composition path. This causes issues when a component's **type changes** based on state:

```go
// ❌ Problem: Button type changes but cached instance is reused
if selected {
    button.Filled(...)(c)  // Uses cached button from first render
} else {
    button.Text(...)(c)    // Never gets its own instance!
}
```

### Solution: Use `c.If()`

`c.If()` generates **different keys** for the true/false branches:

```go
// ✅ Solution: Each branch gets its own cached instance
c.If(
    selected,
    button.Filled(onClick, "Label"),  // Key: true branch
    button.Text(onClick, "Label"),    // Key: false branch
)(c)
```

### When to Use `c.If()`

Use `c.If()` when:
- Switching between different component types (Filled vs Text button)
- Conditionally rendering components where identity matters

Use regular `if` when:
- Just hiding/showing a single component (use `c.When()`)
- The component type doesn't change

## Best Practices

1. **Always read fresh state in callbacks**
   ```go
   onClick := func() {
       // ✅ Get state at click time, not composition time
       current := GetState(stateValue)
       stateValue.Set(current.DoSomething())
   }
   ```

2. **Keep state immutable** - All mutation methods should return new instances

3. **Use `c.If()` for type-switching components** - Prevents cached widget reuse bugs

4. **Prefer passing MutableValue** over multiple callbacks for complex components

## Example: TodoMVC Pattern

```go
// State model with immutable methods
type TodoState struct {
    Todos  []Todo
    Filter Filter
}

func (s *TodoState) AddTodo(text string) *TodoState {
    return &TodoState{
        Todos:  append(s.Todos, Todo{Text: text}),
        Filter: s.Filter,
    }
}

// Component receiving MutableValue
func TodoFooter(stateValue state.MutableValue) api.Composable {
    return func(c api.Composer) api.Composer {
        state := GetTodoState(stateValue)
        
        // Use c.If for buttons that change type based on selection
        c.If(
            state.Filter == FilterAll,
            button.Filled(func() {
                newState := GetTodoState(stateValue).SetFilter(FilterAll)
                stateValue.Set(newState)
            }, "All"),
            button.Text(func() {
                newState := GetTodoState(stateValue).SetFilter(FilterAll)
                stateValue.Set(newState)
            }, "All"),
        )(c)
        
        return c
    }
}
```
