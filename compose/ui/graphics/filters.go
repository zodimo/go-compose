package graphics

import skiaiface "github.com/zodimo/go-skia-support/skia/interfaces"

// Type aliases to re-export Skia's ColorFilter/MaskFilter interfaces
// under the graphics package so we can populate Paint fields with them.
type ColorFilter = skiaiface.ColorFilter
type MaskFilter = skiaiface.MaskFilter
