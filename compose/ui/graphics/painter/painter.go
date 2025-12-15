package painter

import (
	"gioui.org/layout"
)

// Painter defines the API for drawing content that can be scaled and positioned.
// It is an abstraction utilized by the Image composable to handle various types of visual content.
type Painter interface {
	// Draw renders the painter content into the provided ops with the given size.
	// The alpha parameter determines the opacity of the drawing (0.0 to 1.0).
	// colorFilter is reserved for future use and should be ignored for now.
	Draw(gtx layout.Context, size layout.Dimensions, alpha float32, colorFilter any) layout.Dimensions

	// IntrinsicSize returns the intrinsic size of the underlying content.
	// If the content has no intrinsic size, it returns (0, 0).
	IntrinsicSize() layout.Dimensions
}

// ColorFilter is a placeholder for future color manipulation support.
type ColorFilter interface{}
