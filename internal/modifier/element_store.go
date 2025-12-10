package modifier

import "github.com/zodimo/go-maybe"

type ElementStore interface {
	SetElement(string, Element) ElementStore
	GetElement(string) maybe.Maybe[Element]
}
