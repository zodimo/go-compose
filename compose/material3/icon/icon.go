package icon

import (
	"container/list"
	"hash/fnv"
	"image/color"
	"sync"

	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/foundation/text"
	"github.com/zodimo/go-compose/compose/material3"
	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/next/platform"
	uitext "github.com/zodimo/go-compose/compose/ui/text"
	"github.com/zodimo/go-compose/compose/ui/text/font"
	"github.com/zodimo/go-compose/compose/ui/unit"
	"github.com/zodimo/go-compose/internal/layoutnode"
	"github.com/zodimo/go-compose/modifiers/size"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
)

// iconCacheKey is the key used to store the GlobalIconCache in the Composer's state.
var iconCacheKey = "material3_icon_global_cache"

const Material3IconNodeID = "Material3Icon"

// MaterialSymbolsFontFamily is the font family for Material Symbols.
// It assumes the font "Material Symbols Outlined" has been loaded into the shaper.
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
		// Fallback: if unknown type, return empty composable
		return func(c Composer) Composer { return c }
	}
}

// iconFromBytes renders an icon from raw byte data (SVG icons).
func iconFromBytes(iconByte []byte, opts IconOptions) Composable {
	return func(c Composer) Composer {
		opts.Color = opts.Color.TakeOrElse(material3.LocalContentColor.Current(c))

		c.StartBlock(Material3IconNodeID)
		c.Modifier(func(modifier ui.Modifier) ui.Modifier {
			baseModifier := modifier.Then(opts.Modifier)
			// Apply size constraint if specified
			if opts.Size.IsSpecified() {
				sizeInt := int(opts.Size)
				baseModifier = baseModifier.Then(size.Size(sizeInt, sizeInt))
			}
			return baseModifier
		})

		// Retrieve or initialize the global icon cache from the composer's persistent store
		cacheVal := c.State(iconCacheKey, initCache)
		cache := cacheVal.Get().(*GlobalIconCache)

		c.SetWidgetConstructor(iconWidgetConstructor(opts, iconByte, cache))

		return c.EndBlock()
	}
}

// iconFromSymbol renders an icon using the Material Symbols font.
func iconFromSymbol(name SymbolName, opts IconOptions) Composable {
	return func(c Composer) Composer {
		// Resolve color
		contentColor := material3.LocalContentColor.Current(c)
		iconColor := opts.Color.TakeOrElse(contentColor)

		layoutDirection := platform.LocalLayoutDirection.Current(c)
		// Resolve text style with defaults
		defaultTextStyles := uitext.TextStyleResolveDefaults(uitext.TextStyleUnspecified, layoutDirection)

		// Resolve size - prefer unified Size, fallback to FontSize, then default 24sp
		var fontSize unit.TextUnit
		if opts.Size.IsSpecified() {
			// Convert Dp to Sp (treat them as equivalent for font sizing)
			fontSize = unit.Sp(opts.Size.Value())
		} else {
			fontSize = opts.FontSize.TakeOrElse(defaultTextStyles.FontSize())
		}

		return box.Box(text.Text(
			string(name),
			text.WithModifier(opts.Modifier),
			text.WithTextStyle(
				uitext.TextStyleFromOptions(
					uitext.WithFontFamily(MaterialSymbolsFontFamily),
					uitext.WithFontSize(fontSize),
					uitext.WithColor(iconColor),
				),
			),
		),
			box.WithModifier(size.WrapContentSize()),
		)(c)
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

func iconWidgetConstructor(options IconOptions, iconByte []byte, cache *GlobalIconCache) layoutnode.LayoutNodeWidgetConstructor {
	// Pre-calculate hash for this icon data
	h := fnv.New64a()
	h.Write(iconByte)
	dataHash := h.Sum64()

	return layoutnode.NewLayoutNodeWidgetConstructor(func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {
		// We still keep the local widget instance for fallback/misses
		iconWidget := requireIconWidget(iconByte)

		return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {

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
