---
description: Migrate a Material3 component to use ColorDescriptor with theme defaults
---

# Workflow: Migrate Component to ColorDescriptor

This workflow systematically migrates a Material3 component to use the ColorDescriptor pattern with appropriate theme defaults.

## Prerequisites
- Foundation layer (background, border, surface) already migrated
- Component compiles (colors wrapped with SpecificColor if needed)

## Step 1: Analyze Component Structure

**Identify color-related fields in the component:**
```bash
# Find the component's options file
find compose/foundation/material3/<COMPONENT> -name "options.go"

# Search for color-related fields
grep -n "Color\|color\." compose/foundation/material3/<COMPONENT>/options.go
```

**Document current color usage:**
- List all option fields using `color.Color`, `color.NRGBA`, or `token.MatColor`
- Note which colors are container colors, content colors, border colors, etc.
- Check if there are color-related defaults functions

## Step 2: Update Options Structure

**In `options.go`, convert color fields to ColorDescriptor:**

```go
// BEFORE:
type ComponentOptions struct {
    ContainerColor color.Color
    ContentColor   color.NRGBA
    BorderColor    token.MatColor
}

// AFTER:
type ComponentOptions struct {
    ContainerColor theme.ColorDescriptor
    ContentColor   theme.ColorDescriptor  
    BorderColor    theme.ColorDescriptor
}
```

**Update imports:**
```go
import (
    "github.com/zodimo/go-compose/theme"
    // Remove: "image/color" (unless still needed elsewhere)
    // Remove: "git.sr.ht/~schnwalter/gio-mw/token" (for MatColor)
)
```

## Step 3: Set Theme-Aware Defaults

**In `DefaultComponentOptions()` function, use theme roles:**

```go
// BEFORE:
func DefaultComponentOptions() ComponentOptions {
    return ComponentOptions{
        ContainerColor: color.NRGBA{R: 234, G: 221, B: 255, A: 255}, // Hardcoded
        ContentColor:   color.NRGBA{R: 33, G: 0, B: 93, A: 255},
    }
}

// AFTER:
func DefaultComponentOptions() ComponentOptions {
    return ComponentOptions{
        // Use theme color selectors for semantic colors
        ContainerColor: theme.ColorHelper.ColorSelector().SurfaceRoles.SurfaceContainerHigh,
        ContentColor:   theme.ColorHelper.ColorSelector().ContentRoles.OnSurface,
        BorderColor:    theme.ColorHelper.ColorSelector().OutlineRoles.Outline,
    }
}
```

**Common theme role mappings:**
- Container backgrounds → `SurfaceRoles.Surface*` (Surface, SurfaceContainer, SurfaceContainerHigh, etc.)
- Content/text → `ContentRoles.OnSurface`, `OnPrimary`, `OnSecondary`, etc.
- Borders/outlines → `OutlineRoles.Outline`, `OutlineVariant`
- Primary actions → `PrimaryRoles.Primary`, `PrimaryContainer`
- State layers → Use `.SetOpacity()` on base colors

## Step 4: Update Option Setters

**Convert setter functions to accept ColorDescriptor:**

```go
// BEFORE:
func WithContainerColor(c color.Color) ComponentOption {
    return func(o *ComponentOptions) {
        o.ContainerColor = c
    }
}

// AFTER:
func WithContainerColor(colorDesc theme.ColorDescriptor) ComponentOption {
    return func(o *ComponentOptions) {
        o.ContainerColor = colorDesc
    }
}
```

## Step 5: Remove SpecificColor Wrapping

**In the main component file, remove `SpecificColor()` wrappers:**

Since options now already contain `ColorDescriptor`, you can pass them directly:

```go
// BEFORE (temporary wrapper):
surface.WithColor(theme.ColorHelper.SpecificColor(opts.ContainerColor))

// AFTER (options are already ColorDescriptor):
surface.WithColor(opts.ContainerColor)
```

## Step 6: Update Component Logic

**If component resolves colors internally (not via surface):**

