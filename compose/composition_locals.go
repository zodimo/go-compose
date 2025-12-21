package compose

import "github.com/zodimo/go-compose/pkg/api"

// CompositionLocal allows passing data implicitly down the composition tree.
type CompositionLocal[T any] struct {
	defaultValueFactory func() T
}

// CompositionLocalOf creates a new CompositionLocal with a default value factory.
func CompositionLocalOf[T any](defaultValueFactory func() T) *CompositionLocal[T] {
	return &CompositionLocal[T]{defaultValueFactory: defaultValueFactory}
}

// StaticCompositionLocalOf creates a CompositionLocal where the value is unlikely to change.
// For now, this behaves the same as CompositionLocalOf.
func StaticCompositionLocalOf[T any](defaultValueFactory func() T) *CompositionLocal[T] {
	return &CompositionLocal[T]{defaultValueFactory: defaultValueFactory}
}

// Provides creates a ProvidedValue for this CompositionLocal.
func (local *CompositionLocal[T]) Provides(value T) api.ProvidedValue {
	return api.ProvidedValue{
		CompositionLocal: local,
		Value:            value,
	}
}

// Current returns the current value of the CompositionLocal from the Composer.
func (local *CompositionLocal[T]) Current(c api.Composer) T {
	val := c.Consume(local)
	if val == nil {
		return local.defaultValueFactory()
	}
	// We can trust the type because Provides enforces it
	return val.(T)
}

// CompositionLocalProvider scopes the provided values to the content composable.
func CompositionLocalProvider(values []api.ProvidedValue, content Composable) Composable {
	return func(c Composer) Composer {
		c.StartProviders(values)
		content(c)
		c.EndProviders()
		return c
	}
}

// Helper for single value provision
func CompositionLocalProvider1[T any](local *CompositionLocal[T], value T, content Composable) Composable {
	return CompositionLocalProvider([]api.ProvidedValue{local.Provides(value)}, content)
}
