package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"

	"gioui.org/x/explorer"
	fImage "github.com/zodimo/go-compose/compose/foundation/image"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	"github.com/zodimo/go-compose/compose/material3/button"
	"github.com/zodimo/go-compose/compose/material3/surface"
	"github.com/zodimo/go-compose/compose/material3/text"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/unit"

	"github.com/zodimo/go-compose/compose/ui/graphics/shape"
	uilayout "github.com/zodimo/go-compose/compose/ui/layout"
	"github.com/zodimo/go-compose/modifiers/border"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/pkg/x/fileexplorer"
	"github.com/zodimo/go-compose/state"
	"github.com/zodimo/go-maybe"
)

type ImageResult struct {
	Error         error
	Format        string
	Image         image.Image
	ImageResource fImage.ImageResource
}

func UI() api.Composable {
	return func(c api.Composer) api.Composer {
		maybeImageResult := state.MustRemember(c, "maybeImage", func() maybe.Maybe[ImageResult] {
			return maybe.None[ImageResult]()
		})
		maybeSaveResult := state.MustRemember(c, "maybeSaveResult", func() maybe.Maybe[error] {
			return maybe.None[error]()
		})

		isSubscribedImageResult := state.MustRemember(c, "isSubscribedImageResult", func() bool {
			return false
		})

		canSaveFile := state.MustRemember(c, "canSaveFile", func() bool {
			return false
		})

		if !isSubscribedImageResult.Get() {
			maybeImageResult.Subscribe(func() {
				fmt.Printf("Open Image Result: %s\n", maybeImageResult.Get())
				canSaveFile.Set(maybe.Match(
					maybeImageResult.Get(),
					func(imageResult ImageResult) bool {
						return imageResult.Error == nil
					},
					func() bool {
						return false
					},
				))
			})
			isSubscribedImageResult.Set(true)
		}

		isSubscribedSaveResult := state.MustRemember(c, "isSubscribedSaveResult", func() bool {
			return false
		})

		if !isSubscribedSaveResult.Get() {
			maybeSaveResult.Subscribe(func() {
				fmt.Printf("Save Image Result: %s\n", maybeSaveResult.Get())
			})
			isSubscribedSaveResult.Set(true)
		}

		onOpenFile, launchedOpenFile := fileexplorer.RememberExplorer(c, func(expl *explorer.Explorer) {

			file, err := expl.ChooseFile("png", "jpeg", "jpg")
			if err != nil {
				err = fmt.Errorf("failed opening image file: %w", err)
				maybeImageResult.Set(maybe.Some(ImageResult{Error: err}))
				return
			}
			defer file.Close()

			imageBuffer := bytes.NewBuffer(nil)
			_, err = io.Copy(imageBuffer, file)
			reader := bytes.NewReader(imageBuffer.Bytes())
			if err != nil {
				err = fmt.Errorf("failed copying image data: %w", err)
				maybeImageResult.Set(maybe.Some(ImageResult{Error: err}))
				return
			}

			imgData, format, err := image.Decode(reader)
			if err != nil {
				err = fmt.Errorf("failed decoding image data: %w", err)
				maybeImageResult.Set(maybe.Some(ImageResult{Error: err}))
				return
			}
			_, err = reader.Seek(0, io.SeekStart)

			if err != nil {
				err = fmt.Errorf("failed rewinding image data: %w", err)
				maybeImageResult.Set(maybe.Some(ImageResult{Error: err}))
				return
			}

			maybeImageResult.Set(maybe.Some(ImageResult{
				Image:         imgData,
				Format:        format,
				ImageResource: graphics.NewResourceFromImageFile(reader),
			}))
		})

		onSaveFile, launchedSaveFile := fileexplorer.RememberExplorer(c, func(expl *explorer.Explorer) {

			if maybeImageResult.Get().IsNone() {
				fmt.Println("no image loaded, cannot save")
				return
			}

			if maybeImageResult.Get().IsSome() {
				imageResult := maybeImageResult.Get().UnwrapUnsafe()
				if imageResult.Error != nil {
					fmt.Printf("error: %s\n", imageResult.Error)
					return
				}

				extension := "jpg"
				switch imageResult.Format {
				case "png":
					extension = "png"
				}
				file, err := expl.CreateFile("file." + extension)
				if err != nil {
					maybeSaveResult.Set(maybe.Some(fmt.Errorf("failed exporting image file: %w", err)))
					return
				}
				defer func() {
					maybeSaveResult.Set(maybe.Some(file.Close()))
				}()
				switch extension {
				case "jpg":
					if err := jpeg.Encode(file, imageResult.Image, nil); err != nil {
						maybeSaveResult.Set(maybe.Some(fmt.Errorf("failed encoding image file: %w", err)))
						return
					}
				case "png":
					if err := png.Encode(file, imageResult.Image); err != nil {
						maybeSaveResult.Set(maybe.Some(fmt.Errorf("failed encoding image file: %w", err)))
						return
					}
				}
			}

		})

		return c.Sequence(
			column.Column(
				c.Sequence(
					text.HeadlineMedium("File Explorer Demo"),
					spacer.Height(16),
					button.Filled(onOpenFile, "Open File"),
					button.Filled(onSaveFile, "Save File", button.WithEnabled(canSaveFile.Get())),
				),
				column.WithSpacing(column.SpaceSides),
				column.WithAlignment(column.Middle),
				column.WithModifier(size.FillMax()),
			),
			c.When(
				canSaveFile.Get(),
				DisplayMaybeImage(maybeImageResult.Get()),
			),
			launchedOpenFile,
			launchedSaveFile,
		)(c)
	}
}

func DisplayMaybeImage(maybeImageResult maybe.Maybe[ImageResult]) api.Composable {
	return func(c api.Composer) api.Composer {
		return c.Sequence(
			c.WhenLazy(
				maybeImageResult.IsSome(),
				func() api.Composable {
					imageResult := maybeImageResult.UnwrapUnsafe()
					return c.IfLazy(
						imageResult.Error == nil,
						func() api.Composable {
							return surface.Surface(
								fImage.Image(
									maybeImageResult.UnwrapUnsafe().ImageResource,
									fImage.WithContentScale(uilayout.ContentScaleFit),
								),
								surface.WithModifier(
									size.Size(150, 100, size.SizeRequired()).
										Then(
											border.Border(
												unit.Dp(1),
												graphics.FromNRGBA(color.NRGBA{0, 0, 0, 255}),
												shape.ShapeRectangle,
											),
										),
								),
							)
						},
						func() api.Composable {
							return text.BodyLarge(imageResult.Error.Error())
						},
					)

				},
			),
		)(c)

	}
}
