package m3table

import (
	"github.com/zodimo/go-compose/compose/foundation/layout/box"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/row"
	"github.com/zodimo/go-compose/compose/material3/divider"
	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/modifiers/size"
	"github.com/zodimo/go-compose/modifiers/weight"
	"github.com/zodimo/go-compose/pkg/api"
)

// Table constructs a basic material3 table layout given a list of columns,
// the number of data rows, and a factory function to create each cell's content.
func Table(
	columns []Column,
	rowCount int,
	cellContent func(row, col int) api.Composable,
	options ...TableOption,
) api.Composable {

	opts := DefaultTableOptions()
	for _, option := range options {
		if option != nil {
			option(&opts)
		}
	}

	return func(c api.Composer) api.Composer {
		c.StartBlock("Table")

		hasHeaders := false
		for _, col := range columns {
			if col.Header != nil {
				hasHeaders = true
				break
			}
		}

		c.WithComposable(column.Column(
			c.Sequence(
				c.When(hasHeaders, func(c api.Composer) api.Composer {
					return c.Sequence(
						row.Row(
							func(c api.Composer) api.Composer {
								for i, col := range columns {
									c.Key(i, wrapCell(col, col.Header))
								}
								return c
							},
							row.WithModifier(size.MinHeight(int(opts.MinHeaderRowHeight))),
						),
						divider.Divider(),
					)(c)
				}),
				c.Range(rowCount, func(r int) api.Composable {
					return row.Row(
						func(c api.Composer) api.Composer {
							for cIdx, col := range columns {
								c.Key(cIdx, wrapCell(col, cellContent(r, cIdx)))
							}
							return c
						},
						row.WithModifier(size.MinHeight(int(opts.MinRowHeight))),
					)
				}),
			),
			column.WithModifier(opts.Modifier),
		))

		return c.EndBlock()
	}
}

// wrapCell wraps the provided content into a box with proper width modifiers.
func wrapCell(col Column, content api.Composable) api.Composable {
	var mod ui.Modifier = ui.EmptyModifier
	if col.Weight > 0 {
		mod = mod.Then(weight.Weight(col.Weight))
	} else if col.Width > 0 {
		mod = mod.Then(size.Width(int(col.Width)))
	}

	if content == nil {
		content = func(c api.Composer) api.Composer { return c }
	}

	return box.Box(content, box.WithModifier(mod))
}
