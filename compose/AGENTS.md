# Compose Framework

## OVERVIEW
Core declarative UI framework implementing Jetpack Compose patterns on top of gioui.org.

## HIERARCHY
- material3/ - Material Design 3 components (buttons, cards, dialogs, etc.)
- foundation/ - Foundation components (text, layout, lazy lists, icons)
- ui/ - Low-level UI primitives (graphics, text rendering, modifiers)

## HOW TO ADD COMPONENTS
1. Create package in appropriate layer
2. Add alias.go to re-export from implementation
3. Follow "next" pattern for new implementations

## ALIAS PATTERN
The project uses alias.go files to create stable public APIs:

```go
// compose/foundation/next/text/alias.go
package text

import (
    "github.com/zodimo/go-compose/compose/ui/next/text"
    "github.com/zodimo/go-compose/pkg/api"
)

type Composable = api.Composable
type Composer = api.Composer
type TextStyle = text.TextStyle
```

Why: API layer can change implementation location without breaking users.

## "NEXT" DIRECTORIES
Many packages have parallel implementations:
- text/ and next/text/
- button/ and next/button/

The next/ version is newer and often has different APIs. For new development, prefer next/.

## COMPONENT STRUCTURE
Typical component layout:

```
button/
├── alias.go         # Public API re-exports
├── button.go        # Main widget implementation
├── options.go       # Option functions
├── defaults.go      # Default values
└── button_test.go   # Tests
```

## SEE ALSO
- ./material3/AGENTS.md - Material 3 components
- ./foundation/AGENTS.md - Foundation components
- ../pkg/api/types.go - Core API types (Composable, Composer)
