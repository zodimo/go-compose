package icon

import (
	"container/list"
	"image/color"
	"sync"

	"gioui.org/layout"
	"gioui.org/op"
)

// paintKey identifies a painted variant of one geometry: the geometry itself
// plus everything that affects appearance but not shape.
//
// Splitting this from the geometry key is what keeps paint changes cheap: hover
// and theme transitions change color every frame while the axes rarely change,
// so only the cheap op-replay layer is missed.
type paintKey struct {
	geom  geometryKey
	color color.NRGBA
	// strokeWidth is the bit pattern of the resolved stroke width in pixels.
	// Bits rather than an integer conversion: a sub-pixel stroke must never
	// share a key with the unstroked fill, or the cached fill would be
	// replayed for a stroked icon.
	strokeWidth uint32
	// brushKind distinguishes "no brush" from any brush. Brushes are compared by
	// identity pointer: callers that rebuild an equal brush per frame simply miss
	// the cache rather than rendering incorrectly.
	brushID uintptr
}

type paintEntry struct {
	key  paintKey
	ops  *op.Ops
	call op.CallOp
	dims layout.Dimensions
}

// opsCache is a bounded LRU of recorded paint operations.
type opsCache struct {
	mu      sync.Mutex
	entries map[paintKey]*list.Element
	list    *list.List
	limit   int
	hits    int
	misses  int
}

func newOpsCache(limit int) *opsCache {
	if limit <= 0 {
		limit = defaultGeometryCacheLimit
	}
	return &opsCache{
		entries: make(map[paintKey]*list.Element),
		list:    list.New(),
		limit:   limit,
	}
}

// layout returns the cached dimensions for key, replaying the recorded ops into
// gtx.Ops.
func (c *opsCache) layout(gtx layout.Context, key paintKey) (layout.Dimensions, bool) {
	c.mu.Lock()
	elem, ok := c.entries[key]
	if ok {
		c.list.MoveToFront(elem)
	}
	c.mu.Unlock()

	if !ok {
		c.mu.Lock()
		c.misses++
		c.mu.Unlock()
		return layout.Dimensions{}, false
	}

	entry := elem.Value.(*paintEntry)
	entry.call.Add(gtx.Ops)

	c.mu.Lock()
	c.hits++
	c.mu.Unlock()

	return entry.dims, true
}

// record runs draw into a private op.Ops, stores the result under key, and
// replays it into gtx so the current frame draws too.
func (c *opsCache) record(gtx layout.Context, key paintKey, draw func(*op.Ops) layout.Dimensions) layout.Dimensions {
	entryOps := new(op.Ops)
	macro := op.Record(entryOps)

	recordingGtx := gtx
	recordingGtx.Ops = entryOps
	dims := draw(entryOps)

	call := macro.Stop()
	call.Add(gtx.Ops)

	c.mu.Lock()
	for c.list.Len() >= c.limit {
		oldest := c.list.Back()
		if oldest == nil {
			break
		}
		c.list.Remove(oldest)
		delete(c.entries, oldest.Value.(*paintEntry).key)
	}
	elem := c.list.PushFront(&paintEntry{key: key, ops: entryOps, call: call, dims: dims})
	c.entries[key] = elem
	c.misses++
	c.mu.Unlock()

	return dims
}

// stats reports hit and miss counts.
func (c *opsCache) stats() (hits, misses int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.hits, c.misses
}

func (c *opsCache) len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.list.Len()
}

// globalOpsCache caches recorded paint ops for every icon.
var globalOpsCache = newOpsCache(defaultGeometryCacheLimit)

// brushIdentity returns a comparable identity for a brush. A nil brush yields 0,
// which cannot collide with a real pointer.
func brushIdentity(b any) uintptr {
	if b == nil {
		return 0
	}
	return pointerOf(b)
}
