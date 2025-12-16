# ColorDescriptor Migration Progress

This document tracks the migration of Material3 components to use the `ColorDescriptor` pattern with theme-aware defaults.

## Migration Status

Last Updated: 2025-12-16

### ✅ Completed Components

Components that have been fully migrated to use `ColorDescriptor`:

| Component | Status | Date Completed | Notes |
|-----------|--------|----------------|-------|
| **background** | ✅ Complete | - | Foundation layer |
| **border** | ✅ Complete | - | Foundation layer |
| **surface** | ✅ Complete | - | Foundation layer |
| **icon** | ✅ Complete | - | Foundation layer |
| **chip** | ✅ Complete | - | Selection component |
| **segmentedbutton** | ✅ Complete | - | Selection component |
| **tab** | ✅ Complete | - | Navigation component |
| **navigationbar** | ✅ Complete | - | Navigation component |
| **navigationrail** | ✅ Complete | - | Navigation component |
| **navigationdrawer** | ✅ Complete | - | Navigation component |
| **appbar** | ✅ Complete | - | Layout container |
| **scaffold** | ✅ Complete | - | Layout container |
| **bottomappbar** | ✅ Complete | 2025-12-16 | Layout container, verified correct |
| **bottomsheet** | ✅ Complete | - | Feedback component |

### 🔄 In Progress

Components currently being migrated:

| Component | Status | Assignee | Notes |
|-----------|--------|----------|-------|
| - | - | - | - |

### 📋 Pending Migration

Components that still need to be migrated, organized by priority:

#### High Priority (Interactive Primitives)
- [ ] **button** - Core interactive component
- [ ] **iconbutton** - Icon-based button
- [ ] **checkbox** - Selection control
- [ ] **radio** - Selection control
- [ ] **switch** - Toggle control
- [ ] **floatingactionbutton** - Primary action button

#### Medium Priority (Input Components)
- [ ] **textfield** - Text input
- [ ] **textarea** - Multi-line text input

#### Medium Priority (Lists)
- [ ] **listitem** - List item component
- [ ] **divider** - Visual separator

#### Lower Priority (Complex Components)
- [ ] **card** - Container component
- [ ] **menu** - Popup menu
- [ ] **dropdown** - Selection dropdown
- [ ] **snackbar** - Feedback notification
- [ ] **dialog** - Modal dialog

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

## Notes

- Components marked as complete have been verified to use `ColorDescriptor` for theme colors
- `SpecificColor()` should only be used for truly custom/branded colors or transparent overlays
- All theme-based colors should use theme role selectors for proper light/dark theme support
