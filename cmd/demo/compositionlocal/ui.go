package main

import (
	"fmt"

	"gioui.org/unit"
	"github.com/zodimo/go-compose/compose"
	"github.com/zodimo/go-compose/compose/foundation/layout/column"
	"github.com/zodimo/go-compose/compose/foundation/layout/spacer"
	m3text "github.com/zodimo/go-compose/compose/foundation/material3/text"
	"github.com/zodimo/go-compose/pkg/api"
)

// Define a CompositionLocal
var LocalString = compose.CompositionLocalOf(func() string { return "Default (Root)" })

func UI(c api.Composer) api.LayoutNode {
	// Top Level
	root := column.Column(
		func(c api.Composer) api.Composer {

			// 1. Display default value
			val1 := LocalString.Current(c)
			m3text.Text(fmt.Sprintf("1. Outer Value: %s", val1), m3text.TypestyleBodyLarge)(c)

			spacer.Height(int(unit.Dp(16)))(c)

			// 2. Provide a new value
			compose.CompositionLocalProvider(
				[]api.ProvidedValue{LocalString.Provides("Scoped Value")},
				func(c api.Composer) api.Composer {
					val2 := LocalString.Current(c)
					m3text.Text(fmt.Sprintf("2. Inner Value: %s", val2), m3text.TypestyleBodyLarge)(c)

					spacer.Height(int(unit.Dp(16)))(c)

					// 3. Nested Provider
					compose.CompositionLocalProvider(
						[]api.ProvidedValue{LocalString.Provides("Nested Scoped Value")},
						func(c api.Composer) api.Composer {
							val3 := LocalString.Current(c)
							m3text.Text(fmt.Sprintf("3. Nested Value: %s", val3), m3text.TypestyleBodyLarge)(c)
							return c
						},
					)(c)

					return c
				},
			)(c)

			spacer.Height(int(unit.Dp(16)))(c)

			// 4. Verify value is back to default
			val4 := LocalString.Current(c)
			m3text.Text(fmt.Sprintf("4. Outer Again: %s", val4), m3text.TypestyleBodyLarge)(c)

			return c
		},
	)

	return root(c).Build()
}
