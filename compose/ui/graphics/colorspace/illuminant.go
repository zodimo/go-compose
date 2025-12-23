package colorspace

// Illuminant contains standard CIE WhitePoints.
var (
	// IlluminantA is Standard CIE 1931 2° illuminant A, encoded in xyY. This illuminant has a color temperature of 2856K.
	IlluminantA = NewWhitePoint(0.44757, 0.40745)

	// IlluminantB is Standard CIE 1931 2° illuminant B, encoded in xyY. This illuminant has a color temperature of 4874K.
	IlluminantB = NewWhitePoint(0.34842, 0.35161)

	// IlluminantC is Standard CIE 1931 2° illuminant C, encoded in xyY. This illuminant has a color temperature of 6774K.
	IlluminantC = NewWhitePoint(0.31006, 0.31616)

	// IlluminantD50 is Standard CIE 1931 2° illuminant D50, encoded in xyY. This illuminant has a color temperature of 5003K.
	IlluminantD50 = NewWhitePoint(0.34567, 0.35850)

	// IlluminantD55 is Standard CIE 1931 2° illuminant D55, encoded in xyY. This illuminant has a color temperature of 5503K.
	IlluminantD55 = NewWhitePoint(0.33242, 0.34743)

	// IlluminantD60 is Standard CIE 1931 2° illuminant D60, encoded in xyY. This illuminant has a color temperature of 6004K.
	IlluminantD60 = NewWhitePoint(0.32168, 0.33767)

	// IlluminantD65 is Standard CIE 1931 2° illuminant D65, encoded in xyY. This illuminant has a color temperature of 6504K.
	IlluminantD65 = NewWhitePoint(0.31271, 0.32902)

	// IlluminantD75 is Standard CIE 1931 2° illuminant D75, encoded in xyY. This illuminant has a color temperature of 7504K.
	IlluminantD75 = NewWhitePoint(0.29902, 0.31485)

	// IlluminantE is Standard CIE 1931 2° illuminant E, encoded in xyY. This illuminant has a color temperature of 5454K.
	IlluminantE = NewWhitePoint(0.33333, 0.33333)

	// IlluminantD50Xyz is the XYZ representation of D50.
	IlluminantD50Xyz = []float32{0.964212, 1.0, 0.825188}
)
