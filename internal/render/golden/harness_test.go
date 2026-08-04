package golden

import (
	"testing"

	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/foundation/divider"
	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/graphics/shape"
	"github.com/zodimo/go-compose/compose/ui/unit"
	"github.com/zodimo/go-compose/internal/modifier"
	"github.com/zodimo/go-compose/modifiers/alpha"
	"github.com/zodimo/go-compose/modifiers/background"
	"github.com/zodimo/go-compose/modifiers/border"
	"github.com/zodimo/go-compose/modifiers/clip"
	"github.com/zodimo/go-compose/modifiers/offset"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/scale"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/pkg/api"

	_ "github.com/zodimo/go-compose/internal/render/gio"
	_ "github.com/zodimo/go-compose/internal/render/software"
)

// TestBackgroundSizePadding is the core covered composition: a Box with
// Background + Size + Padding modifiers. It exercises the draw modifier
// (background fill via Backend.FillShape), a layout modifier (size), and
// a layout-only modifier (padding) in one tree.
func TestBackgroundSizePadding(t *testing.T) {
	calls := RenderComponent(t, 400, 400, 1.0, func(composer api.Composer) api.Composable {
		return box.Box(
			func(c api.Composer) api.Composer { return c },
			box.WithModifier(
				modifier.EmptyModifier.
					Then(background.Background(graphics.NewColorSrgb(255, 0, 0, 255))).
					Then(size.Size(100, 100)).
					Then(padding.Padding(16, 16, 16, 16)),
			),
		)
	})

	assertNoDiscardedCalls(t, calls)
	AssertGolden(t, "background_size_padding", calls)
}

// TestClipRectangle covers the clip modifier: a Box clipped to a rectangle
// with a background fill. Exercises Backend.Save/ClipPath/Restore around
// child layout.
func TestClipRectangle(t *testing.T) {
	calls := RenderComponent(t, 400, 400, 1.0, func(composer api.Composer) api.Composable {
		return box.Box(
			func(c api.Composer) api.Composer { return c },
			box.WithModifier(
				modifier.EmptyModifier.
					Then(clip.Clip(shape.ShapeRectangle)).
					Then(size.Size(200, 200)).
					Then(background.Background(graphics.NewColorSrgb(0, 0, 255, 255))),
			),
		)
	})

	assertNoDiscardedCalls(t, calls)
	AssertGolden(t, "clip_rectangle", calls)
}

// TestBorderRounded covers the border modifier with a rounded corner shape.
// Exercises Backend.FillShape with a stroke-representing shape on top of
// the content.
func TestBorderRounded(t *testing.T) {
	calls := RenderComponent(t, 400, 400, 1.0, func(composer api.Composer) api.Composable {
		return box.Box(
			func(c api.Composer) api.Composer { return c },
			box.WithModifier(
				modifier.EmptyModifier.
					Then(size.Size(150, 150)).
					Then(border.Border(unit.Dp(2), graphics.NewColorSrgb(255, 0, 0, 255), &shape.RoundedCornerShape{Radius: unit.Dp(8)})),
			),
		)
	})

	assertNoDiscardedCalls(t, calls)
	AssertGolden(t, "border_rounded", calls)
}

// TestAlphaModifier covers the alpha modifier: a Box with a background fill
// under an opacity layer. Exercises Backend.Save/PushOpacity/Restore.
func TestAlphaModifier(t *testing.T) {
	calls := RenderComponent(t, 400, 400, 1.0, func(composer api.Composer) api.Composable {
		return box.Box(
			func(c api.Composer) api.Composer { return c },
			box.WithModifier(
				modifier.EmptyModifier.
					Then(size.Size(100, 100)).
					Then(background.Background(graphics.NewColorSrgb(0, 255, 0, 255))).
					Then(alpha.Alpha(0.5)),
			),
		)
	})

	assertNoDiscardedCalls(t, calls)
	AssertGolden(t, "alpha_modifier", calls)
}

// TestScaleModifier covers the scale modifier: a Box scaled by 1.5.
// Exercises Backend.Save/Scale/Restore.
func TestScaleModifier(t *testing.T) {
	calls := RenderComponent(t, 400, 400, 1.0, func(composer api.Composer) api.Composable {
		return box.Box(
			func(c api.Composer) api.Composer { return c },
			box.WithModifier(
				modifier.EmptyModifier.
					Then(size.Size(100, 100)).
					Then(background.Background(graphics.NewColorSrgb(0, 0, 255, 255))).
					Then(scale.Scale(1.5)),
			),
		)
	})

	assertNoDiscardedCalls(t, calls)
	AssertGolden(t, "scale_modifier", calls)
}

// TestDividerHorizontal covers a foundation component: a horizontal divider
// with a custom thickness and color. This is a leaf widget emitting
// Backend.FillRect-style calls.
func TestDividerHorizontal(t *testing.T) {
	calls := RenderComponent(t, 400, 400, 1.0, func(composer api.Composer) api.Composable {
		return compose.Sequence(
			divider.HorizontalDivider(
				divider.WithThickness(4),
				divider.WithColor(graphics.NewColorSrgb(128, 128, 128, 255)),
			),
		)
	})

	assertNoDiscardedCalls(t, calls)
	AssertGolden(t, "divider_horizontal", calls)
}

// TestOffsetModifier covers the offset modifier: a Box offset by (10, 20) dp
// with a background fill. Exercises Backend.Translate on the software path
// via the offset layout.
func TestOffsetModifier(t *testing.T) {
	calls := RenderComponent(t, 400, 400, 1.0, func(composer api.Composer) api.Composable {
		return box.Box(
			func(c api.Composer) api.Composer { return c },
			box.WithModifier(
				modifier.EmptyModifier.
					Then(offset.Offset(10, 20)).
					Then(size.Size(50, 50)).
					Then(background.Background(graphics.NewColorSrgb(255, 0, 0, 255))),
			),
		)
	})

	assertNoDiscardedCalls(t, calls)
	AssertGolden(t, "offset_modifier", calls)
}