```go
// Get ThemeManager
tm := theme.GetThemeManager()

// Resolve ColorDescriptor to NRGBA at layout time
resolvedColor := tm.ResolveColorDescriptor(opts.ContainerColor)

// Use resolved color for drawing
paint.ColorOp{Color: resolvedColor}.Add(gtx.Ops)
```

## Step 7: Update Tests and Demos

**Update any tests using the component:**
```go
// BEFORE:
Component(
    WithContainerColor(color.NRGBA{R: 255, G: 0, B: 0, A: 255}),
)

// AFTER - for explicit colors:
Component(
    WithContainerColor(theme.ColorHelper.SpecificColor(color.NRGBA{R: 255, G: 0, B: 0, A: 255})),
)

// AFTER - preferred, using theme roles:
Component(
    WithContainerColor(theme.ColorHelper.ColorSelector().PrimaryRoles.Primary),
)
```

## Step 8: Verify Migration

**Run checks:**
```bash
// turbo
# Build the component
go build ./compose/foundation/material3/<COMPONENT>/...

// turbo
# Run component tests if they exist
go test ./compose/foundation/material3/<COMPONENT>/...

// turbo
# Build all dependent code
go build ./...
```

**Visual verification:**
- Run demo if it exists: `go run ./cmd/demo/<COMPONENT>/`
- Verify colors change appropriately with theme switches
- Check that color transformations (opacity, lighten, darken) work

## Step 9: Document Changes

**Update component documentation:**
- Note that colors are now theme-aware
- Document which theme roles are used by default
- Show examples of color customization

## Checklist Template

Use this checklist for each component migration:

```markdown
## Component: <NAME>

- [ ] Step 1: Analyzed color usage (X color fields found)
- [ ] Step 2: Updated Options structure to use ColorDescriptor
- [ ] Step 3: Set theme-aware defaults using color roles
- [ ] Step 4: Updated option setters to accept ColorDescriptor
- [ ] Step 5: Removed SpecificColor() wrappers
- [ ] Step 6: Updated internal color resolution (if applicable)
- [ ] Step 7: Updated tests/demos
- [ ] Step 8: Verified build and visual appearance
- [ ] Step 9: Updated documentation
```

## Common Patterns

### For State-Dependent Colors
```go
// Use conditional ColorDescriptor selection
containerColor := theme.ColorHelper.ColorSelector().SurfaceRoles.Surface
if enabled {
    containerColor = theme.ColorHelper.ColorSelector().PrimaryRoles.PrimaryContainer
}
```

### For Interactive States (hover, pressed, etc.)
```go
// Apply opacity transformations
stateLayerColor := theme.ColorHelper.ColorSelector().ContentRoles.OnSurface.
    SetOpacity(theme.OpacityLevel8) // 8% opacity for hover
```

### For Gradients (future)
```go
// Use theme-aware gradient config
gradient := theme.VerticalGradient(
    theme.ColorHelper.ColorSelector().PrimaryRoles.Primary,
    theme.ColorHelper.ColorSelector().PrimaryRoles.PrimaryContainer,
)
```

## Priority Order

Migrate components in this order (dependencies first):

1. **Core building blocks** - Already done: surface, background, border, icon
2. **Interactive primitives** - button, iconbutton, checkbox, radio, switch
3. **Input components** - textfield, textarea
4. **Selection components** - chip (done), segmentedbutton (done)
5. **Navigation** - tab (done), navigationbar (done), navigationrail (done), navigationdrawer (done)
6. **Layout containers** - appbar (done), scaffold (done), bottomappbar (done)
7. **Feedback** - snackbar, dialog, bottomsheet (done)
8. **Lists** - listitem, divider
9. **Complex components** - card, menu, dropdown

## Notes

- Always prefer theme roles over `SpecificColor()` for component defaults
- Use `SpecificColor()` only for truly custom/branded colors in demos
- Remember to add `theme` import and remove unused `color`/`token` imports
- Test with both light and dark themes to ensure proper contrast

## Progress Tracking

Track migration progress and see the status of all components:
- [ColorDescriptor Migration Progress](../docs/colordescriptor_migration_progress.md)

**After completing a migration**, update the progress document by:
1. Moving the component from "Pending" to "Completed" section
2. Adding the completion date
3. Adding any relevant notes about the migration
