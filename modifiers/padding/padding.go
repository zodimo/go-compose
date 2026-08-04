package padding

import "github.com/zodimo/go-compose/internal/layoutnode"

const NotSet = -1

// TextDirection represents text direction for RTL awareness.
type TextDirection layoutnode.LayoutDirection

const (
	// RTL indicates right-to-left text direction.
	RTL TextDirection = TextDirection(layoutnode.LayoutDirectionRTL)
)

type PaddingData struct {
	Start    int
	Top      int
	End      int
	Bottom   int
	RtlAware bool // future proofing for RTL support
}

func DefaultPadding() PaddingData {
	return PaddingData{
		Start:    NotSet,
		Top:      NotSet,
		End:      NotSet,
		Bottom:   NotSet,
		RtlAware: false,
	}
}
