package colorspace

import (
	"fmt"
)

const (
	MinId = -1
	MaxId = 63
)

// ColorSpace identifies a specific organization of colors.
type ColorSpace interface {
	Name() string
	Model() ColorModel
	Id() int
	ComponentCount() int
	IsWideGamut() bool
	IsSrgb() bool
	MinValue(component int) float32
	MaxValue(component int) float32
	ToXyz(v []float32) []float32
	FromXyz(v []float32) []float32
	ToXy(v0, v1, v2 float32) int64
	ToZ(v0, v1, v2 float32) float32
	XyzaToColor(x, y, z, a float32, colorSpace ColorSpace) (uint64, error) // Returns packed Color
}

// BaseColorSpace provides common fields for ColorSpace implementations.
// Embedding this struct allows reusing Name, Model, Id fields and basic methods.
type BaseColorSpace struct {
	name  string
	model ColorModel
	id    int
}

func NewBaseColorSpace(name string, model ColorModel, id int) BaseColorSpace {
	if len(name) == 0 {
		panic("The name of a color space cannot be null and must contain at least 1 character")
	}
	if id < MinId || id > MaxId {
		panic(fmt.Sprintf("The id must be between %d and %d", MinId, MaxId))
	}
	return BaseColorSpace{
		name:  name,
		model: model,
		id:    id,
	}
}

func (b BaseColorSpace) Name() string {
	return b.name
}

func (b BaseColorSpace) Model() ColorModel {
	return b.model
}

func (b BaseColorSpace) Id() int {
	return b.id
}

func (b BaseColorSpace) ComponentCount() int {
	return b.model.ComponentCount()
}

// Default implementation returns false.
func (b BaseColorSpace) IsSrgb() bool {
	return false
}

func (b BaseColorSpace) ToXy(v0, v1, v2 float32) int64 {
	// This method requires calling ToXyz, which is not implemented in BaseColorSpace.
	// Since Go doesn't have abstract classes with virtual method dispatch on self for un-implemented methods
	// in the base struct unless we use an interface for 'self', we might need to pass the interface implementation to this helper,
	// or implement this in the specific structs (Rgb, etc).
	// But to avoid duplication, we can make this a standalone function or have BaseColorSpace not implement it directly
	// but provide a helper.
	// However, if we embed BaseColorSpace, we can't easily call the embedder's ToXyz from here.
	// So we'll leave it to be implemented by the struct embedding BaseColorSpace, utilizing a helper.
	panic("ToXy must be implemented by the specific ColorSpace")
}

func (b BaseColorSpace) ToZ(v0, v1, v2 float32) float32 {
	panic("ToZ must be implemented by the specific ColorSpace")
}

func (b BaseColorSpace) XyzaToColor(x, y, z, a float32, colorSpace ColorSpace) (uint64, error) {
	panic("XyzaToColor must be implemented by the specific ColorSpace")
}

// Helper for XyzaToColor that can be used by implementations
func XyzaToColorHelper(cs ColorSpace, x, y, z, a float32, destColorSpace ColorSpace) (uint64, error) {
	// This creates a circular dependency if we try to use Color here, because Color depends on ColorSpace.
	// We will return uint64 (packed color) to avoid importing the graphics package.
	// Consumers will cast it to Color.
	// Or we pass the values through.
	// The Kotlin implementation calls `Color(colors[0], colors[1], colors[2], a, colorSpace)`
	// We might need to defer this implementation or keep it simple.
	// actually `graphics` package imports `colorspace`, so `colorspace` cannot import `graphics`.
	// So we return values, or packed uint64.
	// The interface signature I defined returns `uint64`.
	// But `Color` packing logic is in `graphics`.
	// So `colorspace` should probably NOT return `color.Color` (uint64) directly if it requires `graphics` logic to pack.
	// Unless packing logic is here?
	// `Color` packing logic (fp16) is generic.
	// Let's defer implementation or move packing logic to `colorspace` or `util`.
	// The Kotlin `Color` class has the packing logic.
	// Maybe `XyzaToColor` should return `[]float32` (components) instead of packed color?
	// Kotlin version: `internal open fun xyzaToColor(...) : Color`
	// Since it is internal, maybe we don't need it in the public interface.
	// Let's remove it from interface if possible, or return components.
	return 0, nil
}
