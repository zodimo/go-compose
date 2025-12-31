// Package selection provides text selection functionality for Compose Foundation.
//
// This package contains the SelectionContainer composable and related types
// for enabling text selection across multiple text composables.
//
// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/foundation/foundation/src/commonMain/kotlin/androidx/compose/foundation/text/selection/
package selection

import (
	"github.com/zodimo/go-compose/compose/ui/next/text"
	"github.com/zodimo/go-compose/compose/ui/next/text/style"
)

// AnchorInfo contains information about a selection anchor (start/end).
//
// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/foundation/foundation/src/commonMain/kotlin/androidx/compose/foundation/text/selection/Selection.kt
type AnchorInfo struct {
	// Direction is the text direction of the character at selection edge.
	Direction style.ResolvedTextDirection
	// Offset is the character offset for the selection edge.
	// This offset is within the individual child text composable.
	Offset int
	// SelectableId is the id of the Selectable which contains this anchor.
	SelectableId int64
}

// Selection represents the current selection state.
//
// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/foundation/foundation/src/commonMain/kotlin/androidx/compose/foundation/text/selection/Selection.kt
type Selection struct {
	// Start contains information about the start of the selection.
	Start AnchorInfo
	// End contains information about the end of the selection.
	End AnchorInfo
	// HandlesCrossed is true when the user drags one handle to cross the other handle.
	// When selection happens in a single widget, checking TextRange.start > TextRange.end
	// is enough. But when selection happens across multiple widgets, this value needs
	// more complicated calculation.
	HandlesCrossed bool
}

// NewSelection creates a new Selection with the given parameters.
func NewSelection(start, end AnchorInfo, handlesCrossed bool) Selection {
	return Selection{
		Start:          start,
		End:            end,
		HandlesCrossed: handlesCrossed,
	}
}

// Merge merges this selection with another selection.
// If other is nil, returns this selection unchanged.
func (s Selection) Merge(other *Selection) Selection {
	if other == nil {
		return s
	}

	if s.HandlesCrossed || other.HandlesCrossed {
		var newStart, newEnd AnchorInfo
		if other.HandlesCrossed {
			newStart = other.Start
		} else {
			newStart = other.End
		}
		if s.HandlesCrossed {
			newEnd = s.End
		} else {
			newEnd = s.Start
		}
		return Selection{
			Start:          newStart,
			End:            newEnd,
			HandlesCrossed: true,
		}
	}

	return Selection{
		Start:          s.Start,
		End:            other.End,
		HandlesCrossed: false,
	}
}

// ToTextRange returns the selection offset information as a TextRange.
func (s Selection) ToTextRange() text.TextRange {
	return text.NewTextRange(s.Start.Offset, s.End.Offset)
}

// Copy returns a copy of this Selection.
func (s Selection) Copy() Selection {
	return Selection{
		Start:          s.Start,
		End:            s.End,
		HandlesCrossed: s.HandlesCrossed,
	}
}

// IsCollapsed returns true if the selection start and end are at the same position.
func (s Selection) IsCollapsed() bool {
	return s.Start.Offset == s.End.Offset && s.Start.SelectableId == s.End.SelectableId
}
