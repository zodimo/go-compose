package tree

import (
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	"github.com/zodimo/go-compose/compose/foundation/lazy"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/modifiers/background"
	"github.com/zodimo/go-compose/modifiers/clickable"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/modifiers/weight"
	"github.com/zodimo/go-compose/pkg/api"
)

type TreeScope interface {
	// Node adds a leaf node to the tree.
	Node(key any, content api.Composable)

	// Branch adds a collapsible branch node to the tree.
	// header is the content displayed for the branch itself.
	// children is a function that defines the children of this branch.
	Branch(key any, header api.Composable, children func(TreeScope))
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

	isSelected := state.IsSelected(key)
	rowBackground := graphics.ColorUnspecified
	if isSelected {
		rowBackground = graphics.Selected(s.options.SelectedNodeOnColor)
	}

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
			row.WithModifier(
				background.Background(rowBackground),
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
			row.WithAlignment(row.Middle),
		)(c)
	})
}

func (s *treeScopeImpl) Branch(key any, header api.Composable, children func(TreeScope)) {
	isExpanded := s.state.IsExpanded(key)
	indentSize := s.options.IndentSize

	// Capture state and options for closure
	state := s.state
	opts := s.options

	isSelected := state.IsSelected(key)
	rowBackground := graphics.ColorUnspecified
	if isSelected {
		rowBackground = graphics.Selected(s.options.SelectedBranchOnColor)
	}

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
							opts.BranchIcons.OpenIcon,
							opts.BranchIcons.ClosedIcon,
						),
					),
					row.WithModifier(
						background.Background(rowBackground),
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
					row.WithModifier(
						background.Background(rowBackground),
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
			row.WithAlignment(row.Middle),
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
