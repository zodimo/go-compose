package dialog

import (
	"fmt"
	"go-compose-dev/internal/layoutnode"

	"git.sr.ht/~schnwalter/gio-mw/wdk"
	"git.sr.ht/~schnwalter/gio-mw/widget/button"
	"git.sr.ht/~schnwalter/gio-mw/widget/dialog"
)

const Material3DialogNodeID = "Material3Dialog"

// AlertDialog creates a Material3 alert dialog.
func AlertDialog(
	onDismissRequest func(),
	onConfirm func(),
	confirmLabel string,
	options ...DialogOption,
) Composable {
	return func(c Composer) Composer {
		opts := DefaultDialogOptions()
		for _, option := range options {
			option(&opts)
		}

		// Persist the buttons so they maintain state (clicks/animations) across frames.
		key := c.GenerateID()
		path := c.GetPath()

		buttonsStatePath := fmt.Sprintf("%d/%s/dialog_buttons", key, path)
		buttonsValue := c.State(buttonsStatePath, func() any {
			return &DialogButtonsState{
				ConfirmButton: button.Text(),
				CancelButton:  button.Text(),
			}
		})
		buttonsState := buttonsValue.Get().(*DialogButtonsState)

		// Create constructor args with all necessary data
		constructorArgs := DialogConstructorArgs{
			ButtonsState:     buttonsState,
			Title:            opts.Title,
			Text:             opts.Text,
			ConfirmLabel:     confirmLabel,
			DismissLabel:     opts.DismissLabel,
			Icon:             opts.Icon,
			OnConfirm:        onConfirm,
			OnDismiss:        opts.OnDismiss,
			OnDismissRequest: onDismissRequest,
		}

		c.StartBlock(Material3DialogNodeID)
		c.Modifier(func(modifier Modifier) Modifier {
			return modifier.Then(opts.Modifier)
		})
		c.SetWidgetConstructor(dialogWidgetConstructor(constructorArgs))

		return c.EndBlock()
	}
}

type DialogButtonsState struct {
	ConfirmButton *button.Button
	CancelButton  *button.Button
}

type DialogConstructorArgs struct {
	ButtonsState     *DialogButtonsState
	Title            string
	Text             string
	ConfirmLabel     string
	DismissLabel     string
	Icon             wdk.IconWidget
	OnConfirm        func()
	OnDismiss        func()
	OnDismissRequest func()
}

func dialogWidgetConstructor(args DialogConstructorArgs) layoutnode.LayoutNodeWidgetConstructor {
	return layoutnode.NewLayoutNodeWidgetConstructor(func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {
		return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
			// 1. Create a fresh Dialog struct using the factory to get the correct Theme from gtx.
			//    We use Confirm(gtx) as the base since it supports both buttons.
			dlg := dialog.Confirm(gtx)

			// 2. Inject persisted buttons to maintain state
			dlg.ConfirmButton = args.ButtonsState.ConfirmButton
			dlg.CancelButton = args.ButtonsState.CancelButton

			// 3. Update Properties
			dlg.Headline = args.Title
			dlg.Label = args.Text
			dlg.ConfirmText = args.ConfirmLabel
			dlg.CancelText = args.DismissLabel

			if args.Icon != nil {
				dlg.WithIcon(args.Icon)
			}

			// 4. Handle Events using the persisted buttons (which are now in dlg)
			if dlg.ConfirmButton.Clicked(gtx) {
				if args.OnConfirm != nil {
					args.OnConfirm()
				}
			}

			// Only check cancel if label is present (implying it's shown)
			if args.DismissLabel != "" {
				if dlg.CancelButton.Clicked(gtx) {
					if args.OnDismiss != nil {
						args.OnDismiss()
					} else if args.OnDismissRequest != nil {
						args.OnDismissRequest()
					}
				}
			} else {
				// If no dismiss label, ensure CancelText is empty so it doesn't render?
				// gio-mw checks CancelText string.
				dlg.CancelText = ""
			}

			// 5. Layout
			// The Layout method will use the injected buttons and the theme from factory.
			return dlg.Layout(gtx)
		}
	})
}
