package navigation

import "github.com/zodimo/go-compose/pkg/api"

type NavGraphBuilder struct {
	destinations map[string]api.Composable
}

func NewNavGraphBuilder() *NavGraphBuilder {
	return &NavGraphBuilder{
		destinations: make(map[string]api.Composable),
	}
}

func (b *NavGraphBuilder) Composable(route string, content api.Composable) {
	b.destinations[route] = content
}
