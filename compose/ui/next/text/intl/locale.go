package intl

// TextDirection represents the direction of text flow.
type TextDirection string

const (
	// TextDirectionLTR is left-to-right text.
	TextDirectionLTR TextDirection = "ltr"
	// TextDirectionRTL is right-to-left text.
	TextDirectionRTL TextDirection = "rtl"
)

// Locale provides language information for the current system.
type Locale struct {
	// Language is the BCP-47 tag for the primary language of the system.
	Language string
	// Direction indicates the primary direction of text and layout
	// flow for the system.
	Direction TextDirection
}
