
# Roadmap: github.com/zodimo/go-compose

This roadmap outlines the path towards a complete, production-ready implementation of Material 3 components in Go, powered by [Gio](https://gioui.org/).

## Vision
To provide a feature-complete, idiomatic Go implementation of Jetpack Compose's Material 3 API, enabling developers to build beautiful, cross-platform applications with a familiar declarative syntax.

## Phases

### Phase 1: Structure & Navigation (App Shell)
*Focus: Enabling full application layouts.*

- [x] **Scaffold**: The core layout structure for screens.
- [x] **App Bars**:
    - [x] Top App Bar (Small, CenterAligned, Medium, Large)
    - [x] Bottom App Bar
- [x] **Bottom Navigation**: `NavigationBar` and `NavigationBarItem`.
- [x] **Tabs** (`TabRow`, `Tab`)
- [x] **Navigation Drawer**: Polish existing implementation and ensure full M3 compliance (standard vs modal).

### Phase 2: Lists & Data Display
*Focus: Efficiently displaying collections of data.*

- [ ] **Lists**: 
    - [x] Wrappers for `gio.List` to match Compose `LazyColumn`/`LazyRow` API.
    - [ ] Item spacing and content padding support.
- [ ] **Grids**: Lazy grids.
- [ ] **Chips**:
    - [x] Chips (Assist, Filter, Input, Suggestion) chips.
- [x] **Badges**: Integration with navigation items and icons.
- [ ] **Image**:
    - [ ] `Image` Composable. Must support all Modifiers (not just scaling). (Missing in `gio-mw`)
- [x] **Tooltips**: Plain and rich tooltips for desktop/mouse users. (Basic implementation done, known issues with buttons)

### Phase 3: Extended Inputs
*Focus: Complex user inputs.*

- [x] **Sliders**:
    - [x] Continuous and Discrete sliders.
    - [ ] Range sliders.
    - [x] Custom thumb and track support.
- [ ] **Pickers**:
    - [ ] Date Picker (Modal and Docked).
    - [ ] Time Picker (Dial and Input).
- [ ] **Segmented Button**: Single-select and multi-select variants.
- [ ] **Menus**:
    - [x] Dropdown Menu (polish existing implementation).
    - [ ] Exposed Dropdown Menu (ComboBox).

### Phase 4: Polish & Advanced Features
*Focus: Animation, accessibility, and desktop specifics.*

- [x] **Bottom Sheets**:
  - [x] Modal bottom sheet
  - [ ] Standard bottom sheet (Persistent)
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
| **Actions** | 🟡 Partial | [x] Floating Action Button (Basic implementation done), Segmented Button, Ext. FAB |
| **Communication** | 🟡 Partial | Tooltips, Badges |
| **Containment** | 🟡 Partial | Scaffold, Bottom Sheet |
| **Navigation** | 🔴 Early | App Bars, Bottom Nav |
| **Selection** | 🟡 Partial | Chips, Pickers |
| **Text Inputs** | 🟢 Good | Search Bar |

## Recent Milestones (Completed)
- [x] **Core Inputs**: Button, TextField, Checkbox, Radio, Switch.
- [x] **Surfaces**: Card, Surface (with shapes and shadows), Divider.
- [x] **Feedback**: Snackbar, Progress Indicators, Dialogs.
- [x] **Text**: Full typography support.
