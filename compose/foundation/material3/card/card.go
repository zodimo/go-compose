package card

import (
	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/ui/graphics/shape"
	"github.com/zodimo/go-compose/internal/modifier"
	"github.com/zodimo/go-compose/modifiers/background"
	"github.com/zodimo/go-compose/modifiers/border"
	"github.com/zodimo/go-compose/modifiers/clip"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/shadow"
	"github.com/zodimo/go-compose/theme"

	"gioui.org/unit"
)

// Card corner radius (Material 3 Medium = 12dp)
var cardCornerShape = shape.RoundedCornerShape{Radius: unit.Dp(12)}

func Elevated(contents CardContentContainer, options ...CardOption) Composable {
	return cardComposable(cardElevated, contents, options...)
}

func Filled(contents CardContentContainer, options ...CardOption) Composable {
	return cardComposable(cardFilled, contents, options...)
}

func Outlined(contents CardContentContainer, options ...CardOption) Composable {
	return cardComposable(cardOutlined, contents, options...)
}

func cardComposable(kind cardKind, contents CardContentContainer, options ...CardOption) Composable {

	// Get card colors from theme
	colorRoles := theme.ColorHelper.ColorSelector()

	opts := DefaultCardOptions()
	for _, option := range options {
		if option == nil {
			continue
		}
		option(&opts)

	}

	// Build modifier chain:
	// We split user modifiers into two groups:
	// 1. Outer (Layout, Sizing): Wraps the card's style (background/border) to constrain it.
	// 2. Inner (Interaction/Clickable): Wrapped BY the card's style so overlay draws ON TOP.
	//    Also wrapped by Clip so overlay is clipped.

	userOuter, userInner := partitionModifiers(opts.Modifier)

	// 1. Outer Chain (User Layout + Card Shadow + Card Clip)
	// Order: UserOuter -> Shadow -> Clip
	cardModifier := modifier.EmptyModifier
	cardModifier = cardModifier.Then(userOuter)

	if kind == cardElevated {
		cardModifier = cardModifier.Then(shadow.Simple(unit.Dp(2), cardCornerShape))
	}

	// Always apply clip (Outlined/Filled/Elevated all use rounded corners)
	cardModifier = cardModifier.Then(clip.Clip(cardCornerShape))

	// 2. Inner Styling (Background, Border)
	// Applied inside Clip, but outside Inner User Modifiers
	styleModifier := modifier.EmptyModifier

	switch kind {
	case cardElevated:
		styleModifier = styleModifier.Then(background.Background(colorRoles.SurfaceRoles.ContainerLow, background.WithShape(cardCornerShape)))
	case cardOutlined:
		styleModifier = styleModifier.Then(background.Background(colorRoles.SurfaceRoles.Surface, background.WithShape(cardCornerShape)))
		styleModifier = styleModifier.Then(border.Border(unit.Dp(1), colorRoles.OutlineRoles.OutlineVariant, cardCornerShape))
	default: // Filled
		styleModifier = styleModifier.Then(background.Background(colorRoles.SurfaceRoles.ContainerHighest, background.WithShape(cardCornerShape)))
	}

	// 3. Assemble: Outer -> Style -> Inner
	finalModifier := cardModifier.Then(styleModifier).Then(userInner)

	// Apply user modifier AFTER clip so that clickable hover effect is clipped to rounded corners
	composables := []Composable{}

	for _, child := range contents.children {
		if child.cover {
			composables = append(composables, child.composable)
		} else {
			// add padding
			composables = append(composables, box.Box(child.composable, box.WithModifier(padding.All(16))))
		}
	}

	return column.Column(compose.Sequence(composables...), column.WithModifier(finalModifier))
}

// partitionModifiers splits the modifier chain into outer (layout/constraints)
// and inner (interaction/overlay) segments to ensure correct layering.
func partitionModifiers(m Modifier) (outer, inner Modifier) {
	outer = modifier.EmptyModifier
	inner = modifier.EmptyModifier

	// FoldOut walks Head -> Tail (Left -> Right)
	m.AsChain().FoldOut(nil, func(_ interface{}, elt Modifier) interface{} {
		if inspectable, ok := elt.(modifier.InspectableModifier); ok {
			if inspectable.InspectorInfo().Name == "clickable" {
				inner = inner.Then(elt)
				return nil
			}
		}
		outer = outer.Then(elt)
		return nil
	})

	return outer, inner
}
