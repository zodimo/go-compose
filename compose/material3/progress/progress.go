package progress
// Engine-bound (Door 2): this component's public API is gioui-free; its
// implementation uses the engine directly behind the seam (internal/render).


import (
	"fmt"

	"github.com/zodimo/go-compose/compose/ui"
	"github.com/zodimo/go-compose/internal/layoutnode"

	"git.sr.ht/~schnwalter/gio-mw/widget/indicator"
)

const (
	Material3LinearProgressNodeID   = "Material3LinearProgress"
	Material3CircularProgressNodeID = "Material3CircularProgress"
)

func LinearProgressIndicator(progress float32, options ...IndicatorOption) Composable {
	// Initialize with a Linear indicator
	return indicatorComposable(indicator.Linear(), progress, Material3LinearProgressNodeID, options...)
}

func CircularProgressIndicator(progress float32, options ...IndicatorOption) Composable {
	// Initialize with a Circular indicator
	return indicatorComposable(indicator.Circular(), progress, Material3CircularProgressNodeID, options...)
}

func indicatorComposable(defaultIndicator *indicator.Indicator, progress float32, nodeID string, options ...IndicatorOption) Composable {
	return func(c Composer) Composer {
		opts := DefaultIndicatorOptions()
		for _, option := range options {
			if option == nil {
				continue
			}
			option(&opts)
		}

		if opts.Indicator == nil {
			key := c.GenerateID()
			path := c.GetPath()

			// We persist the indicator to maintain animation state (lastProgress, startTime)
			statePath := fmt.Sprintf("%d/%s/indicator", key, path)
			opts.Indicator = c.State(statePath, func() any { return defaultIndicator }).Get().(*indicator.Indicator)
		}

		// Update progress on every recomposition
		opts.Indicator.Progress = progress

		c.StartBlock(nodeID)
		c.Modifier(func(modifier ui.Modifier) ui.Modifier {
			return modifier.Then(opts.Modifier)
		})
		c.SetWidgetConstructor(indicatorWidgetConstructor(opts.Indicator))

		return c.EndBlock()
	}
}

func indicatorWidgetConstructor(ind *indicator.Indicator) layoutnode.LayoutNodeWidgetConstructor {
	return layoutnode.NewLayoutNodeWidgetConstructor(func(node layoutnode.LayoutNode) layoutnode.GioLayoutWidget {
	return func(gtx layoutnode.LayoutContext) layoutnode.LayoutDimensions {
		return layoutnode.FromGioDimensions(ind.Layout(*gtx.ToGio()))
	}
	})
}
