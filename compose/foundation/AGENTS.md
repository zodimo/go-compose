# Foundation Components

## OVERVIEW
Low-level UI primitives: text rendering, layouts, lazy lists, icons, images.

## MODULES
- border_stroke.go - Border stroke definitions
- icon/ - Icon rendering from Material Design icons
- image/ - Image display and loading
- layout/ - Layout primitives
  - box/ - Box container
  - column/ - Vertical stacking
  - row/ - Horizontal stacking
  - overlay/ - Stack children on top of each other
  - spacer/ - Empty space
  - tree/ - Tree data structure layout
- lazy/ - LazyColumn and LazyRow for virtualization
- next/ - Next-generation implementations
  - text/ - Advanced text with selection, input
  - textfield/ - Text input fields
- stopclickthru/ - Click event blocking
- text/ - Basic text rendering
- textfield/ - Basic text input

## TEXT SYSTEM
Three text implementations:
1. text/ - Basic text (stable)
2. next/text/ - Advanced text with selection, auto-sizing
3. next/textfield/ - Text input with complex editing

## LAYOUT SYSTEM
Layout modifiers in modifiers/ directory:
- size/ - Width, height, fill
- padding/ - Insets
- weight/ - Flex distribution
- background/ - Background colors/shapes

## KNOWN ISSUES
- Word navigation only whitespace-delimited (affects CJK languages)
- Selection system has 11+ TODOs for handle rendering, position calc
- See TODOs in next/text/selection/
