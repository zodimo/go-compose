package bottomsheet

import (
	"time"

	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/material3/surface"
	"github.com/zodimo/go-compose/compose/ui/graphics/shape"
	"github.com/zodimo/go-compose/internal/animation"
	"github.com/zodimo/go-compose/internal/modifier"
	animMod "github.com/zodimo/go-compose/modifiers/animation"
	"github.com/zodimo/go-compose/modifiers/clickable"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/theme"

	"gioui.org/unit"
	"git.sr.ht/~schnwalter/gio-mw/token"
)

type Modifier = modifier.Modifier
type Composer = api.Composer

func ModalBottomSheet(
	sheetContent Composable,
	content Composable,
	options ...ModalBottomSheetOption,
) Composable {
	opts := DefaultModalBottomSheetOptions()
	for _, opt := range options {
		opt(&opts)
	}

	return func(c Composer) Composer {
		tm := theme.GetThemeManager()
		m3 := tm.GetMaterial3Theme()

		// defaults
		containerColor := opts.ContainerColor
		if containerColor == (token.MatColor{}) {
			containerColor = m3.Scheme.SurfaceContainerLow
		}

		scrimColor := opts.ScrimColor
		if scrimColor == (token.MatColor{}) {
			scrimColor = m3.Scheme.Scrim
		}

		sheetShape := opts.Shape
		if sheetShape == (token.CornerShape{}) {
			// Default M3 sheet shape: Top corners rounded 28.dp
			sheetShape = token.CornerShape{
				Kind: token.CornerKindRound,
				Size: unit.Dp(28),
			}
		}

		// Animation state
		// If SheetState provided, use it. Otherwise create internal one.
		// Usually we want persistent state to control close.

		var anim *animation.VisibilityAnimation
		if opts.SheetState != nil {
			anim = opts.SheetState.visibleAnim
		} else {
			// Internal state sync with IsOpen
			state := c.State(c.GenerateID().String()+"/anim", func() any {
				return &animation.VisibilityAnimation{
					Duration: time.Millisecond * 300,
					State:    animation.Invisible,
				}
			}).Get().(*animation.VisibilityAnimation)
			anim = state
		}

		// Sync animation state with IsOpen if no external SheetState manipulation is expected
		// OR if IsOpen is the driver.
		// Simplest: Check IsOpen.
		if opts.IsOpen {
			anim.Appear(time.Now())
		} else {
			// Only hide if not managed by SheetState manually?
			// If opts.SheetState is set, the caller calls Show/Hide.
			// But IsOpen is also an option.
			// Let's assume IsOpen drives it if it changes.
			// But since we are declarative, IsOpen IS the truth.
			anim.Disappear(time.Now())
		}

		baseScrim := scrimColor.SetOpacity(token.OpacityLevel8).AsNRGBA()

		return box.Box(
			func(c Composer) Composer {
				// 1. Main Content (Background)
				content(c)

				// Determine if we should render overlay elements
				shouldRender := opts.IsOpen || anim.Visible()

				if shouldRender {
					// 2. Scrim (Animated Opacity)
					box.Box(
						func(c Composer) Composer { return c },
						box.WithModifier(
							modifier.EmptyModifier.
								Then(size.FillMax()).
								Then(animMod.AnimatedBackground(anim, baseScrim, shape.ShapeRectangle)).
								Then(clickable.OnClick(func() {
									if opts.OnDismissRequest != nil {
										opts.OnDismissRequest()
									}
								})),
						),
					)(c)

					// 3. Bottom Sheet Surface
					// Needs to align bottom.
					box.Box(
						func(c Composer) Composer {
							return surface.Surface(
								func(c Composer) Composer {
									return column.Column(
										func(c Composer) Composer {
											// 1. Drag Handle
											if opts.DragHandle != nil {
												opts.DragHandle(c)
											} else {
												// Default Drag Handle
												box.Box(
													func(c Composer) Composer { return c },
													box.WithModifier(
														modifier.EmptyModifier.
															Then(size.Width(32)).
															Then(size.Height(4)).
															Then(animMod.AnimatedBackground(anim, m3.Scheme.SurfaceVariant.OnColor.SetOpacity(token.OpacityLevel(0.4)).AsNRGBA(), shape.RoundedCornerShape{Radius: unit.Dp(2)})),
													),
												)(c)
											}

											// Spacing
											box.Box(func(c Composer) Composer { return c }, box.WithModifier(size.Height(22)))(c)

											// 2. Content
											sheetContent(c)
											return c
										},
										// Ensure column centers the drag handle
										column.WithAlignment(column.Middle),
									)(c)
								},
								surface.WithColor(theme.ColorHelper.SpecificColor(containerColor.AsNRGBA())),
								// Use the custom shape. Surface expects Shape type, so we might need adapter or use corner radius if uniform?
								// M3 sheet usually has only Top corners rounded.
								// Our Surface/Shape implementation might need check.
								// For now, let's use a RoundedCornerShape adapter or just pass specific radius if uniform.
								// But `token.CornerShape` is usually single corner.
								// Let's use `shape.RoundedCornerShape{Radius: sheetShape.Size}` assuming uniform for now,
								// or better, make a Custom shape if possible.
								// The `drawer.go` used `unit.Dp(16)`.
								surface.WithShape(shape.RoundedCornerShape{Radius: sheetShape.Size}),
								surface.WithModifier(
									modifier.EmptyModifier.
										Then(animMod.AnimatedHeight(anim, 0)). // Animate height from 0 (or slide up)
										Then(size.FillMaxWidth()),
								),
							)(c)
						},
						box.WithModifier(
							modifier.EmptyModifier.Then(size.FillMax()),
						),
						box.WithAlignment(box.S), // Align the surface to the bottom
					)(c)
				}
				return c
			},
		)(c)
	}
}
