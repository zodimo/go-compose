package state

import (
	"github.com/zodimo/go-compose/internal/immap"
)

type Memo = immap.ImmutableMap[any]

type MemoTyped[T any] = immap.ImmutableMap[T]

func EmptyMemo[T any]() MemoTyped[T] {
	return immap.EmptyImmutableMap[T]()
}
