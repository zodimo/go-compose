package icon

import (
	"fmt"
	"image"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"

	"github.com/go-text/typesetting/font"

	"github.com/zodimo/go-compose/compose/ui/unit"
	"github.com/zodimo/go-compose/internal/layoutnode"
)

// SymbolErrorHandler, when non-nil, is called during composition for every
// symbol icon that cannot be rendered — an unknown name, a missing style asset,
// or a glyph with no outline. Rendering then draws nothing, so a typo is
// reported rather than silently drawn as text.
//
// It is a package-level hook because icon rendering has no error channel: icons
// compose as widgets, not as values that can return an error.
var SymbolErrorHandler func(SymbolName, error)

// reportSymbolError forwards err to SymbolErrorHandler when one is installed.
func reportSymbolError(name SymbolName, err error) {
	if SymbolErrorHandler != nil {
		SymbolErrorHandler(name, err)
	}
}

func errNoLineMetrics(name SymbolName) error {
	return fmt.Errorf("icon: %q has no horizontal line metrics in the symbol font", name)
}

func errNoOutline(name SymbolName) error {
	return fmt.Errorf("icon: %q resolved to a glyph with no vector outline", name)
}

// symbolRender is everything needed to draw one symbol icon: the converted
// geometry plus the placement and size derived from the font's metrics.
type symbolRender struct {
	// geometry is the converted outline. Its coordinates are relative to the
	// glyph origin on the baseline, Y-down.
	geometry pathGeometry

	// baselineFromTop is the distance in pixels from the top of the icon box to
	// the text baseline. Drawing translates by this before filling.
	baselineFromTop int

	// size is the reported icon box, width by height.
	size image.Point

	// baselineFromBottom is the reported layout baseline, measured upward from
	// the bottom of the box as layout.Dimensions expects.
	baselineFromBottom int
}

func (s symbolRender) dimensions() layout.Dimensions {
	return layout.Dimensions{
		Size:     s.size,
		Baseline: s.baselineFromBottom,
	}
}

// draw paints the icon geometry at the resolved position, recording into ops.
//
// ops is passed explicitly rather than read from gtx.Ops because the paint cache
// records into a private op.Ops and replays it later. Drawing into gtx.Ops would
// record into the wrong list and the cached frame would replay empty.
func (s symbolRender) draw(ops *op.Ops, f fill, strokeWidthPx float32) layout.Dimensions {
	if !s.geometry.empty() {
		// The glyph origin sits on the baseline, which is baselineFromTop below
		// the top of the icon box.
		trans := op.Offset(image.Pt(0, s.baselineFromTop)).Push(ops)

		var shape clip.Op
		if strokeWidthPx > 0 {
			shape = strokeShape(s.geometry.path, strokeWidthPx)
		} else {
			shape = clip.Outline{Path: s.geometry.path}.Op()
		}

		f.paint(ops, image.Rectangle{Max: s.size}, shape)
		trans.Pop()
	}
	return s.dimensions()
}

// buildSymbolRender resolves a symbol into drawable geometry at sizePx.
//
// The reported box matches what the previous text-based path measured: width is
// the requested size, and height is the font's ascender plus descender at that
// size. That is what makes the vector path a drop-in replacement — a 24dp icon
// still lays out 24x30 with a 3px baseline, not 24x24.
func buildSymbolRender(style SymbolStyle, name SymbolName, sizePx int, axes axisValues) (symbolRender, error) {
	key := geometryKey{symbol: name, style: style, axes: axes, sizePx: sizePx}
	if cached, ok := globalGeometryCache.get(key); ok {
		return cached, nil
	}

	render, err := extractSymbolRender(style, name, sizePx, axes)
	if err != nil {
		return symbolRender{}, err
	}

	globalGeometryCache.noteConversion()
	return globalGeometryCache.put(key, render), nil
}

// extractSymbolRender does the actual work of resolving a symbol and converting
// its outline. It is separated from buildSymbolRender so the cache boundary is
// explicit.
func extractSymbolRender(style SymbolStyle, name SymbolName, sizePx int, axes axisValues) (symbolRender, error) {
	faces := sharedSymbolFonts()

	renderStyle, fellBack := faces.resolveStyle(style)
	if fellBack {
		// The requested style has no font asset: render outlined and say so.
		if err, ok := faces.errs[style]; ok {
			reportSymbolError(name, err)
		}
	}

	sf, err := faces.face(renderStyle)
	if err != nil {
		return symbolRender{}, err
	}

	face := font.NewFace(sf.font)
	// Apply the axes before extracting anything, so both the shaped glyph and its
	// outline come from the requested instance.
	face.SetVariations(axes.variations())

	// Line metrics place the glyph and size the box. The font's ascender and
	// descender are constant across its variation axes, so this is stable.
	heightPx, baselineFromTop, baselineFromBottom, ok := boxMetrics(face, sf.upem, sizePx)
	if !ok {
		return symbolRender{}, errNoLineMetrics(name)
	}

	glyph, err := resolveGlyph(face, renderStyle, name, sizePx)
	if err != nil {
		return symbolRender{}, err
	}

	outline, ok := glyph.outline()
	if !ok {
		return symbolRender{}, errNoOutline(name)
	}

	geometry := outlineToPath(outline.Segments, sf.upem, sizePx)
	if geometry.empty() {
		return symbolRender{}, errNoOutline(name)
	}

	return symbolRender{
		geometry:           geometry,
		baselineFromTop:    baselineFromTop,
		size:               image.Pt(sizePx, heightPx),
		baselineFromBottom: baselineFromBottom,
	}, nil
}

// strokeWidthPxFromDp converts a stroke width in Dp to pixels.
func strokeWidthPxFromDp(width unit.Dp, density unit.Density) float32 {
	if width <= 0 {
		return 0
	}
	return density.DpToPx(width)
}

// boxMetrics computes the icon box height and the baseline's position, in
// pixels, from the font's ascender and descender.
//
// Ascender and descender are each rounded up independently, which is what the
// text shaper did: at 24px the ascender is 26.4px and the descender 2.4px, and
// the measured box was 27+3 = 30px tall with a 3px baseline.
func boxMetrics(face *font.Face, upem uint16, sizePx int) (height, fromTop, fromBottom int, ok bool) {
	if upem == 0 {
		return 0, 0, 0, false
	}
	extents, ok := face.FontHExtents()
	if !ok {
		return 0, 0, 0, false
	}

	scale := float32(sizePx) / float32(upem)
	ascender := int(ceil32(extents.Ascender * scale))
	descender := int(ceil32(-extents.Descender * scale))
	if ascender <= 0 || descender < 0 {
		return 0, 0, 0, false
	}

	return ascender + descender, ascender, descender, true
}

// symbolWidgetConstructor adapts a symbolRender to the layout node system, so
// icons compose alongside every other widget.
func symbolWidgetConstructor(geomKey geometryKey, render symbolRender, f fill, strokeWidthPx float32) layoutnode.LayoutNodeWidgetConstructor {
	key := paintKey{
		geom:        geomKey,
		color:       f.color,
		strokeWidth: int32(strokeWidthPx),
		brushID:     brushIdentity(f.brush),
	}

	return layoutnode.NewLayoutNodeWidgetConstructor(func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {
		return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
			if dims, ok := globalOpsCache.layout(gtx, key); ok {
				return dims
			}
			return globalOpsCache.record(gtx, key, func(ops *op.Ops) layoutnode.LayoutDimensions {
				// Draw into the ops we are handed: the cache records into a
				// private op.Ops and replays it on later frames.
				return render.draw(ops, f, strokeWidthPx)
			})
		}
	})
}
