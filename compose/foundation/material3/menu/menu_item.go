package menu

import (
	"go-compose-dev/compose/foundation/layout/row"
	"go-compose-dev/internal/modifiers/clickable"
	"go-compose-dev/internal/modifiers/padding"
	"go-compose-dev/internal/modifiers/size"
	"go-compose-dev/pkg/api"
)

// DropdownMenuItem Composable
func DropdownMenuItem(
	onClick func(),
	content api.Composable,
) api.Composable {
	return func(c api.Composer) api.Composer {
		return row.Row(
			content,
			row.WithModifier(clickable.OnClick(onClick)),
			row.WithModifier(padding.Horizontal(16, 16)),
			row.WithModifier(padding.Vertical(8, 8)),
			row.WithModifier(size.FillMaxWidth()),
			row.WithModifier(size.Height(48)),
			row.WithAlignment(row.Middle),
		)(c)
	}
}
