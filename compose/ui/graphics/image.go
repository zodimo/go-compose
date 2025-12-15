package graphics

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"io/fs"
	"log"
	"os"

	"gioui.org/op/paint"
)

type ImageResource struct {
	ImageOp paint.ImageOp
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
	decodedImage, format, err := image.Decode(imageFile)
	if err != nil {
		panic(fmt.Errorf("failed to decode image file: %v", err))
	}
	size := decodedImage.Bounds().Size()
	log.Printf("decoded iamge format: %s, size: %s\n", format, size)
	return ImageResource{
		ImageOp: paint.NewImageOp(decodedImage),
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
