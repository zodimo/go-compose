package skiagraphics

import (
	"fmt"

	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-skia-support/skia/enums"
)

func graphicsPathOperationToSkiaPathOp(op graphics.PathOperation) enums.PathOp {
	switch op {
	case graphics.PathOperationDifference:
		return enums.PathOpDifference
	case graphics.PathOperationIntersect:
		return enums.PathOpIntersect
	case graphics.PathOperationUnion:
		return enums.PathOpUnion
	case graphics.PathOperationXor:
		return enums.PathOpXor
	case graphics.PathOperationReverseDifference:
		return enums.PathOpReverseDifference
	default:
		panic(fmt.Sprintf("Unknown PathOperation: %s", op.String()))
	}
}
