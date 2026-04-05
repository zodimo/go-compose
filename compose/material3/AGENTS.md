# Material 3 Components

## OVERVIEW
Material Design 3 implementation ported from Android Jetpack Compose.

## COMPONENTS

| Directory | Component | Description |
|-----------|-----------|-------------|
| appbar/ | AppBar | Top and bottom app bars |
| badge/ | Badge | Status indicators |
| bottomappbar/ | BottomAppBar | Bottom application bar |
| bottomsheet/ | BottomSheet | Modal bottom sheets |
| button/ | Button | Filled, FilledTonal, Outlined, Text, Elevated variants |
| card/ | Card | Filled, Elevated, Outlined variants |
| checkbox/ | Checkbox | Checkboxes |
| chip/ | Chip | Assist, Filter, Input, Suggestion variants |
| dialog/ | Dialog | Dialogs and alerts |
| divider/ | Divider | Horizontal dividers |
| floatingactionbutton/ | FAB | Floating action buttons |
| icon/ | Icon | Icon display |
| iconbutton/ | IconButton | Icon buttons |
| menu/ | Menu | Dropdown menus |
| navigationbar/ | NavigationBar | Bottom navigation |
| navigationdrawer/ | NavigationDrawer | Side navigation drawer |
| navigationrail/ | NavigationRail | Side navigation rail |
| progress/ | ProgressIndicator | Progress indicators |
| radiobutton/ | RadioButton | Radio buttons |
| scaffold/ | Scaffold | App scaffold layout |
| segmentedbutton/ | SegmentedButton | Segmented button groups |
| slider/ | Slider | Sliders and range sliders |
| snackbar/ | Snackbar | Snackbar notifications |
| surface/ | Surface | Surface containers |
| switch/ | Switch | Toggle switches |
| tab/ | Tab | Tab navigation |

### Theme Files
- color_scheme.go, color_set.go - Theme colors
- motion_scheme.go - Animation schemes
- shapes.go - Shape definitions

## PATTERNS

- next/ subdirectories hold v2 implementations (e.g., next/button/, next/textfield/)
- Each component package follows the structure:
  - Main widget file
  - options.go - Configuration options
  - defaults.go - Default values
- Color schemes are split across color_scheme.go and color_scheme_next.go

## DEMOS

See cmd/demo/<component>/ for usage examples:

- cmd/demo/appbar/ - App bar variations
- cmd/demo/card/ - Card examples
- cmd/demo/kitchen/ - Complete component showcase

To run a demo:
```bash
cd cmd/demo/<component>
go run .
```
