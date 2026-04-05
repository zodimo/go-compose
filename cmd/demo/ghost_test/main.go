package main

import (
	"image"
	"image/color"
	"log"
	"os"

	"github.com/zodimo/go-compose/compose"
	fImage "github.com/zodimo/go-compose/compose/foundation/image"
	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	"github.com/zodimo/go-compose/compose/foundation/lazy"
	"github.com/zodimo/go-compose/compose/material3/card"
	m3text "github.com/zodimo/go-compose/compose/material3/text"
	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	uilayout "github.com/zodimo/go-compose/compose/ui/layout"
	"github.com/zodimo/go-compose/modifiers/background"
	"github.com/zodimo/go-compose/modifiers/padding"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/runtime"
	"github.com/zodimo/go-compose/store"
	"github.com/zodimo/go-compose/theme"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

type TestCase int

const (
	TestLazyBoxes TestCase = iota
	TestLazyImages
	TestLazyCardsNoImages
	TestLazyCardsWithImages
	TestStaticColumnImages
)

var currentTest = TestLazyCardsWithImages

func testName(tc TestCase) string {
	names := []string{
		"Lazy Boxes",
		"Lazy Images",
		"Lazy Cards (No Images)",
		"Lazy Cards + Images (GHOST)",
		"Static Column + Images",
	}
	if int(tc) < len(names) {
		return names[tc]
	}
	return "Unknown"
}

func main() {
	go func() {
		w := new(app.Window)
		w.Option(app.Title("Ghost Test - " + testName(currentTest)))
		w.Option(app.Size(unit.Dp(720), unit.Dp(800)))

		if err := Run(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	app.Main()
}

func Run(window *app.Window) error {
	enLocale := system.Locale{Language: "en", Direction: system.LTR}
	var ops op.Ops
	store := store.NewPersistentState()
	store.Subscribe(func() { window.Invalidate() })
	runtime := runtime.NewRuntime()
	themeManager := theme.GetThemeManager()

	for {
		switch frameEvent := window.Event().(type) {
		case app.DestroyEvent:
			return frameEvent.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, frameEvent)
			gtx.Locale = enLocale
			gtx = themeManager.Material3ThemeInit(gtx)

			composer := compose.NewComposer(api.ComposerWithStore(store))
			callOp := runtime.Run(gtx, composer, UI())
			callOp.Add(gtx.Ops)
			frameEvent.Frame(gtx.Ops)
		}
	}
}

func UI() api.Composable {
	return func(c api.Composer) api.Composer {
		switch currentTest {
		case TestLazyBoxes:
			c = TestLazyBoxesUI()(c)
		case TestLazyImages:
			c = TestLazyImagesUI()(c)
		case TestLazyCardsNoImages:
			c = TestLazyCardsNoImagesUI()(c)
		case TestLazyCardsWithImages:
			c = TestLazyCardsWithImagesUI()(c)
		case TestStaticColumnImages:
			c = TestStaticColumnImagesUI()(c)
		}
		c.Modifier(func(m ui.Modifier) ui.Modifier {
			return m.Then(padding.All(24))
		})
		return c
	}
}

func TestLazyBoxesUI() api.Composable {
	return func(c api.Composer) api.Composer {
		return lazy.LazyColumn(
			func(scope lazy.LazyListScope) {
				for i := 0; i < 20; i++ {
					idx := i
					key := "box_" + string(rune('a'+idx%26))
					scope.Item(key, box.Box(
						m3text.BodyMedium("Item "+string(rune('A'+idx%26))),
						box.WithModifier(
							size.Width(300, size.SizeRequired()).
								Then(size.Height(100)).
								Then(background.Background(graphics.NewColorSrgb(100+idx*5, 150-idx*3, 200, 255))).
								Then(padding.All(16)),
						),
					))
				}
			},
			lazy.WithModifier(padding.All(24).Then(size.FillMax())),
		)(c)
	}
}

func TestLazyImagesUI() api.Composable {
	return func(c api.Composer) api.Composer {
		return lazy.LazyColumn(
			func(scope lazy.LazyListScope) {
				scope.Item("title", m3text.HeadlineMedium("Lazy Images Only"))

				for i := 0; i < 10; i++ {
					idx := i
					key := "img_" + string(rune('a'+idx%26))
					scope.Item(key, column.Column(
						c.Sequence(
							m3text.TitleMedium("Image "+string(rune('A'+idx%26))),
							spacer.Height(8),
							fImage.Image(
								CreateTestImage(idx),
								fImage.WithContentScale(uilayout.ContentScaleCrop),
								fImage.WithModifier(size.Height(120)),
							),
						),
					))
				}
			},
			lazy.WithModifier(size.FillMax()),
		)(c)
	}
}

func TestLazyCardsNoImagesUI() api.Composable {
	return func(c api.Composer) api.Composer {
		return lazy.LazyColumn(
			func(scope lazy.LazyListScope) {
				scope.Item("title", m3text.HeadlineMedium("Lazy Cards (No Images)"))

				for i := 0; i < 15; i++ {
					idx := i
					key := "card_" + string(rune('a'+idx%26))
					scope.Item(key, card.Elevated(
						card.CardContents(
							card.Content(
								column.Column(
									c.Sequence(
										m3text.TitleMedium("Card "+string(rune('A'+idx%26))),
										spacer.Height(8),
										m3text.BodyMedium("This card has no image content."),
									),
								),
							),
						),
						card.WithModifier(size.Width(340, size.SizeRequired())),
					))
				}
			},
			lazy.WithModifier(size.FillMax()),
		)(c)
	}
}

func TestLazyCardsWithImagesUI() api.Composable {
	return func(c api.Composer) api.Composer {
		return lazy.LazyColumn(
			func(scope lazy.LazyListScope) {
				scope.Item("title", m3text.HeadlineMedium("Lazy Cards + Images (GHOST TEST)"))
				scope.Item("desc", m3text.BodySmall("Scroll to bottom and check for ghost images"))

				for i := 0; i < 10; i++ {
					idx := i
					key := "card_img_" + string(rune('a'+idx%26))
					scope.Item(key, card.Elevated(
						card.CardContents(
							card.ContentCover(
								fImage.Image(
									CreateTestImage(idx),
									fImage.WithContentScale(uilayout.ContentScaleCrop),
									fImage.WithModifier(size.Height(120)),
								),
							),
							card.Content(
								column.Column(
									c.Sequence(
										m3text.TitleMedium("Card with Image "+string(rune('A'+idx%26))),
										spacer.Height(4),
										m3text.BodySmall("This card has an image in ContentCover"),
									),
								),
							),
						),
						card.WithModifier(size.Width(340, size.SizeRequired())),
					))
				}
			},
			lazy.WithModifier(size.FillMax()),
		)(c)
	}
}

func TestStaticColumnImagesUI() api.Composable {
	return func(c api.Composer) api.Composer {
		items := []api.Composable{
			m3text.HeadlineMedium("Static Column + Images"),
			m3text.BodySmall("No lazy loading - all items composed upfront"),
		}

		for i := 0; i < 10; i++ {
			idx := i
			items = append(items,
				card.Elevated(
					card.CardContents(
						card.ContentCover(
							fImage.Image(
								CreateTestImage(idx),
								fImage.WithContentScale(uilayout.ContentScaleCrop),
								fImage.WithModifier(size.Height(120)),
							),
						),
						card.Content(
							column.Column(
								c.Sequence(
									m3text.TitleMedium("Card with Image "+string(rune('A'+idx%26))),
									spacer.Height(4),
									m3text.BodySmall("Static composition"),
								),
							),
						),
					),
					card.WithModifier(size.Width(340, size.SizeRequired())),
				),
			)
		}

		return column.Column(
			c.Sequence(items...),
			column.WithModifier(padding.All(16).Then(size.FillMax())),
		)(c)
	}
}

func CreateTestImage(index int) graphics.ImageResource {
	width, height := 340, 120
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	baseR := uint8(80 + index*15)
	baseG := uint8(60 + index*10)
	baseB := uint8(140 + index*8)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r := baseR + uint8(x*30/width)
			g := baseG + uint8(y*30/height)
			b := baseB
			img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
		}
	}

	return graphics.ImageResource{
		ImageOp: paint.NewImageOp(img),
	}
}
