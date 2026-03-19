package m3table

import (
	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/compose/ui/unit"
	"github.com/zodimo/go-compose/pkg/api"
)

// Column defines the configuration for a single column in the table.
type Column struct {
	// Header is an optional composable that defines the column header.
	Header api.Composable
	// Weight is the flex weight of the column. If greater than 0, the column
	// scales proportionally based on its weight.
	Weight int
	// Width is the fixed width of the column. This is used if Weight is 0.
	Width unit.Dp
}

type TableOptions struct {
	Modifier ui.Modifier
}

type TableOption func(*TableOptions)

func DefaultTableOptions() TableOptions {
	return TableOptions{
		Modifier: ui.EmptyModifier,
	}
}

func WithModifier(modifier ui.Modifier) TableOption {
	return func(o *TableOptions) {
		o.Modifier = modifier
	}
}
