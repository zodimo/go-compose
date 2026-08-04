package graphics

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"io/fs"
	"os"
)

// ImageResource represents a decoded image that can be used for rendering.
// It stores the decoded image as a standard library image.Image, keeping
// gioui.org/op/paint types out of the public API surface.
type ImageResource struct {
	img image.Image
}

// NewImageResource creates an ImageResource from a decoded image.Image.
func NewImageResource(img image.Image) ImageResource {
	return ImageResource{img: img}
}

// Image returns the underlying decoded image.
func (ir ImageResource) Image() image.Image {
	return ir.img
}

func NewResourceFromImageFile(imageFile io.Reader) ImageResource {
	return requireImage(imageFile)
}
func NewResourceFromImageFS(assetsFS fs.ReadFileFS, imagePath string) ImageResource {
	return requireImageFromFS(assetsFS, imagePath)
}

func NewResourceFromImageByPath(imagePath string) ImageResource {
	return requireImageByPath(imagePath)
}

// imageFile

func requireImage(imageFile io.Reader) ImageResource {
	decodedImage, _, err := image.Decode(imageFile)
	if err != nil {
		panic(fmt.Errorf("failed to decode image file: %v", err))
	}
	return ImageResource{
		img: decodedImage,
	}
}

func requireImageFromFS(assetsFS fs.ReadFileFS, imagePath string) ImageResource {
	imageBytes, err := assetsFS.ReadFile(imagePath)
	if err != nil {
		panic(fmt.Errorf("failed to open image file: %v", err))
	}
	return requireImage(bytes.NewReader(imageBytes))
}

func requireImageByPath(imagePath string) ImageResource {
	imageFile, err := os.Open(imagePath)
	if err != nil {
		panic(fmt.Errorf("failed to open image file: %v", err))
	}
	defer func(imageFile *os.File) {
		err := imageFile.Close()
		if err != nil {
			panic(fmt.Errorf("failed to close image: %v", err))
		}
	}(imageFile)
	return requireImage(imageFile)
}
