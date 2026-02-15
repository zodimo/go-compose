package state

import "github.com/zodimo/go-compose/internal/immap"

type SlotStore = immap.ImmutableMap[any]

type SlotStoreTyped[T any] = immap.ImmutableMap[T]

func EmptySlotStore[T any]() SlotStoreTyped[T] {
	return immap.EmptyImmutableMap[T]()
}
