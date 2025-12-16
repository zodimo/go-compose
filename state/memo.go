package state

import (
	"github.com/zodimo/go-compose/internal/immap"
)

// Memo is an alias for an immutable map holding values of any type.
// It is used for memoization and efficient state storage.
type Memo = immap.ImmutableMap[any]

// MemoTyped is an alias for an immutable map holding values of a specific type T.
type MemoTyped[T any] = immap.ImmutableMap[T]

// EmptyMemo returns an empty typed immutable map.
func EmptyMemo[T any]() MemoTyped[T] {
	return immap.EmptyImmutableMap[T]()
}
