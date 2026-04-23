package tree

import (
	"github.com/zodimo/go-compose/compose/material3/icon"
	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/compose/ui/text"
	"github.com/zodimo/go-compose/internal/modifier"
	"github.com/zodimo/go-compose/pkg/api"
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

	//will be in SP
	TextStyle *text.TextStyle

	// Callbacks (inspired by Fyne)

	// OnSelected is called when a node is selected.
	OnSelected func(id any)

	// OnUnselected is called when a node is unselected.
	OnUnselected func(id any)

	// OnBranchOpened is called when a branch is expanded.
	OnBranchOpened func(id any)

	// OnBranchClosed is called when a branch is collapsed.
	OnBranchClosed func(id any)

	BranchIcons *TreeBranchIcons
}

// DefaultTreeOptions returns the default TreeOptions.
func DefaultTreeOptions() TreeOptions {
	return TreeOptions{
		Modifier:       modifier.EmptyModifier,
		IndentSize:     24,
		TextStyle:      text.TextStyleUnspecified,
		HideSeparators: false,
		BranchIcons: newTreeBranchIcons(
			icon.Icon(icon.SymbolIndeterminateCheckBox),
			icon.Icon(icon.SymbolAddBox),
		),
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

func WithTextStyle(style *text.TextStyle) TreeOption {
	return func(o *TreeOptions) {
		o.TextStyle = style
	}
}

func WithTextStyleOption(textStyleOption text.TextStyleOption) TreeOption {
	return func(o *TreeOptions) {
		o.TextStyle = text.CopyTextStyle(o.TextStyle, textStyleOption)
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
func WithBranchIcons(icons *TreeBranchIcons) TreeOption {
	return func(o *TreeOptions) {
		o.BranchIcons = icons
	}
}

func WithBranchOpenIcon(open api.Composable) TreeOption {
	return func(o *TreeOptions) {
		o.BranchIcons.OpenIcon = open
	}
}

func WithBranchCloseIcon(closed api.Composable) TreeOption {
	return func(o *TreeOptions) {
		o.BranchIcons.ClosedIcon = closed
	}
}
