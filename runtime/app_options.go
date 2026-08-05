package runtime

import (
	"context"

	"github.com/zodimo/go-compose/compose/ui/unit"
	"github.com/zodimo/go-compose/pkg/api"
)

// AppConfig carries configuration options for runtime.App.
type AppConfig struct {
	Title      string
	Size       unit.DpSize
	Content    api.Composable
	RootCtx    context.Context
	QuitOnLast bool
}

// AppOption configures an App.
type AppOption func(*AppConfig)

// DefaultAppConfig returns the default AppConfig.
func DefaultAppConfig() AppConfig {
	return AppConfig{
		QuitOnLast: true,
		RootCtx:    context.Background(),
	}
}

// Title sets the window title for the default single window.
func Title(t string) AppOption {
	return func(c *AppConfig) {
		c.Title = t
	}
}

// Size sets the initial window size for the default single window.
func Size(size unit.DpSize) AppOption {
	return func(c *AppConfig) {
		c.Size = size
	}
}

// Content sets the root composable UI for the default single window.
func Content(ui api.Composable) AppOption {
	return func(c *AppConfig) {
		c.Content = ui
	}
}

// WithRootContext sets the root context for the application.
// When cancelled, App.Run() exits with the context error.
func WithRootContext(ctx context.Context) AppOption {
	return func(c *AppConfig) {
		c.RootCtx = ctx
	}
}

// QuitWhenLastWindowCloses configures whether closing the last window exits the application.
// Defaults to true.
func QuitWhenLastWindowCloses(b bool) AppOption {
	return func(c *AppConfig) {
		c.QuitOnLast = b
	}
}

// WindowOptions carries per-window settings for AddWindow.
type WindowOptions struct {
	Title   string
	Size    unit.DpSize
	Content api.Composable
}
