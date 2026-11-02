package icon

import (
	"container/list"
	"sync"
)

// defaultGeometryCacheLimit bounds how many distinct icon geometries are kept.
// Icons are small and highly repetitive — a screen shows tens at most — so a few
// hundred entries covers normal use while keeping memory flat.
const defaultGeometryCacheLimit = 256

// geometryKey identifies one converted icon geometry. It is built from the
// canonical, post-clamp axis values and the integer pixel size, never from raw
// Dp or float sizes, so configurations that should be identical share an entry
// and configurations that differ cannot collide.
type geometryKey struct {
	symbol SymbolName
	style  SymbolStyle
	axes   axisValues
	sizePx int
}

type geometryEntry struct {
	key    geometryKey
	render symbolRender
}

// geometryCache is a bounded least-recently-used cache of converted geometry.
// It mirrors the shape of GlobalIconCache (mutex + map + container/list) so the
// two icon paths behave the same way.
type geometryCache struct {
	mu          sync.Mutex
	entries     map[geometryKey]*list.Element
	list        *list.List
	limit       int
	extractions int
}

func newGeometryCache(limit int) *geometryCache {
	if limit <= 0 {
		limit = defaultGeometryCacheLimit
	}
	return &geometryCache{
		entries: make(map[geometryKey]*list.Element),
		list:    list.New(),
		limit:   limit,
	}
}

// get returns the cached render for key, moving it to the front.
func (c *geometryCache) get(key geometryKey) (symbolRender, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.entries[key]; ok {
		c.list.MoveToFront(elem)
		return elem.Value.(*geometryEntry).render, true
	}
	return symbolRender{}, false
}

// put stores render for key, evicting the least-recently-used entry if full.
func (c *geometryCache) put(key geometryKey, render symbolRender) symbolRender {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.entries[key]; ok {
		c.list.MoveToFront(elem)
		return elem.Value.(*geometryEntry).render
	}

	for c.list.Len() >= c.limit {
		oldest := c.list.Back()
		if oldest == nil {
			break
		}
		c.list.Remove(oldest)
		delete(c.entries, oldest.Value.(*geometryEntry).key)
	}

	elem := c.list.PushFront(&geometryEntry{key: key, render: render})
	c.entries[key] = elem
	return render
}

// len reports the number of cached entries.
func (c *geometryCache) len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.list.Len()
}

// conversions reports how many geometries were actually built, which is how
// tests observe cache hits.
func (c *geometryCache) conversions() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.extractions
}

func (c *geometryCache) noteConversion() {
	c.mu.Lock()
	c.extractions++
	c.mu.Unlock()
}

// globalGeometryCache is the process-wide cache. Geometry depends only on the
// key, not on the frame or the composer, so one cache serves every caller.
var globalGeometryCache = newGeometryCache(defaultGeometryCacheLimit)
