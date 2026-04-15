package tree

import (
	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/internal/modifier"
)

// TreeOption configures a Tree component.
type TreeOption func(*TreeOptions)

// TreeOptions holds configuration for a Tree component.
type TreeOptions struct {
	// Modifier applies styling to the tree container.
	Modifier ui.Modifier

	// IndentSize is the horizontal indentation per depth level (default: 24).
	IndentSize int

	// HideSeparators hides the visual separators between tree nodes.
	HideSeparators bool

	IconSize int

	// Callbacks (inspired by Fyne)

	// OnSelected is called when a node is selected.
	OnSelected func(id any)

	// OnUnselected is called when a node is unselected.
	OnUnselected func(id any)

	// OnBranchOpened is called when a branch is expanded.
	OnBranchOpened func(id any)

	// OnBranchClosed is called when a branch is collapsed.
	OnBranchClosed func(id any)
}

// DefaultTreeOptions returns the default TreeOptions.
func DefaultTreeOptions() TreeOptions {
	return TreeOptions{
		Modifier:       modifier.EmptyModifier,
		IndentSize:     24,
		IconSize:       16,
		HideSeparators: false,
	}
}

// WithModifier sets a modifier for the tree container.
func WithModifier(m ui.Modifier) TreeOption {
	return func(o *TreeOptions) {
		o.Modifier = o.Modifier.Then(m)
	}
}

// WithIndentSize sets the indentation size per depth level.
func WithIndentSize(size int) TreeOption {
	return func(o *TreeOptions) {
		o.IndentSize = size
	}
}

func WithIconSize(size int) TreeOption {
	return func(o *TreeOptions) {
		o.IconSize = size
	}
}

// WithHideSeparators hides separators between tree nodes.
func WithHideSeparators(hide bool) TreeOption {
	return func(o *TreeOptions) {
		o.HideSeparators = hide
	}
}

// WithOnSelected sets the callback for when a node is selected.
func WithOnSelected(callback func(id any)) TreeOption {
	return func(o *TreeOptions) {
		o.OnSelected = callback
	}
}

// WithOnUnselected sets the callback for when a node is unselected.
func WithOnUnselected(callback func(id any)) TreeOption {
	return func(o *TreeOptions) {
		o.OnUnselected = callback
	}
}

// WithOnBranchOpened sets the callback for when a branch is expanded.
func WithOnBranchOpened(callback func(id any)) TreeOption {
	return func(o *TreeOptions) {
		o.OnBranchOpened = callback
	}
}

// WithOnBranchClosed sets the callback for when a branch is collapsed.
func WithOnBranchClosed(callback func(id any)) TreeOption {
	return func(o *TreeOptions) {
		o.OnBranchClosed = callback
	}
}
