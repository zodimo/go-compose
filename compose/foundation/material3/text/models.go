package text

import "git.sr.ht/~schnwalter/gio-mw/token"

// Typestyle aliases token.Typestyle to avoid direct dependency on gio-mw for consumers.
type Typestyle = token.Typestyle

const (
	TypestyleDefault Typestyle = token.TypestyleDefault

	TypestyleLabelSmall     = token.TypestyleLabelSmall
	TypestyleLabelMedium    = token.TypestyleLabelMedium
	TypestyleLabelLarge     = token.TypestyleLabelLarge
	TypestyleBodySmall      = token.TypestyleBodySmall
	TypestyleBodyMedium     = token.TypestyleBodyMedium
	TypestyleBodyLarge      = token.TypestyleBodyLarge
	TypestyleTitleSmall     = token.TypestyleTitleSmall
	TypestyleTitleMedium    = token.TypestyleTitleMedium
	TypestyleTitleLarge     = token.TypestyleTitleLarge
	TypestyleHeadlineSmall  = token.TypestyleHeadlineSmall
	TypestyleHeadlineMedium = token.TypestyleHeadlineMedium
	TypestyleHeadlineLarge  = token.TypestyleHeadlineLarge
	TypestyleDisplaySmall   = token.TypestyleDisplaySmall
	TypestyleDisplayMedium  = token.TypestyleDisplayMedium
	TypestyleDisplayLarge   = token.TypestyleDisplayLarge

	TypestyleLabelSmallEmphasized     = token.TypestyleLabelSmallEmphasized
	TypestyleLabelMediumEmphasized    = token.TypestyleLabelMediumEmphasized
	TypestyleLabelLargeEmphasized     = token.TypestyleLabelLargeEmphasized
	TypestyleBodySmallEmphasized      = token.TypestyleBodySmallEmphasized
	TypestyleBodyMediumEmphasized     = token.TypestyleBodyMediumEmphasized
	TypestyleBodyLargeEmphasized      = token.TypestyleBodyLargeEmphasized
	TypestyleTitleSmallEmphasized     = token.TypestyleTitleSmallEmphasized
	TypestyleTitleMediumEmphasized    = token.TypestyleTitleMediumEmphasized
	TypestyleTitleLargeEmphasized     = token.TypestyleTitleLargeEmphasized
	TypestyleHeadlineSmallEmphasized  = token.TypestyleHeadlineSmallEmphasized
	TypestyleHeadlineMediumEmphasized = token.TypestyleHeadlineMediumEmphasized
	TypestyleHeadlineLargeEmphasized  = token.TypestyleHeadlineLargeEmphasized
	TypestyleDisplaySmallEmphasized   = token.TypestyleDisplaySmallEmphasized
	TypestyleDisplayMediumEmphasized  = token.TypestyleDisplayMediumEmphasized
	TypestyleDisplayLargeEmphasized   = token.TypestyleDisplayLargeEmphasized
)
