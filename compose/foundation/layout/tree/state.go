package tree

import (
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/state"
)

// TreeState controls the state of a Tree component, such as which branches are expanded.
type TreeState struct {
	expandedItems state.MutableValueTyped[map[any]bool]
	selectedItems state.MutableValueTyped[map[any]bool] // Using map for O(1) lookup, even for single selection
}

// RememberTreeState creates a TreeState that is remembered across compositions.
func RememberTreeState(c api.Composer) *TreeState {
	expanded := c.State("treeState/expanded", func() any {
		return make(map[any]bool)
	})
	selected := c.State("treeState/selected", func() any {
		return make(map[any]bool)
	})

	expandedTyped, _ := state.MutableValueToTyped[map[any]bool](expanded)
	selectedTyped, _ := state.MutableValueToTyped[map[any]bool](selected)

	return &TreeState{
		expandedItems: expandedTyped,
		selectedItems: selectedTyped,
	}
}

// IsExpanded checks if the given node ID is expanded.
func (s *TreeState) IsExpanded(id any) bool {
	expanded := s.expandedItems.Get()
	return expanded[id]
}

// Expand marks the given node ID as expanded.
func (s *TreeState) Expand(id any) {
	s.expandedItems.Update(func(current map[any]bool) map[any]bool {
		if current[id] {
			return current
		}
		newMap := make(map[any]bool, len(current)+1)
		for k, v := range current {
			newMap[k] = v
		}
		newMap[id] = true
		return newMap
	})
}

// Collapse marks the given node ID as collapsed.
func (s *TreeState) Collapse(id any) {
	s.expandedItems.Update(func(current map[any]bool) map[any]bool {
		if !current[id] {
			return current
		}
		newMap := make(map[any]bool, len(current))
		for k, v := range current {
			if k != id {
				newMap[k] = v
			}
		}
		return newMap
	})
}

// Toggle flips the expansion state of the given node ID.
func (s *TreeState) Toggle(id any) {
	if s.IsExpanded(id) {
		s.Collapse(id)
	} else {
		s.Expand(id)
	}
}

// IsSelected checks if the given node ID is selected.
func (s *TreeState) IsSelected(id any) bool {
	selected := s.selectedItems.Get()
	return selected[id]
}

// Select marks the given node ID as selected.
// If singleSelection is true, all other nodes are unselected.
func (s *TreeState) Select(id any, singleSelection bool) {
	s.selectedItems.Update(func(current map[any]bool) map[any]bool {
		if singleSelection {
			return map[any]bool{id: true}
		}
		if current[id] {
			return current
		}
		newMap := make(map[any]bool, len(current)+1)
		for k, v := range current {
			newMap[k] = v
		}
		newMap[id] = true
		return newMap
	})
}

// Unselect marks the given node ID as unselected.
func (s *TreeState) Unselect(id any) {
	s.selectedItems.Update(func(current map[any]bool) map[any]bool {
		if !current[id] {
			return current
		}
		newMap := make(map[any]bool, len(current))
		for k, v := range current {
			if k != id {
				newMap[k] = v
			}
		}
		return newMap
	})
}

// ClearSelection removes all selections.
func (s *TreeState) ClearSelection() {
	s.selectedItems.Set(make(map[any]bool))
}

// GetSelectedItems returns a slice of all selected item IDs.
func (s *TreeState) GetSelectedItems() []any {
	selected := s.selectedItems.Get()
	result := make([]any, 0, len(selected))
	for id := range selected {
		result = append(result, id)
	}
	return result
}

// GetFirstSelected returns the first selected item ID, or nil if none selected.
// Useful for single-selection mode.
func (s *TreeState) GetFirstSelected() any {
	selected := s.selectedItems.Get()
	for id := range selected {
		return id
	}
	return nil
}

// OpenAllBranches expands all branches in the tree.
// roots: the root node IDs.
// childUIDs: function to get children IDs for a given ID.
// isBranch: function to check if a node is a branch (optional, if nil, checks children > 0).
func (s *TreeState) OpenAllBranches(roots []any, childUIDs func(any) []any, isBranch func(any) bool) {
	newExpanded := make(map[any]bool)

	var walk func(ids []any)
	walk = func(ids []any) {
		for _, id := range ids {
			isB := false
			if isBranch != nil {
				isB = isBranch(id)
			} else {
				children := childUIDs(id)
				isB = len(children) > 0
			}

			if isB {
				newExpanded[id] = true
				walk(childUIDs(id))
			}
		}
	}

	walk(roots)
	s.expandedItems.Set(newExpanded)
}

// CloseAllBranches collapses all branches in the tree.
func (s *TreeState) CloseAllBranches() {
	s.expandedItems.Set(make(map[any]bool))
}

// ExpandToNode opens all ancestors of the given node to make it visible.
// parentLookup: function to get the parent ID for a given ID, returns nil if no parent.
func (s *TreeState) ExpandToNode(id any, parentLookup func(any) any) {
	// Collect all ancestors
	ancestors := make([]any, 0)
	current := parentLookup(id)
	for current != nil {
		ancestors = append(ancestors, current)
		current = parentLookup(current)
	}

	// Expand all ancestors
	for _, ancestorID := range ancestors {
		s.Expand(ancestorID)
	}
}
