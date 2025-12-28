package bottomsheet

import (
	"github.com/zodimo/go-compose/pkg/api"

	"git.sr.ht/~schnwalter/gio-mw/token"
	"github.com/zodimo/go-compose/theme"
)

type Composable = api.Composable

type ModalBottomSheetOptions struct {
	IsOpen           bool // Controlled by parent usually, or we can use visible state?
	OnDismissRequest func()
	SheetState       *SheetState
	ContainerColor   theme.ColorDescriptor // Will use default if not set
	ScrimColor       theme.ColorDescriptor // Will use default if not set
	Shape            token.CornerShape     // Will use default if not set
	DragHandle       Composable            // Optional custom drag handle
	// WindowInsets     column.WindowInsets // For handling safe areas if needed - Removed for compilation
}

type ModalBottomSheetOption func(*ModalBottomSheetOptions)

func DefaultModalBottomSheetOptions() ModalBottomSheetOptions {
	return ModalBottomSheetOptions{
		IsOpen:         false,
		ContainerColor: theme.ColorHelper.UnspecifiedColor(),
		ScrimColor:     theme.ColorHelper.UnspecifiedColor(),
	}
}

func WithSheetState(state *SheetState) ModalBottomSheetOption {
	return func(o *ModalBottomSheetOptions) {
		o.SheetState = state
	}
}

// WithIsOpen is useful if the parent controls the state specifically without a SheetState object,
// but usually SheetState is preferred for imperative show/hide.
// Let's align with Drawer: it uses `IsOpen` and `OnClose`.
func WithIsOpen(isOpen bool) ModalBottomSheetOption {
	return func(o *ModalBottomSheetOptions) {
		o.IsOpen = isOpen
	}
}

func WithOnDismissRequest(onDismiss func()) ModalBottomSheetOption {
	return func(o *ModalBottomSheetOptions) {
		o.OnDismissRequest = onDismiss
	}
}

func WithContainerColor(color theme.ColorDescriptor) ModalBottomSheetOption {
	return func(o *ModalBottomSheetOptions) {
		o.ContainerColor = color
	}
}

func WithScrimColor(color theme.ColorDescriptor) ModalBottomSheetOption {
	return func(o *ModalBottomSheetOptions) {
		o.ScrimColor = color
	}
}

// Additional options for Shape, DragHandle, etc.
