package floatingactionbutton

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/foundation/material3/surface"
	"github.com/zodimo/go-compose/internal/modifier"
	"github.com/zodimo/go-compose/modifiers/clickable"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/pkg/api"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"git.sr.ht/~schnwalter/gio-mw/token"
)

// FloatingActionButton is a Material 3 Floating Action Button.
// It sits above the content and represents the primary action.
func FloatingActionButton(
	onClick func(),
	content api.Composable,
	options ...FloatingActionButtonOption,
) api.Composable {
	return func(c api.Composer) api.Composer {
		opts := DefaultFloatingActionButtonOptions()
		for _, option := range options {
			if option == nil {
				continue
			}
			option(&opts)
		}

		// Managing state for interaction (Pressed/Hovered)
		key := c.GenerateID()
		path := c.GetPath()
		statePath := fmt.Sprintf("%d/%s/fab_clickable", key, path)
		clickableState := c.State(statePath, func() any { return &widget.Clickable{} })
		fabClickable := clickableState.Get().(*widget.Clickable)

		// Determine Elevation based on state
		elevation := opts.Elevation
		if fabClickable.Pressed() {
			elevation = token.ElevationLevel3
		} else if fabClickable.Hovered() {
			elevation = token.ElevationLevel4
		}

		// Construct modifier chain
		fabModifier := opts.Modifier.Then(
			clickable.OnClick(onClick, clickable.WithClickable(fabClickable)),
		).Then(
			GetSizeModifier(opts.Size),
		)

		return SurfaceWithThemeDefaults(
			fabClickable,
			elevation,
			opts,
			fabModifier,
			box.Box(
				content,
				box.WithAlignment(layout.Center),
				box.WithModifier(size.FillMax()),
			),
		)(c)
	}
}

// SurfaceWithThemeDefaults wraps Surface.
func SurfaceWithThemeDefaults(
	fabClickable *widget.Clickable,
	elevation token.ElevationLevel,
	opts FloatingActionButtonOptions,
	fabModifier modifier.Modifier,
	content api.Composable,
) api.Composable {
	return func(c api.Composer) api.Composer {

		surfaceOpts := []surface.SurfaceOption{
			surface.WithShadowElevation(ElevationToDp(elevation)),
			surface.WithShape(opts.Shape),
			surface.WithColor(opts.ContainerColor),
			surface.WithContentColor(opts.ContentColor),
			surface.WithModifier(fabModifier),
		}

		return surface.Surface(
			content,
			surfaceOpts...,
		)(c)
	}
}

// ElevationToDp converts Material 3 elevation tokens to Dp values.
func ElevationToDp(level token.ElevationLevel) unit.Dp {
	switch level {
	case token.ElevationLevel0:
		return 0
	case token.ElevationLevel1:
		return 1
	case token.ElevationLevel2:
		return 3
	case token.ElevationLevel3:
		return 6
	case token.ElevationLevel4:
		return 8
	case token.ElevationLevel5:
		return 12
	default:
		return 6
	}
}

// GetSizeModifier returns the size modifier for the FAB based on the FabSize option.
func GetSizeModifier(fabSize FabSize) modifier.Modifier {
	switch fabSize {
	case FabSizeSmall:
		return size.Size(40, 40)
	case FabSizeLarge:
		return size.Size(96, 96)
	default:
		// Medium and default
		return size.Size(56, 56)
	}
}
