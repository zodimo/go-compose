---
description: How to apply styling (clip, border, colors) to UI elements idiomatically
---

# GoCompose Styling Patterns

When building UI with GoCompose, follow these patterns for styling.

## 1. Always Prefer Modifier Chains for Single Element Styling

For clip, border, background on a single element, use the modifier chain pattern:

```go
fImage.Image(
    imageResource,
    fImage.WithModifier(
        size.Size(100, 100).
            Then(clip.Clip(shape.ShapeCircle)).
            Then(border.Border(4, color, shape.ShapeCircle)),
    ),
)
```

**Do NOT** wrap in Surface just for clip/border styling.

## 2. Use Surface Only When Needed

Use Surface when you need:
- A container with background + content
- Multiple children with shared styling
- Semantic elevation (cards, dialogs, FABs)

## 3. Use ColorSelector for Theme Colors

```go
colors := theme.ColorHelper.ColorSelector()
colors.PrimaryRoles.Primary
colors.SurfaceRoles.OnSurface
colors.SurfaceRoles.OnVariant
```

**Do NOT** use:
```go
// Avoid this pattern
theme.ColorHelper.SpecificColor(m3.Scheme.Primary.Color.AsNRGBA())
```

## Reference

See [docs/styling_patterns.md](file:///home/jaco/SecondBrain/1-Projects/GoCompose/development/go-compose/docs/styling_patterns.md) for the full guide.
