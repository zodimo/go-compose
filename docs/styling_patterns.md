# GoCompose Styling Patterns Guide

This guide documents idiomatic patterns for styling components in GoCompose, following Jetpack Compose principles.

## Core Principle: Prefer Modifiers Over Wrappers

When styling a single element, **chain modifiers directly** rather than wrapping in a container component like `Surface`.

### ✅ Idiomatic Pattern: Modifier Chain

```go
// Circle-clipped image with border
fImage.Image(
    imageResource,
    fImage.WithContentScale(uilayout.ContentScaleCrop),
    fImage.WithAlignment(size.Center),
    fImage.WithModifier(
        size.Size(100, 100).
            Then(clip.Clip(shape.ShapeCircle)).
            Then(border.Border(4, colors.PrimaryRoles.Primary, shape.ShapeCircle)),
    ),
)
```

### ❌ Avoid: Unnecessary Surface Wrapper

```go
// Don't do this for simple styling
surface.Surface(
    fImage.Image(imageResource, ...),
    surface.WithShape(shape.ShapeCircle),
    surface.WithBorder(4, borderColor),
    surface.WithModifier(size.Size(108, 108)),
)
```

## When to Use Each Pattern

| Use Case | Pattern |
|----------|---------|
| Single element with clip/border | **Modifier chain** |
| Element with background color | **Modifier chain** with `background.Background()` |
| Container with multiple children + shared styling | **Surface** |
| Card with elevation, padding, and content | **Surface** |
| Semantic elevation (FAB, Dialog, etc.) | **Surface** |

## Available Styling Modifiers

| Modifier | Purpose | Example |
|----------|---------|---------|
| `clip.Clip(shape)` | Clip to shape | `clip.Clip(shape.ShapeCircle)` |
| `border.Border(width, color, shape)` | Draw border | `border.Border(2, color, shape.ShapeCircle)` |
| `background.Background(color)` | Fill background | `background.Background(colors.SurfaceRoles.Surface)` |
| `shadow.Simple(elevation, shape)` | Add shadow | `shadow.Simple(4, shape.RoundedCornerShape{Radius: 8})` |
| `padding.All(dp)` | Add padding | `padding.All(16)` |

## Color Pattern: Use ColorSelector

Always use the `ColorSelector()` pattern for theme-reactive colors:

```go
colors := theme.ColorHelper.ColorSelector()

// Theme role references (reactive to theme changes)
colors.PrimaryRoles.Primary
colors.SurfaceRoles.OnSurface
colors.SurfaceRoles.OnVariant

// For explicit colors, use SpecificColor
theme.ColorHelper.SpecificColor(color.NRGBA{R: 255, G: 0, B: 0, A: 255})
```

## Common Shape Options

```go
shape.ShapeRectangle                    // No rounding
shape.ShapeCircle                       // Full circle/ellipse
shape.RoundedCornerShape{Radius: 16}    // Uniform corners
shape.RoundedCornerShape{              // Per-corner control
    TopStart: 16,
    TopEnd:   16,
    BottomStart: 0,
    BottomEnd: 0,
}
```
