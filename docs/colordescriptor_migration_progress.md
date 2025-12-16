# ColorDescriptor Migration Progress

This document tracks the migration of Material3 components to use the `ColorDescriptor` pattern with theme-aware defaults.

## Migration Status

Last Updated: 2025-12-16 (Complete Component Audit)

### ✅ Migrated to ColorDescriptor (4 components)

Components that use `theme.ColorDescriptor`:

| Component | Date Completed | Notes |
|-----------|----------------|-------|
| **surface** | - | Foundation component, uses ColorDescriptor for Color, ContentColor, BorderColor |
| **bottomappbar** | 2025-12-16 | Uses ColorDescriptor for ContainerColor, ContentColor |
| **navigationrail** | - | Uses surface wrapper with SpecificColor (relies on surface's ColorDescriptor) |
| **tab** | 2025-12-16 | Uses ColorDescriptor for TabRow ContainerColor/ContentColor and Tab SelectedContentColor/UnselectedContentColor |

### 📋 Pending ColorDescriptor Migration (9 components)

Native go-compose components using `color.Color` or `color.NRGBA` that should be migrated:

| Component | Color Fields | Priority | Notes |
|-----------|--------------|----------|-------|
| **appbar** | `TopAppBarColors` struct with 5 color.Color fields | High | Layout container |
| **badge** | `ContainerColor`, `ContentColor` (NRGBA) | Medium | Simple component |
| **chip** | `Color`, `BorderColor` (NRGBA) | High | Selection component |
| **floatingactionbutton** | `ContainerColor`, `ContentColor` (Color) | High | Primary action |
| **navigationbar** | `ContainerColor`, `ContentColor` (Color) | High | Navigation |
| **scaffold** | `ContainerColor`, `ContentColor` (Color) | High | Layout container |
| **segmentedbutton** | `SelectedColor`, `UnselectedColor`, `SelectedContentColor`, `UnselectedContentColor`, `BorderColor` (NRGBA) | High | Selection component |
| **slider** | `SliderColors` struct with 10 color.NRGBA fields | Medium | Input control |
| **divider** | `Color` (Color) | Low | Simple component, but also uses gio-mw |

### 🔗 External Widget Components (13 confirmed)

Components that wrap gio-mw widgets and don't expose color options:

| Component | Widget/Package | Notes |
|-----------|----------------|-------|
| **button** | `gio-mw/widget/button` | No color options exposed |
| **card** | `gio-mw/widget/card` | No color options exposed |
| **checkbox** | `gio-mw/widget/checkbox` | No color options exposed |
| **dialog** | `gio-mw/widget/dialog` | No color options exposed |
| **iconbutton** | `gio-mw/widget/button` | No color options exposed |
| **menu** | `gio-mw/token` | Uses gio-mw tokens |
| **progress** | `gio-mw/widget/indicator` | Wrapper around gio-mw progress indicator |
| **radiobutton** | `gio-mw/wdk` | No color options exposed |
| **snackbar** | `gio-mw/widget/snackbar` | Wrapper around gio-mw snackbar widget |
| **switch** | `gio-mw/widget/toggle` | No color options exposed |
| **text** | `gio-mw/token`, `gio-mw/wdk` | Uses gio-mw typography helpers |
| **textfield** | `gio-mw/widget/input` | Wrapper around gio-mw input widget |
| **tooltip** | `gio-mw/widget/tooltip` | Wrapper around gio-mw tooltip widget |

> [!NOTE]
> These components use the gio-mw widget library's Material theme system. They automatically respond to theme changes but do not expose color customization options.

### ⚠️ Special Cases (2 components)

| Component | Status | Notes |
|-----------|--------|-------|
| **bottomsheet** | Uses `token.MatColor` (gio-mw type) | Hybrid: native layout but uses gio-mw color types. Needs investigation whether to migrate or keep as-is |
| **navigationdrawer** | No color options | Native component but doesn't expose color customization |

## Migration Workflow

For step-by-step migration instructions, see:
- [Migration Workflow](.agent/workflows/migrate_component_colordescriptor.md)

## Quick Reference

### Common Theme Role Mappings
- **Container backgrounds** → `SurfaceRoles.Surface*` (Surface, SurfaceContainer, SurfaceContainerHigh, etc.)
- **Content/text** → `ContentRoles.OnSurface`, `OnPrimary`, `OnSecondary`, etc.
- **Borders/outlines** → `OutlineRoles.Outline`, `OutlineVariant`
- **Primary actions** → `PrimaryRoles.Primary`, `PrimaryContainer`
- **State layers** → Use `.SetOpacity()` on base colors

### Migration Checklist
- [ ] Analyze color usage
- [ ] Update Options structure to use ColorDescriptor
- [ ] Set theme-aware defaults using color roles
- [ ] Update option setters to accept ColorDescriptor
- [ ] Remove SpecificColor() wrappers (keep only for non-theme colors)
- [ ] Update internal color resolution if needed
- [ ] Update tests/demos
- [ ] Verify build and visual appearance
- [ ] Update documentation

## Architecture Notes

### Complete Audit Summary (2025-12-16)

**Total Material3 Components: 28**

**ColorDescriptor Migration Progress: 4 of 13 native components (31%)**

| Category | Count | Components |
|----------|-------|------------|
| ✅ **Migrated** | 4 | surface, bottomappbar, navigationrail, tab |
| 📋 **Pending** | 9 | appbar, badge, chip, floatingactionbutton, navigationbar, scaffold, segmentedbutton, slider, divider |
| 🔗 **External widgets** | 13 | button, card, checkbox, dialog, iconbutton, menu, progress, radiobutton, snackbar, switch, text, textfield, tooltip |
| ⚠️ **Special cases** | 2 | bottomsheet (uses token.MatColor), navigationdrawer (no colors) |

### Migration Priority

Based on component importance and usage:

**High Priority (6 components):**
- `appbar` - Layout container, widely used
- `chip` - Selection component
- `floatingactionbutton` - Primary action button
- `navigationbar` - Navigation component
- `scaffold` - Core layout container
- `segmentedbutton` - Selection component

**Medium Priority (2 components):**
- `badge` - Notification component
- `slider` - Input control (10 color fields)

**Low Priority (1 component):**
- `divider` - Simple visual separator

### Component Type Breakdown

**Native Components (13):** Built in go-compose, support or should support ColorDescriptor
- 4 migrated: surface, bottomappbar, navigationrail, tab
- 9 pending migration: appbar, badge, chip, floatingactionbutton, navigationbar, scaffold, segmentedbutton, slider, divider

**External Widget Components (13):** Wrap gio-mw widgets, use gio-mw themes
- No color customization exposed
- Automatically theme-aware via gio-mw

**Special Cases (2):**
- bottomsheet: Hybrid (native layout, gio-mw color types)
- navigationdrawer: Native but minimal color customization

## Notes

- 4 of 13 native components now use `theme.ColorDescriptor` (31% complete)
- 9 high-value native components still need migration
- External widget components (13) use gio-mw themes - no migration needed
- All pending components have been verified by file inspection
- `SpecificColor()` should only wrap non-theme colors in migrated components
- Theme role selectors ensure proper light/dark theme support
