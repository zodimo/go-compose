package icon

import (
	"encoding/binary"
	"hash"
	"hash/fnv"
	"image"
	"math"

	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/clip"

	"github.com/go-text/typesetting/font"
	"github.com/go-text/typesetting/font/opentype"
)

// pathGeometry is a glyph outline converted to a drawable path, together with
// the measurements a renderer needs. clip.PathSpec exposes no accessors, so
// command count and ink bounds are carried alongside it rather than recovered
// from it.
type pathGeometry struct {
	// path is the converted outline. It records into its own op.Ops, so it is
	// independent of the frame that produced it and safe to cache and replay.
	path clip.PathSpec

	// commands is the number of path commands emitted. Zero means the outline
	// was empty or unusable, which callers must treat as "nothing to draw".
	commands int

	// ink is the outline's ink bounds in pixels, relative to the glyph's origin
	// on its baseline. Y increases downward, matching Gio's convention, so ink
	// above the baseline has a negative Min.Y.
	ink image.Rectangle

	// signature is a checksum of every converted coordinate. Two renders with the
	// same signature have identical geometry; clip.PathSpec is opaque, so this is
	// the observable proof that axis settings actually changed the outline.
	signature uint64
}

// empty reports whether the outline produced nothing drawable.
func (g pathGeometry) empty() bool { return g.commands == 0 }

// outlineToPath converts a glyph outline into a Gio path scaled to sizePx.
//
// Coordinates are mapped as screenX = x*scale and screenY = -y*scale: the Y axis
// is flipped from the font's Y-up convention to Gio's Y-down convention, and the
// origin sits on the glyph's baseline at its left origin. Callers place the
// result by translating by the ascent, since the outline carries no line-box
// information of its own.
func outlineToPath(segments []font.Segment, upem uint16, sizePx int) pathGeometry {
	// A zero upem would make every coordinate +Inf and produce a path that
	// silently draws nothing.
	if upem == 0 {
		return pathGeometry{}
	}

	scale := float32(sizePx) / float32(upem)
	point := func(p font.SegmentPoint) f32.Point {
		return f32.Pt(p.X*scale, -p.Y*scale)
	}

	var ops op.Ops
	var path clip.Path
	path.Begin(&ops)
	commands := 0
	sig := fnv.New64a()

	for _, seg := range segments {
		switch seg.Op {
		case opentype.SegmentOpMoveTo:
			path.MoveTo(point(seg.Args[0]))
			hashPoint(sig, seg.Args[0])
			commands++
		case opentype.SegmentOpLineTo:
			path.LineTo(point(seg.Args[0]))
			hashPoint(sig, seg.Args[0])
			commands++
		case opentype.SegmentOpQuadTo:
			path.QuadTo(point(seg.Args[0]), point(seg.Args[1]))
			hashPoint(sig, seg.Args[0])
			hashPoint(sig, seg.Args[1])
			commands++
		case opentype.SegmentOpCubeTo:
			path.CubeTo(point(seg.Args[0]), point(seg.Args[1]), point(seg.Args[2]))
			hashPoint(sig, seg.Args[0])
			hashPoint(sig, seg.Args[1])
			hashPoint(sig, seg.Args[2])
			commands++
		}
	}

	// Glyph outlines are implicitly closed by the rasteriser: a contour's final
	// point joins its first. Path.MoveTo ends the preceding contour, so
	// multi-contour glyphs (most icons have holes) convert correctly without an
	// explicit Close. This mirrors Gio's own shaper.Shape conversion.
	spec := path.End()

	return pathGeometry{
		path:      spec,
		commands:  commands,
		ink:       inkBounds(segments, upem, sizePx),
		signature: sig.Sum64(),
	}
}

// inkBounds returns the outline's ink bounds in pixels, Y-down.
func inkBounds(segments []font.Segment, upem uint16, sizePx int) image.Rectangle {
	minX, minY, maxX, maxY, ok := outlineBounds(segments)
	if !ok || upem == 0 {
		return image.Rectangle{}
	}

	scale := float32(sizePx) / float32(upem)
	// Font space is Y-up, so the topmost point is the largest Y and becomes the
	// smallest screen Y.
	return image.Rect(
		int(floor32(minX*scale)),
		int(floor32(-maxY*scale)),
		int(ceil32(maxX*scale)),
		int(ceil32(-minY*scale)),
	)
}

// outlineBounds returns the ink bounds of segments in font units, which the
// renderer needs in order to report dimensions matching what the glyph actually
// occupies.
func outlineBounds(segments []font.Segment) (minX, minY, maxX, maxY float32, ok bool) {
	minX, minY = float32(1e9), float32(1e9)
	maxX, maxY = float32(-1e9), float32(-1e9)

	for _, seg := range segments {
		for _, p := range seg.ArgsSlice() {
			if p.X < minX {
				minX = p.X
			}
			if p.X > maxX {
				maxX = p.X
			}
			if p.Y < minY {
				minY = p.Y
			}
			if p.Y > maxY {
				maxY = p.Y
			}
		}
	}

	return minX, minY, maxX, maxY, maxX >= minX && maxY >= minY
}

// hashPoint folds one segment point into the geometry signature.
func hashPoint(h hash.Hash64, p font.SegmentPoint) {
	var buf [8]byte
	binary.LittleEndian.PutUint32(buf[0:4], math.Float32bits(p.X))
	binary.LittleEndian.PutUint32(buf[4:8], math.Float32bits(p.Y))
	h.Write(buf[:])
}

func floor32(v float32) float32 {
	i := float32(int(v))
	if i > v {
		return i - 1
	}
	return i
}

func ceil32(v float32) float32 {
	i := float32(int(v))
	if i < v {
		return i + 1
	}
	return i
}
