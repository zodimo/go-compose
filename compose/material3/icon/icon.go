package icon

import (
	"container/list"
	"hash/fnv"
	"image/color"
	"sync"

	"github.com/zodimo/go-compose/compose/material3"
	"github.com/zodimo/go-compose/internal/layoutnode"
	"github.com/zodimo/go-compose/theme"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
)

// iconCacheKey is the key used to store the GlobalIconCache in the Composer's state.
var iconCacheKey = "icon_global_cache"

const Material3IconNodeID = "Material3Icon"

// icons from golang.org/x/exp/shiny/materialdesign/icons
func Icon(iconByte []byte, options ...IconOption) Composable {
	opts := DefaultIconOptions()
	for _, option := range options {
		if option == nil {
			continue
		}
		option(&opts)
	}

	return func(c Composer) Composer {
		opts.Color = theme.TakeOrElseColor(
			opts.Color,
			theme.ColorHelper.SpecificColor(material3.LocalContentColor.Current(c)),
		)

		c.StartBlock(Material3IconNodeID)
		c.Modifier(func(modifier Modifier) Modifier {
			return modifier.Then(opts.Modifier)
		})

		// Retrieve or initialize the global icon cache from the composer's persistent store
		cacheVal := c.State(iconCacheKey, initCache)
		cache := cacheVal.Get().(*GlobalIconCache)

		c.SetWidgetConstructor(iconWidgetConstructor(opts, iconByte, cache))

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

func iconWidgetConstructor(options IconOptions, iconByte []byte, cache *GlobalIconCache) layoutnode.LayoutNodeWidgetConstructor {
	// Pre-calculate hash for this icon data
	h := fnv.New64a()
	h.Write(iconByte)
	dataHash := h.Sum64()

	return layoutnode.NewLayoutNodeWidgetConstructor(func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {
		// We still keep the local widget instance for fallback/misses
		iconWidget := requireIconWidget(iconByte)

		return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {

			nrgba := theme.GetThemeManager().ResolveColorDescriptor(options.Color).AsNRGBA()

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
