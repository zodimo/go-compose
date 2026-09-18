package formengine

import (
	"strconv"
	"strings"
)

// ResolvePath resolves a dotted path such as "a.b.c", "items[0].name", or mixed
// forms like "group.items[2].field" against the given root node.
//
// A path segment is either a plain map key ("a") or an array step ("items[0]"):
// the key part resolves from the current Group, and the "[index]" part resolves
// the nth child of the resulting Array. A resolution miss (unknown key,
// out-of-range index, wrong node kind, or malformed brackets) returns
// (nil, false). The empty path resolves to the root itself.
//
// Array indices are positional and shift on Remove/Insert. Resolve paths early
// and hold the returned node reference; do not cache path strings across tree
// mutations.
func ResolvePath(root FormNode, path string) (FormNode, bool) {
	if root == nil {
		return nil, false
	}
	if path == "" {
		return root, true
	}

	cur := root
	for _, seg := range parsePath(path) {
		if seg.isIndex {
			grp, ok := cur.(*Group)
			if !ok {
				return nil, false
			}
			arrNode, ok := grp.child(seg.key)
			if !ok {
				return nil, false
			}
			arr, ok := arrNode.(*Array)
			if !ok {
				return nil, false
			}
			next, ok := arr.At(seg.index)
			if !ok {
				return nil, false
			}
			cur = next
			continue
		}

		grp, ok := cur.(*Group)
		if !ok {
			return nil, false
		}
		next, ok := grp.child(seg.key)
		if !ok {
			return nil, false
		}
		cur = next
	}
	return cur, true
}

type pathSegment struct {
	key     string
	index   int
	isIndex bool
}

// parsePath splits a path on "." and classifies each segment as a plain map key
// or an array step ("key[index]"). Malformed brackets are kept as a plain key,
// which will not resolve.
func parsePath(path string) []pathSegment {
	parts := strings.Split(path, ".")
	segs := make([]pathSegment, 0, len(parts))
	for _, p := range parts {
		if p == "" {
			// An empty segment (double dot or trailing dot) can never resolve.
			segs = append(segs, pathSegment{key: p})
			continue
		}
		if i := strings.IndexByte(p, '['); i >= 0 {
			key := p[:i]
			rest := p[i+1:]
			if j := strings.IndexByte(rest, ']'); j >= 0 {
				if idx, err := strconv.Atoi(rest[:j]); err == nil {
					segs = append(segs, pathSegment{key: key, index: idx, isIndex: true})
					continue
				}
			}
			segs = append(segs, pathSegment{key: p})
			continue
		}
		segs = append(segs, pathSegment{key: p})
	}
	return segs
}

// segmentOf returns the path segment that identifies child relative to parent.
// For a Group parent the segment is the map key; for an Array parent it is "[i]".
func segmentOf(parent FormNode, child FormNode) (string, bool) {
	switch p := parent.(type) {
	case *Group:
		for name, node := range p.childrenMap() {
			if node == child {
				return name, true
			}
		}
	case *Array:
		for i, node := range p.childrenList() {
			if node == child {
				return "[" + strconv.Itoa(i) + "]", true
			}
		}
	}
	return "", false
}

// pathOf builds the dotted path of node from the root by walking parent links.
// Array child segments ("[i]") are joined to the preceding segment without an
// intervening dot, yielding forms like "items[0].name".
func pathOf(node FormNode) string {
	var segs []string
	cur := node
	for {
		p := cur.Parent()
		if p == nil {
			break
		}
		seg, ok := segmentOf(p, cur)
		if !ok {
			break
		}
		segs = append(segs, seg)
		cur = p
	}

	var b strings.Builder
	for i := len(segs) - 1; i >= 0; i-- {
		seg := segs[i]
		if b.Len() > 0 && !strings.HasPrefix(seg, "[") {
			b.WriteByte('.')
		}
		b.WriteString(seg)
	}
	return b.String()
}
