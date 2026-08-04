// Package render implements the draw-IR seam: the framework-facing surface of
// the rendering backends. It defines the Backend interface that each rendering
// backend (gio, software, etc.) implements, plus the opaque types that flow
// across the seam (Ops, CallOp, DrawCommand).
//
// The seam is designed so that the framework public API (compose/, modifiers/,
// theme/, runtime/, pkg/) never references gioui.org types. Only this package
// and internal/layoutnode may import engine packages.
//
// Backend registration uses blank imports (design D6):
//
//	import _ "github.com/zodimo/go-compose/internal/render/gio"
package render
