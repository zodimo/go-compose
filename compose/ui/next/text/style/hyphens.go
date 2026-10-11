package style

import (
	"github.com/zodimo/go-sentinel-helper/sentinel/companions"
)

// Automatic hyphenation configuration.
//
// Hyphenation is a dash-like punctuation mark used to join two-words into one or separate
// syl-lab-les of a word.
//
// Automatic hyphenation is added between syllables at appropriate hyphenation points, following
// language rules.
//
// However, user can override automatic break point selection, suggesting line break opportunities
// (see Suggesting line break opportunities below).
//
// Suggesting line break opportunities:
//   - "\u2010" ("hard" hyphen) Indicates a visible line break opportunity. Even if the
//     line is not actually broken at that point, the hyphen is still rendered.
//   - "\u00AD" ("soft" hyphen) This character is not rendered visibly; instead, it marks a
//     place where the word can be broken if hyphenation is necessary.
//
// The default configuration for [Hyphens] = [HyphensNone]
type Hyphens int

const (
	// hyphensUnspecified represents an unset value, a usual replacement for "null" when a
	// primitive value is desired. It is deliberately package-private: the only public
	// way to obtain it is HyphensCompanion.Unspecified().
	hyphensUnspecified Hyphens = 0

	// Lines will break with no hyphenation.
	//
	// "Hard" hyphens will still be respected. However, no automatic hyphenation will be
	// attempted. If a word must be broken due to being longer than a line, it will break at any
	// character and will not attempt to break at a syllable boundary.
	//
	// +---------+
	// | Experim |
	// | ental   |
	// +---------+
	HyphensNone Hyphens = 1

	// The words will be automatically broken at appropriate hyphenation points.
	//
	// However, suggested line break opportunities (see Suggesting line break opportunities
	// above) will override automatic break point selection when present.
	//
	// +---------+
	// | Experi- |
	// | mental  |
	// +---------+
	HyphensAuto Hyphens = 2
)

func (h Hyphens) String() string {
	switch h {
	case HyphensNone:
		return "Hyphens.None"
	case HyphensAuto:
		return "Hyphens.Auto"
	case hyphensUnspecified:
		return "Hyphens.Unspecified"
	default:
		return "Invalid"
	}
}

// Returns true if it is not [HyphensUnspecified].
func (h Hyphens) IsSpecified() bool {
	return isSpecifiedHyphens(h)
}

// === Sentinel operations ===

func isSpecifiedHyphens(h Hyphens) bool {
	return h != hyphensUnspecified
}

func isUnspecifiedHyphens(h Hyphens) bool {
	return h == hyphensUnspecified
}

func takeOrElseHyphens(h, def Hyphens) Hyphens {
	if isSpecifiedHyphens(h) {
		return h
	}
	return def
}

func mergeHyphens(a, b Hyphens) Hyphens {
	if isUnspecifiedHyphens(b) {
		return a
	}
	return b
}

func stringHyphens(h Hyphens) string {
	return h.String()
}

func sameHyphens(a, b Hyphens) bool {
	return a == b
}

func semanticEqualHyphens(a, b Hyphens) bool {
	return a == b
}

func equalHyphens(a, b Hyphens) bool {
	return a == b
}

func copyHyphens(h Hyphens) Hyphens {
	return h
}

// HyphensCompanion is the single exported sentinel-family symbol for Hyphens.
var HyphensCompanion = companions.NewScalarCompanion(
	func() Hyphens { return hyphensUnspecified },
	stringHyphens,
	isSpecifiedHyphens,
	isUnspecifiedHyphens,
	takeOrElseHyphens,
	mergeHyphens,
	sameHyphens,
	semanticEqualHyphens,
	equalHyphens,
	copyHyphens,
)
