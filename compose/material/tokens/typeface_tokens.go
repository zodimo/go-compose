package tokens

import "github.com/zodimo/go-compose/compose/ui/text/font"

var TypefaceTokens = TypefaceTokensData{
	Brand:         font.FontFamilySansSerif,
	Plain:         font.FontFamilySansSerif,
	WeightBold:    font.FontWeightBold,
	WeightMedium:  font.FontWeightMedium,
	WeightRegular: font.FontWeightNormal,
}

type TypefaceTokensData struct {
	Brand         font.FontFamily
	Plain         font.FontFamily
	WeightBold    font.FontWeight
	WeightMedium  font.FontWeight
	WeightRegular font.FontWeight
}
