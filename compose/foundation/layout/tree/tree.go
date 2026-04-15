package tree

import (
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	"github.com/zodimo/go-compose/compose/foundation/lazy"
	"github.com/zodimo/go-compose/compose/material3/icon"
	"github.com/zodimo/go-compose/compose/ui/unit"
	"github.com/zodimo/go-compose/modifiers/clickable"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/modifiers/weight"

	"github.com/zodimo/go-compose/pkg/api"
)

// TreeScope defines the operations available within a Tree.
type TreeScope interface {
	// Node adds a leaf node to the tree.
	Node(key any, content api.Composable)

	// Branch adds a collapsible branch node to the tree.
	// header is the content displayed for the branch itself.
	// children is a function that defines the children of this branch.
	Branch(key any, header api.Composable, children func(TreeScope))
}

// Tree creates a tree component that efficiently renders hierarchical data using LazyColumn.
// It supports expanding and collapsing branches.
func Tree(
	state *TreeState,
	content func(TreeScope),
	options ...TreeOption,
) api.Composable {
	opts := DefaultTreeOptions()
	for _, opt := range options {
		opt(&opts)
	}

	// Convert TreeOption modifier to lazy option if present
	lazyOpts := make([]lazy.LazyListOption, 0)
	if opts.Modifier != nil {
		lazyOpts = append(lazyOpts, lazy.WithModifier(opts.Modifier))
	}

	return lazy.LazyColumn(
		func(scope lazy.LazyListScope) {
			tScope := &treeScopeImpl{
				listScope: scope,
				state:     state,
				depth:     0,
				options:   &opts,
			}
			content(tScope)
		},
		lazyOpts...,
	)
}

type treeScopeImpl struct {
	listScope lazy.LazyListScope
	state     *TreeState
	depth     int
	options   *TreeOptions
}

func (s *treeScopeImpl) Node(key any, content api.Composable) {
	indentSize := s.options.IndentSize
	opts := s.options
	state := s.state

	s.listScope.Item(key, func(c api.Composer) api.Composer {
		return row.Row(
			c.Sequence(
				// Indentation
				spacer.Width(s.depth*indentSize),
				// Spacer for the expander icon alignment
				// spacer.Width(indentSize),
				// Node content
				content,
			),
			row.WithAlignment(row.Middle),
			row.WithModifier(
				size.FillMaxWidth().
					Then(clickable.OnClick(func() {
						SelectNodeWithCallback(state, key, opts)
					}).
						Then(padding.All(4)),
					),
			),
		)(c)
	})
}

func (s *treeScopeImpl) Branch(key any, header api.Composable, children func(TreeScope)) {
	isExpanded := s.state.IsExpanded(key)
	indentSize := s.options.IndentSize

	// Capture state and options for closure
	state := s.state
	opts := s.options

	rightArrowIcon := icon.SymbolArrowRight
	downArrowIcon := icon.SymbolArrowDropDown

	// Branch Header
	s.listScope.Item(key, func(c api.Composer) api.Composer {
		return row.Row(
			c.Sequence(
				row.Row(
					c.Sequence(
						// Indentation
						spacer.Width(s.depth*indentSize),
						// Expander Icon
						// Toggle Button - only toggles expand/collapse

						c.If(
							isExpanded,
							icon.Icon(
								downArrowIcon,
								icon.WithSize(unit.Dp(s.options.IconSize)),
							),

							icon.Icon(
								rightArrowIcon,
								icon.WithSize(unit.Dp(s.options.IconSize)),
							),
						),
					),
					row.WithAlignment(row.Middle),
					row.WithModifier(
						clickable.OnClick(func() {
							// toggle branch node
							toggleBranchWithCallback(state, key, opts)

						}).
							Then(padding.All(4)),
					),
				),
				row.Row(
					c.Sequence(

						// Header Content
						header,
					),
					row.WithAlignment(row.Middle),
					row.WithModifier(
						weight.Weight(1).
							Then(clickable.OnClick(func() {
								// Select the branch node
								SelectNodeWithCallback(state, key, opts)
							}).
								Then(padding.All(4)),
							),
					),
				),
			),
			row.WithModifier(
				size.FillMaxWidth(),
			),
		)(c)
	})

	// Children
	if isExpanded {
		childScope := &treeScopeImpl{
			listScope: s.listScope,
			state:     s.state,
			depth:     s.depth + 1,
			options:   s.options,
		}
		children(childScope)
	}
}

// toggleBranchWithCallback toggles the branch and invokes appropriate callbacks.
func toggleBranchWithCallback(state *TreeState, key any, opts *TreeOptions) {
	wasExpanded := state.IsExpanded(key)
	state.Toggle(key)

	if wasExpanded {
		// Branch was closed
		if opts.OnBranchClosed != nil {
			opts.OnBranchClosed(key)
		}
	} else {
		// Branch was opened
		if opts.OnBranchOpened != nil {
			opts.OnBranchOpened(key)
		}
	}
}

// TreeFromData creates a Tree from a data structure, similar to fyne's data-driven Tree.
// roots: List of root IDs.
// childUIDs: Function to get children IDs for a given ID.
// isBranch: Function to strict check if a node is a branch (optional, if nil, checks if children > 0).
// createNode: Composable factory for a node ID.
func TreeFromData(
	state *TreeState,
	roots []any,
	childUIDs func(any) []any,
	isBranch func(any) bool,
	createNode func(any) api.Composable,
	options ...TreeOption,
) api.Composable {

	return Tree(
		state,
		func(scope TreeScope) {
			var traverse func(s TreeScope, ids []any)
			traverse = func(s TreeScope, ids []any) {
				for _, id := range ids {
					isB := false
					if isBranch != nil {
						isB = isBranch(id)
					} else {
						children := childUIDs(id)
						isB = len(children) > 0
					}

					if isB {
						s.Branch(
							id,
							createNode(id),
							func(innerS TreeScope) {
								traverse(innerS, childUIDs(id))
							},
						)
					} else {
						s.Node(id, createNode(id))
					}
				}
			}

			traverse(scope, roots)
		},
		options...,
	)
}

// SelectNodeWithCallback selects a node and invokes the appropriate callbacks.
// This is a helper function for use in node click handlers.
func SelectNodeWithCallback(state *TreeState, id any, opts *TreeOptions) {
	// Get previously selected items for callback
	previouslySelected := state.GetSelectedItems()

	// Select the new item (single selection mode)
	state.Select(id, true)

	// Call OnUnselected for previously selected items
	if opts != nil && opts.OnUnselected != nil {
		for _, prevID := range previouslySelected {
			if prevID != id {
				opts.OnUnselected(prevID)
			}
		}
	}

	// Call OnSelected for newly selected item
	if opts != nil && opts.OnSelected != nil {
		opts.OnSelected(id)
	}
}
