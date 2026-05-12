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
	Modifier           ui.Modifier
	MinHeaderRowHeight unit.Dp
	MinRowHeight       unit.Dp
}

type TableOption func(*TableOptions)

func DefaultTableOptions() TableOptions {
	return TableOptions{
		Modifier:           ui.EmptyModifier,
		MinHeaderRowHeight: 56, // Material 3 standard data table header row height
		MinRowHeight:       52, // Material 3 standard data table row height
	}
}

func WithModifier(modifier ui.Modifier) TableOption {
	return func(o *TableOptions) {
		o.Modifier = modifier
	}
}

// WithMinHeaderRowHeight sets the minimum height for the header row.
func WithMinHeaderRowHeight(height unit.Dp) TableOption {
	return func(o *TableOptions) {
		o.MinHeaderRowHeight = height
	}
}

// WithMinRowHeight sets the minimum height for data rows.
func WithMinRowHeight(height unit.Dp) TableOption {
	return func(o *TableOptions) {
		o.MinRowHeight = height
	}
}
