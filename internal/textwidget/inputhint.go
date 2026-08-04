package textwidget

import "gioui.org/io/key"

// InputHint wraps gioui key.InputHint, hiding it from the public API surface.
type InputHint uint8

const (
	HintAny       InputHint = InputHint(key.HintAny)
	HintText      InputHint = InputHint(key.HintText)
	HintNumeric   InputHint = InputHint(key.HintNumeric)
	HintEmail     InputHint = InputHint(key.HintEmail)
	HintURL       InputHint = InputHint(key.HintURL)
	HintTelephone InputHint = InputHint(key.HintTelephone)
	HintPassword  InputHint = InputHint(key.HintPassword)
)

// ToGio converts to the gioui key.InputHint value.
func (h InputHint) ToGio() key.InputHint {
	return key.InputHint(h)
}
