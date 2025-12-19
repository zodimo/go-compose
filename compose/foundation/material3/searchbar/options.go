package searchbar

import (
	"git.sr.ht/~schnwalter/gio-mw/widget/search"
)

type SearchBarOptions struct {
	Modifier       Modifier
	Enabled        bool
	SupportingText string
	LeadingIcon    search.Icon
	TrailingIcon   search.Icon
	OnSubmit       func()
}

func DefaultSearchBarOptions() SearchBarOptions {
	return SearchBarOptions{
		Modifier: EmptyModifier,
		Enabled:  true,
	}
}

type SearchBarOption func(*SearchBarOptions)

func WithModifier(m Modifier) SearchBarOption {
	return func(o *SearchBarOptions) {
		o.Modifier = m
	}
}

func WithEnabled(enabled bool) SearchBarOption {
	return func(o *SearchBarOptions) {
		o.Enabled = enabled
	}
}

func WithSupportingText(text string) SearchBarOption {
	return func(o *SearchBarOptions) {
		o.SupportingText = text
	}
}

func WithLeadingIcon(icon search.Icon) SearchBarOption {
	return func(o *SearchBarOptions) {
		o.LeadingIcon = icon
	}
}

func WithTrailingIcon(icon search.Icon) SearchBarOption {
	return func(o *SearchBarOptions) {
		o.TrailingIcon = icon
	}
}

func WithOnSubmit(f func()) SearchBarOption {
	return func(o *SearchBarOptions) {
		o.OnSubmit = f
	}
}
