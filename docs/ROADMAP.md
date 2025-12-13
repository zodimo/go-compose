
# Roadmap: go-compose-dev

This roadmap outlines the path towards a complete, production-ready implementation of Material 3 components in Go, powered by [Gio](https://gioui.org/).

## Vision
To provide a feature-complete, idiomatic Go implementation of Jetpack Compose's Material 3 API, enabling developers to build beautiful, cross-platform applications with a familiar declarative syntax.

## Phases

### Phase 1: Structure & Navigation (App Shell)
*Focus: Enabling full application layouts.*

- [ ] **Scaffold**: The core layout structure for screens.
- [ ] **App Bars**:
    - [ ] Top App Bar (Small, CenterAligned, Medium, Large)
    - [ ] Bottom App Bar
- [ ] **Bottom Navigation**: `NavigationBar` and `NavigationBarItem`.
- [ ] **Navigation Drawer**: Polish existing implementation and ensure full M3 compliance (standard vs modal).

### Phase 2: Lists & Data Display
*Focus: Efficiently displaying collections of data.*

- [ ] **Lists**: 
    - [ ] Wrappers for `gio.List` to match Compose `LazyColumn`/`LazyRow` API.
    - [ ] Item spacing and content padding support.
- [ ] **Grids**: Lazy grids.
- [ ] **Chips**:
    - [ ] Assist, Filter, Input, and Suggestion chips.
- [ ] **Badges**: Integration with navigation items and icons.
- [ ] **Tooltips**: Plain and rich tooltips for desktop/mouse users.

### Phase 3: Extended Inputs
*Focus: Complex user inputs.*

- [ ] **Sliders**:
    - [ ] Continuous and Discrete sliders.
    - [ ] Range sliders.
    - [ ] Custom thumb and track support.
- [ ] **Pickers**:
    - [ ] Date Picker (Modal and Docked).
    - [ ] Time Picker (Dial and Input).
- [ ] **Segmented Button**: Single-select and multi-select variants.
- [ ] **Menus**:
    - [ ] Dropdown Menu (polish existing implementation).
    - [ ] Exposed Dropdown Menu (ComboBox).

### Phase 4: Polish & Advanced Features
*Focus: Animation, accessibility, and desktop specifics.*

- [ ] **Animations**:
    - [ ] Shared element transitions.
    - [ ] `AnimatedContent` wrappers.
- [ ] **Accessibility**: Ensure all components export semantic information correctly for screen readers.
- [ ] **Desktop Support**:
    - [ ] Keyboard shortcuts integration.
    - [ ] Context menus.
    - [ ] Cursor handling improvements.

## Component Status Overview

| Component Group | Status | Key Missing Items |
| :--- | :--- | :--- |
| **Actions** | 🟡 Partial | Segmented Button, Ext. FAB |
| **Communication** | 🟡 Partial | Tooltips, Badges |
| **Containment** | 🟡 Partial | Scaffold, Bottom Sheet |
| **Navigation** | 🔴 Early | App Bars, Bottom Nav |
| **Selection** | 🟡 Partial | Chips, Pickers, Sliders |
| **Text Inputs** | 🟢 Good | Search Bar |

## Recent Milestones (Completed)
- [x] **Core Inputs**: Button, TextField, Checkbox, Radio, Switch.
- [x] **Surfaces**: Card, Surface (with shapes and shadows), Divider.
- [x] **Feedback**: Snackbar, Progress Indicators, Dialogs.
- [x] **Text**: Full typography support.
