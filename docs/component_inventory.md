# Component Inventory & Gap Analysis

This document tracks the status of Jetpack Compose components within `go-compose-dev` and their availability in the upstream `gio-mw` library.

**Status Legend:**
- ✅ **Implemented**: Available in `go-compose-dev/compose/foundation/material3`.
- 📦 **Available in gio-mw**: Exists in `gio-mw` but needs porting/wrapping in `go-compose-dev`.
- ❌ **Missing**: Not currently implemented in either.
- 🚧 **Partial**: Started but incomplete.

## Actions

| Component | Status | `gio-mw` | Notes |
| :--- | :--- | :--- | :--- |
| **Button** | ✅ Implemented | `widget/button` | `compose/foundation/material3/button` |
| **Floating Action Button** | ✅ Implemented | `widget/button` | Can be derived from Button or typically styled in `gio-mw`. |
| **Icon Button** | ✅ Implemented | `widget/button` | `compose/foundation/material3/iconbutton` |
| **Segmented Button** | ❌ Missing | - | |
| **Extended FAB** | ❌ Missing | - | |

## Communication

| Component | Status | `gio-mw` | Notes |
| :--- | :--- | :--- | :--- |
| **Badges** | ✅ Implemented | `compose/foundation/material3/badge` | |
| **Progress Indicators** | ✅ Implemented | `widget/indicator` | `compose/foundation/material3/progress` |
| **Snackbar** | ✅ Implemented | `widget/snackbar` | `compose/foundation/material3/snackbar` |
| **Tooltips** | 📦 Available | `widget/tooltip` | |

## Containment

| Component | Status | `gio-mw` | Notes |
| :--- | :--- | :--- | :--- |
| **Bottom Sheets** | 📦 Available | `widget/sheet` | |
| **Cards** | ✅ Implemented | `widget/card` | `compose/foundation/material3/card` |
| **Carousel** | ❌ Missing | - | |
| **Dialogs** | ✅ Implemented | `widget/dialog` | `compose/foundation/material3/dialog` |
| **Dividers** | ✅ Implemented | `widget/divider` | `compose/foundation/material3/divider` |
| **Lists** | ✅ Implemented | Core Gio | Implemented `LazyColumn` and `LazyRow` wrappers (Eager composition, Lazy layout). |
| **Scaffold** | ✅ Implemented | `compose/foundation/material3/scaffold` | High priority for app structure. |
| **Surface** | ✅ Implemented | - | `compose/foundation/material3/surface`. Fundamental building block. |

## Navigation

| Component | Status | `gio-mw` | Notes |
| :--- | :--- | :--- | :--- |
| **App Bars** | ✅ Implemented | - | Top and Bottom App Bars implemented. |
| **Navigation Bar** | ✅ Implemented | - | Bottom Navigation (`navigationbar`). |
| **Navigation Drawer** | ✅ Implemented | - | - [x] Navigation Drawer (Modal) - [x] Navigation Drawer Item |
| **Navigation Rail** | ✅ Implemented | `widget/rail` | `compose/foundation/material3/navigationrail` (Prototype Implemented) |
| **Tabs** | 📦 Available | `widget/tab` | |

## Selection

| Component | Status | `gio-mw` | Notes |
| :--- | :--- | :--- | :--- |
| **Checkbox** | ✅ Implemented | `widget/checkbox` | `compose/foundation/material3/checkbox` |
| **Chips** | ✅ Implemented | `compose/foundation/material3/chip` | Assist, Filter, Input, Suggestion chips. |
| **Date Picker** | ❌ Missing | - | |
| **Menus** | ✅ Implemented | `widget/overlay` | `compose/foundation/material3/menu`. Using Card for now (inline). |
| **Radio Button** | ✅ Implemented | `widget/radio` | `compose/foundation/material3/radiobutton` |
| **Sliders** | 📦 Available | `widget/slider` | |
| **Switch** | ✅ Implemented | `widget/toggle` | `compose/foundation/material3/switch` |
| **Time Pickers** | ❌ Missing | - | |

## Text Inputs

| Component | Status | `gio-mw` | Notes |
| :--- | :--- | :--- | :--- |
| **Search** | 📦 Available | `widget/search` | |
| **Text Fields** | ✅ Implemented | `widget/input` | `compose/foundation/material3/textfield` |
| **Text** | ✅ Implemented | - | `compose/foundation/material3/text`. Renders text with typography. |

## Summary
- **Strong Foundation**: Core inputs (Text, Checkbox, Radio, Switch) and containers (Card, Surface, Dialog) are ready.
- **Rich Middleware**: `gio-mw` offers a lot of "low hanging fruit" to port: Progress, Slider, Tabs, Snackbar.
- **Structural Gaps**: `Scaffold` and Navigation components (App Bars, Drawers) are major missing pieces for full app shells.
