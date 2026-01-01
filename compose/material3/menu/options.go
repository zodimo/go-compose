package menu

import (
	"github.com/zodimo/go-compose/internal/modifier"
	"github.com/zodimo/go-compose/pkg/api"

	"github.com/zodimo/go-compose/compose/ui/unit"
)

type DropdownMenuOptions struct {
	Modifier modifier.Modifier
	OffsetX  unit.Dp
	OffsetY  unit.Dp
}

func DefaultDropdownMenuOptions() DropdownMenuOptions {
	return DropdownMenuOptions{
		Modifier: modifier.EmptyModifier,
		OffsetX:  0,
		OffsetY:  0,
	}
}

type DropdownMenuOption func(*DropdownMenuOptions)

func WithModifier(m modifier.Modifier) DropdownMenuOption {
	return func(opts *DropdownMenuOptions) {
		opts.Modifier = m
	}
}

func WithOffset(x, y unit.Dp) DropdownMenuOption {
	return func(opts *DropdownMenuOptions) {
		opts.OffsetX = x
		opts.OffsetY = y
	}
}

// DropdownMenuItemOptions

type DropdownMenuItemOptions struct {
	Modifier     modifier.Modifier
	LeadingIcon  api.Composable
	TrailingIcon api.Composable
	Enabled      bool
}

func DefaultDropdownMenuItemOptions() DropdownMenuItemOptions {
	return DropdownMenuItemOptions{
		Modifier:     modifier.EmptyModifier,
		LeadingIcon:  nil,
		TrailingIcon: nil,
		Enabled:      true,
	}
}

type DropdownMenuItemOption func(*DropdownMenuItemOptions)

func WithMenuItemModifier(mod modifier.Modifier) DropdownMenuItemOption {
	return func(opts *DropdownMenuItemOptions) {
		opts.Modifier = mod
	}
}

func WithLeadingIcon(icon api.Composable) DropdownMenuItemOption {
	return func(opts *DropdownMenuItemOptions) {
		opts.LeadingIcon = icon
	}
}

func WithTrailingIcon(icon api.Composable) DropdownMenuItemOption {
	return func(opts *DropdownMenuItemOptions) {
		opts.TrailingIcon = icon
	}
}

func WithEnabled(enabled bool) DropdownMenuItemOption {
	return func(opts *DropdownMenuItemOptions) {
		opts.Enabled = enabled
	}
}
