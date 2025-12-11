package dialog

import "git.sr.ht/~schnwalter/gio-mw/wdk"

type DialogOptions struct {
	Modifier     Modifier
	Title        string
	Text         string
	DismissLabel string
	Icon         wdk.IconWidget
	OnDismiss    func()
}

type DialogOption func(*DialogOptions)

func DefaultDialogOptions() DialogOptions {
	return DialogOptions{
		Modifier: EmptyModifier,
	}
}

func WithModifier(m Modifier) DialogOption {
	return func(o *DialogOptions) {
		o.Modifier = m
	}
}

func WithTitle(title string) DialogOption {
	return func(o *DialogOptions) {
		o.Title = title
	}
}

func WithText(text string) DialogOption {
	return func(o *DialogOptions) {
		o.Text = text
	}
}

func WithIcon(icon wdk.IconWidget) DialogOption {
	return func(o *DialogOptions) {
		o.Icon = icon
	}
}

func WithDismissButton(label string, onDismiss func()) DialogOption {
	return func(o *DialogOptions) {
		o.DismissLabel = label
		o.OnDismiss = onDismiss
	}
}
