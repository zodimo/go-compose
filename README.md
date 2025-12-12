# Testbed for go-compose

## Material 3 Component Inventory

Comparison of official Material 3 components with their implementation status in `go-compose-dev` and availability in the upstream `gio-mw` library.

| Material 3 Component | `gio-mw` Availability | `go-compose-dev` Status | Notes |
| :--- | :--- | :--- | :--- |
| **App bars: Bottom** | 🔴 Missing | 🔴 Todo | |
| **App bars: Top** | 🔴 Missing | 🔴 Todo | |
| **Badge** | 🟢 `widget/badge` | 🔴 Todo | |
| **Bottom sheets** | 🟢 `widget/sheet` | 🔴 Todo | |
| **Buttons** | 🟢 `widget/button` | 🟢 Implemented | |
| **Buttons: Extended FAB** | 🟢 `widget/button` | 🔴 Todo | Likely a variant of Button |
| **Buttons: FAB** | 🟢 `widget/button` | 🔴 Todo | Likely a variant of Button |
| **Buttons: Icon button** | 🟢 `widget/button` | 🟢 Implemented | Likely a variant of Button |
| **Buttons: Segmented** | 🔴 Missing | 🔴 Todo | |
| **Cards** | 🟢 `widget/card` | 🟢 Implemented | |
| **Carousel** | 🔴 Missing | 🔴 Todo | |
| **Checkbox** | 🟢 `widget/checkbox` | 🟢 Implemented | |
| **Chips** | 🔴 Missing | 🔴 Todo | |
| **Date pickers** | 🔴 Missing | 🔴 Todo | |
| **Dialogs** | 🟢 `widget/dialog` | 🟢 Implemented | |
| **Divider** | 🟢 `widget/divider` | 🟢 Implemented | |
| **Lists** | 🔴 Missing | 🔴 Todo | Core `gio.List` exists, but M3 styles missing |
| **Menus** | 🔴 Missing | 🟢 Implemented | |
| **Navigation bar** | 🔴 Missing | 🔴 Todo | Bottom Navigation |
| **Navigation drawer** | 🔴 Missing | 🟢 Implemented | Modal |
| **Navigation rail** | 🟢 `widget/rail` | 🟢 Implemented | |
| **Progress indicators** | 🟢 `widget/indicator` | 🟢 Implemented | |
| **Radio button** | 🟢 `widget/radio` | 🟢 Implemented | |
| **Search** |  `widget/search` | 🔴 Todo | |
| **Side sheets** |  `widget/sheet` | 🔴 Todo | |
| **Sliders** |  `widget/slider` | 🔴 Todo | |
| **Snackbars** | 🟢 `widget/snackbar` | 🟢 Implemented | |
| **Switch** | 🟢 `widget/toggle` | 🟢 Implemented | |
| **Surface** | 🟢 `exp` | 🟢 Implemented | Core primitive (Shadows, Borders, Shapes) |
| **Tabs** |  `widget/tab` | 🔴 Todo | |
| **Text fields** | 🟢 `widget/input` | 🟢 Implemented | |
| **Time pickers** | 🔴 Missing | 🔴 Todo | |
| **Tooltips** | � `widget/tooltip` | 🔴 Todo | |
