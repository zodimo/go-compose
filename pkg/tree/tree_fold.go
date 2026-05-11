package tree

//  fold in , fold out
// fold left, fold right
// depth first, breath first

// AKA FoldL or Depth first - > left to right
func FoldIn[T any](root Tree, initial T, operation func(carry T, item Tree) T) T {
	// return opertion(initial, root)
	untypedOperation := func(carry any, item Tree) any {
		carryT := carry.(T)
		return operation(carryT, item)
	}

	return root.FoldIn(initial, untypedOperation).(T)
}

// AKA FoldR or Beath first -> right to left
func FoldOut[T any](root Tree, initial T, operation func(item Tree, carry T) T) T {
	untypedOperation := func(item Tree, carry any) any {
		carryT := carry.(T)
		return operation(item, carryT)
	}
	return root.FoldOut(initial, untypedOperation).(T)
}
