package button

// we need persistant storage for internal state

// func RegisterInput(node compose.Node, gtx layout.Context) {
// 	button := GetButton(node)
// 	clickable := GetClickable(node)
// 	onClick := GetOnClick(node)

// 	if button.WithClickable(clickable).Clicked(gtx) {
// 		onClick()
// 	}
// }

// func GetButton(node compose.Node) *button.Button {
// 	buttonAny, ok := node.Slots.Find(Material3ButtonButtonNodeID)
// 	if !ok {
// 		panic("material3 button button not found")
// 	}
// 	button, ok := buttonAny.(*button.Button)
// 	if !ok {
// 		panic("material3 button not found")
// 	}
// 	return button
// }

// func GetOnClick(node compose.Node) func() {
// 	onClickAny, ok := node.Slots.Find(Material3ButtonOnClickNodeID)
// 	if !ok {
// 		panic("material3 button onClick not found")
// 	}
// 	onClick, ok := onClickAny.(func())
// 	if !ok {
// 		panic("material3 button onClick is not a func()")
// 	}
// 	return onClick
// }

// func MeasureButton(node compose.Node, gtx layout.Context) layout.Dimensions {
// 	defer op.Record(gtx.Ops).Stop()

// 	button := GetButton(node)
// 	return button.Layout(gtx, GetLabelContent(node))
// }

// func GetLabelContent(node compose.Node) string {
// 	labelContent, ok := node.Slots.Find(Material3ButtonLabelNodeID)
// 	if !ok {
// 		panic("material3 button label not found")
// 	}
// 	labelContentString, ok := labelContent.(string)
// 	if !ok {
// 		panic("material3 button label is not a string")
// 	}
// 	return labelContentString
// }

// func GetClickable(node compose.Node) *widget.Clickable {
// 	clickableAny, ok := node.Slots.Find(Material3ButtonClickableNodeID)
// 	if !ok {
// 		panic("material3 button clickable not found")
// 	}
// 	clickable, ok := clickableAny.(*widget.Clickable)
// 	if !ok {
// 		panic("material3 button clickable is not a widget.Clickable")
// 	}
// 	return clickable
// }

// func GetButtonOptions(node compose.Node) ButtonOptions {
// 	buttonOptionsAny, ok := node.Slots.Find(MaterialButtonOptionsNodeID)
// 	if !ok {
// 		panic("material button options not found")
// 	}
// 	buttonOptions, ok := buttonOptionsAny.(ButtonOptions)
// 	if !ok {
// 		panic("material button options is not a ButtonOptions")
// 	}
// 	return buttonOptions
// }
