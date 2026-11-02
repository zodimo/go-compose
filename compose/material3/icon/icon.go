package icon

import (
	"container/list"
	"fmt"
	"hash/fnv"
	"image/color"
	"sync"

	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/material3"
	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/platform"
	uitext "github.com/zodimo/go-compose/compose/ui/text"
	"github.com/zodimo/go-compose/compose/ui/text/font"
	"github.com/zodimo/go-compose/compose/ui/unit"
	"github.com/zodimo/go-compose/internal/layoutnode"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
)

// iconCacheKey is the key used to store the GlobalIconCache in the Composer's state.
var iconCacheKey = "material3_icon_global_cache"

const Material3IconNodeID = "Material3Icon"

// MaterialSymbolsFontFamily is the font family for Material Symbols.
// It assumes the font "Material Symbols Outlined" has been loaded into the shaper.
//
// Deprecated: symbol icons no longer render through the text shaper, so they do
// not consult the shaper's font collection. This remains exported only for
// callers that referenced it directly; icons need no font registration.
var MaterialSymbolsFontFamily = font.NewGenericFontFamily("Material Symbols Outlined", "Material Symbols Outlined")

// Icon renders an icon from an IconSource.
// IconSource can be:
//   - IconBytes: raw icon data from golang.org/x/exp/shiny/materialdesign/icons
//   - SymbolName: a Material Symbol name (e.g. SymbolSearch, SymbolHome)
//
// Examples:
//
//	icon.Icon(icon.IconBytes(icons.ActionHome), icon.WithColor(color))
//	icon.Icon(icon.SymbolSearch, icon.WithSymbolSize(unit.Sp(24)))
func Icon(source IconSource, options ...IconOption) Composable {
	opts := DefaultIconOptions()
	for _, option := range options {
		if option == nil {
			continue
		}
		option(&opts)
	}

	switch src := source.(type) {
	case IconBytes:
		return iconFromBytes([]byte(src), opts)
	case SymbolName:
		return iconFromSymbol(src, opts)
	default:
		panic(fmt.Sprintf("Unknown IconSource type: %T", src))
	}
}

// iconFromBytes renders an icon from raw byte data (SVG icons).
func iconFromBytes(iconByte []byte, opts IconOptions) Composable {
	return func(c Composer) Composer {
		// Resolve color
		contentColor := material3.LocalContentColor.Current(c)
		opts.Color = graphics.ColorCompanion.TakeOrElse(opts.Color, contentColor)

		localTextStyle := compose.LocalTextStyle.Current(c)
		localDensity := compose.LocalDensity.Current(c)
		layoutDirection := platform.LocalLayoutDirection.Current(c)

		// Resolve text style with defaults
		defaultTextStyles := uitext.ResolveTextStyleDefaults(localTextStyle, layoutDirection)

		iconSizeInPx := localDensity.DpRoundToPx(unit.DpCompanion.TakeOrElse(opts.Size, localDensity.TextUnitToDp(defaultTextStyles.FontSize())))

		c.StartBlock(Material3IconNodeID)
		c.Modifier(func(modifier ui.Modifier) ui.Modifier {
			baseModifier := modifier.Then(opts.Modifier)
			return baseModifier
		})

		// Retrieve or initialize the global icon cache from the composer's persistent store
		cacheVal := c.State(iconCacheKey, initCache)
		cache := cacheVal.Get().(*GlobalIconCache)

		c.SetWidgetConstructor(iconWidgetConstructor(opts, iconByte, cache, iconSizeInPx))

		return c.EndBlock()
	}
}

