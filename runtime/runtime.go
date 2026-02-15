package runtime

import (
	"gioui.org/op"
	"github.com/zodimo/go-compose/pkg/api"
)

type Runtime interface {
	Run(LayoutContext, api.Composer, api.Composable) op.CallOp
}
