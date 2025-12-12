package colorhelper

import "image/color"

func ToNRGBA(input color.Color) color.NRGBA {
	nrgbaModel := color.NRGBAModel
	return nrgbaModel.Convert(input).(color.NRGBA)
}