// iconFromSymbol renders a Material Symbol as a vector outline.
//
// Unlike the text-based rendering this replaces, the glyph is extracted from the
// variable font directly, so its axes are reachable and its geometry can be
// filled or stroked like any other shape.
func iconFromSymbol(name SymbolName, opts IconOptions) Composable {
	return func(c Composer) Composer {
		contentColor := material3.LocalContentColor.Current(c)
		iconColor := graphics.ColorCompanion.TakeOrElse(opts.Color, contentColor)

		localTextStyle := compose.LocalTextStyle.Current(c)
		localDensity := compose.LocalDensity.Current(c)
		layoutDirection := platform.LocalLayoutDirection.Current(c)

		// Unspecified size falls back to the inherited font size, as before.
		// The deprecated WithSymbolSize is honoured as the first fallback so
		// callers that used it keep the size they asked for.
		defaultTextStyles := uitext.ResolveTextStyleDefaults(localTextStyle, layoutDirection)
		fallbackDp := localDensity.TextUnitToDp(defaultTextStyles.FontSize())
		if unit.TextUnitCompanion.IsSpecified(opts.FontSize) {
			fallbackDp = localDensity.TextUnitToDp(opts.FontSize)
		}
		sizeDp := unit.DpCompanion.TakeOrElse(opts.Size, fallbackDp)
		sizePx := localDensity.DpRoundToPx(sizeDp)

		c.StartBlock(Material3IconNodeID)
		c.Modifier(func(modifier ui.Modifier) ui.Modifier {
			return modifier.Then(opts.Modifier)
		})

		// Optical size is a design-space axis, so it follows the logical size
		// rather than the device pixel size: the same icon selects the same
		// optical design on every screen density.
		axes := resolveAxes(opts, sizeDp.Value())
		render, err := buildSymbolRender(opts.SymbolStyle, name, sizePx, axes)
		if err != nil {
			// Nothing to draw; report rather than render an empty box silently.
			// The box keeps the usual metrics so a typo does not move layout.
			reportSymbolError(name, err)
			render = emptySymbolRender(opts.SymbolStyle, sizePx)
		}

		// Brush wins over an explicit color, which wins over the ambient
		// content color resolved above.
		paintStyle := fill{
			brush: opts.Brush,
			color: graphics.ColorToNRGBA(iconColor),
		}
		var strokeWidthPx float32
		if opts.hasStroke {
			strokeWidthPx = strokeWidthPxFromDp(opts.StrokeWidth, localDensity)
		}

		geomKey := geometryKey{symbol: name, style: opts.SymbolStyle, axes: axes, sizePx: sizePx}
		c.SetWidgetConstructor(symbolWidgetConstructor(geomKey, render, paintStyle, strokeWidthPx))

		return c.EndBlock()
	}
}

// GlobalIconCache manages the LRU cache for icon rendering operations
type GlobalIconCache struct {
	mu    sync.Mutex
	cache map[cacheKey]*list.Element
	list  *list.List
	limit int
}

func NewGlobalIconCache() any {
	return &GlobalIconCache{
		cache: make(map[cacheKey]*list.Element),
		list:  list.New(),
		limit: 100,
	}
}

// Global cache structures
type cacheKey struct {
	dataHash    uint64
	constraints layout.Constraints
	color       color.NRGBA
}

type cacheEntry struct {
	key  cacheKey
	ops  *op.Ops
	call op.CallOp
	dims layout.Dimensions
}

func initCache() any {
	return NewGlobalIconCache()
}

func iconWidgetConstructor(options IconOptions, iconByte []byte, cache *GlobalIconCache, iconSizeInPx int) layoutnode.LayoutNodeWidgetConstructor {
	// Pre-calculate hash for this icon data
	h := fnv.New64a()
	h.Write(iconByte)
	dataHash := h.Sum64()

	return layoutnode.NewLayoutNodeWidgetConstructor(func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {
		// We still keep the local widget instance for fallback/misses
		iconWidget := requireIconWidget(iconByte)

		return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {

			gtx.Constraints.Min.X = iconSizeInPx

			nrgba := graphics.ColorToNRGBA(options.Color)

			key := cacheKey{
				dataHash:    dataHash,
				constraints: gtx.Constraints,
				color:       nrgba,
			}

			cache.mu.Lock()
			if elem, hit := cache.cache[key]; hit {
				cache.list.MoveToFront(elem)
				entry := elem.Value.(*cacheEntry)
				// Replay the cached op
				entry.call.Add(gtx.Ops)
				cache.mu.Unlock()
				// fmt.Println("GLOBAL CACHE HIT")
				return entry.dims
			}
			cache.mu.Unlock()

			// fmt.Println("GLOBAL CACHE MISS")

			// Record new ops
			// We need a fresh op.Ops for the cache entry
			entryOps := new(op.Ops)
			macro := op.Record(entryOps)

			// Use a context targeting the entryOps
			gtxCache := gtx
			gtxCache.Ops = entryOps

			dims := iconWidget(gtxCache, nrgba)
			call := macro.Stop()

			// Store in cache
			cache.mu.Lock()
			// Check capacity
			if cache.list.Len() >= cache.limit {
				oldest := cache.list.Back()
				if oldest != nil {
					cache.list.Remove(oldest)
					oldEntry := oldest.Value.(*cacheEntry)
					delete(cache.cache, oldEntry.key)
				}
			}

			newEntry := &cacheEntry{
				key:  key,
				ops:  entryOps,
				call: call,
				dims: dims,
			}
			elem := cache.list.PushFront(newEntry)
			cache.cache[key] = elem
			cache.mu.Unlock()

			// Add to current frame
			call.Add(gtx.Ops)

			return dims
		}
	})
}

func requireIconWidget(data []byte) IconWidget {
	iconWidget, err := widget.NewIcon(data)
	if err != nil {
		panic(err)
	}
	return func(gtx layout.Context, foreground color.NRGBA) layout.Dimensions {
		return iconWidget.Layout(gtx, foreground)
	}
}
