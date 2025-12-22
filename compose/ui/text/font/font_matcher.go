package font

// FontMatcher matches fonts from a FontFamily to a weight/style request.
// It applies the rules at CSS 4 Font Matching:
// https://www.w3.org/TR/css-fonts-4/#font-style-matching
type FontMatcher struct{}

// NewFontMatcher creates a new FontMatcher.
func NewFontMatcher() *FontMatcher {
	return &FontMatcher{}
}

// MatchFont returns fonts matching the requested weight and style.
// If there is not a font that exactly satisfies the given constraints,
// the best match will be returned following CSS 4 Font Matching rules.
// Returns an empty slice if no fonts match.
func (m *FontMatcher) MatchFont(fonts []Font, weight FontWeight, style FontStyle) []Font {
	if len(fonts) == 0 {
		return nil
	}

	// Check for exact match first
	exactMatches := filterFonts(fonts, func(f Font) bool {
		return f.Weight().Equals(weight) && f.Style() == style
	})
	if len(exactMatches) > 0 {
		return exactMatches
	}

	// If no exact match, filter by style first
	fontsToSearch := filterFonts(fonts, func(f Font) bool {
		return f.Style() == style
	})
	if len(fontsToSearch) == 0 {
		fontsToSearch = fonts
	}

	var result []Font
	switch {
	case weight.Compare(FontWeightW400) < 0:
		// If the desired weight is less than 400:
		// - weights less than or equal to the desired weight are checked in descending order
		// - followed by weights above the desired weight in ascending order
		result = m.filterByClosestWeight(fontsToSearch, weight, true, nil, nil)

	case weight.Compare(FontWeightW500) > 0:
		// If the desired weight is greater than 500:
		// - weights greater than or equal to the desired weight are checked in ascending order
		// - followed by weights below the desired weight in descending order
		result = m.filterByClosestWeight(fontsToSearch, weight, false, nil, nil)

	default:
		// If the desired weight is inclusively between 400 and 500:
		// - weights greater than or equal to the target weight are checked in ascending order until 500
		// - followed by weights less than the target weight in descending order
		// - followed by weights greater than 500
		result = m.filterByClosestWeight(fontsToSearch, weight, false, nil, &FontWeightW500)
		if len(result) == 0 {
			result = m.filterByClosestWeight(fontsToSearch, weight, true, nil, nil)
		}
		if len(result) == 0 {
			result = m.filterByClosestWeight(fontsToSearch, weight, false, &FontWeightW500, nil)
		}
	}

	return result
}

// MatchFontFromFamily matches fonts from a FontListFontFamily.
func (m *FontMatcher) MatchFontFromFamily(family *FontListFontFamily, weight FontWeight, style FontStyle) []Font {
	return m.MatchFont(family.Fonts, weight, style)
}

// filterByClosestWeight finds the closest weight match.
// preferBelow: if true, prefer weights below the target; otherwise prefer above.
// minSearchRange: if not nil, exclude weights below this value.
// maxSearchRange: if not nil, exclude weights above this value.
func (m *FontMatcher) filterByClosestWeight(
	fonts []Font,
	weight FontWeight,
	preferBelow bool,
	minSearchRange *FontWeight,
	maxSearchRange *FontWeight,
) []Font {
	var bestWeightAbove *FontWeight
	var bestWeightBelow *FontWeight

	for _, font := range fonts {
		possibleWeight := font.Weight()

		// Apply range filters
		if minSearchRange != nil && possibleWeight.Compare(*minSearchRange) < 0 {
			continue
		}
		if maxSearchRange != nil && possibleWeight.Compare(*maxSearchRange) > 0 {
			continue
		}

		cmp := possibleWeight.Compare(weight)
		if cmp < 0 {
			// possibleWeight < target weight
			if bestWeightBelow == nil || possibleWeight.Compare(*bestWeightBelow) > 0 {
				w := possibleWeight
				bestWeightBelow = &w
			}
		} else if cmp > 0 {
			// possibleWeight > target weight
			if bestWeightAbove == nil || possibleWeight.Compare(*bestWeightAbove) < 0 {
				w := possibleWeight
				bestWeightAbove = &w
			}
		} else {
			// Exact weight match
			w := possibleWeight
			bestWeightAbove = &w
			bestWeightBelow = &w
			break
		}
	}

	var bestWeight *FontWeight
	if preferBelow {
		if bestWeightBelow != nil {
			bestWeight = bestWeightBelow
		} else {
			bestWeight = bestWeightAbove
		}
	} else {
		if bestWeightAbove != nil {
			bestWeight = bestWeightAbove
		} else {
			bestWeight = bestWeightBelow
		}
	}

	if bestWeight == nil {
		return nil
	}

	return filterFonts(fonts, func(f Font) bool {
		return f.Weight().Equals(*bestWeight)
	})
}

// filterFonts returns fonts that match the predicate.
func filterFonts(fonts []Font, predicate func(Font) bool) []Font {
	var result []Font
	for _, f := range fonts {
		if predicate(f) {
			result = append(result, f)
		}
	}
	return result
}
