package runtime

import (
	"context"
	"errors"
	"fmt"
	"os"
	goRuntime "runtime"
	"sync"

	"github.com/zodimo/go-compose/compose/ui/unit"
	"github.com/zodimo/go-compose/pkg/api"
)

// windowConfig holds the internal configuration for a registered window.
type windowConfig struct {
	Title          string
	Size           unit.DpSize
	Content        api.Composable
	useLocalWindow bool
}

// windowDriverFunc drives a single window lifecycle.
type windowDriverFunc func(app *Application, cfg *windowConfig)

// Application is the framework-owned application entry. It owns the window event
// loop(s), the per-window store, and the frame render pipeline.
type Application struct {
	ctx              context.Context
	cancel           context.CancelFunc
	quitOnLast       bool
	windows          []*windowConfig
	windowExit       chan error
	allWindowsClosed chan struct{}
	wg               sync.WaitGroup
	driver           windowDriverFunc
	appMain          func()
}

// App creates a single-window application with the given options.
// If Title, Size, or Content are set in opts, a single window is registered.
func App(opts ...AppOption) *Application {
	cfg := DefaultAppConfig()
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}

	app := NewApp(
		WithRootContext(cfg.RootCtx),
		QuitWhenLastWindowCloses(cfg.QuitOnLast),
	)

	// If Content, Title, or Size was provided, register the primary window.
	if cfg.Content != nil || cfg.Title != "" || cfg.Size != unit.DpSizeZero {
		app.AddWindow(WindowOptions{
			Title:   cfg.Title,
			Size:    cfg.Size,
			Content: cfg.Content,
		})
	}

	return app
}

// NewApp creates a (potentially) multi-window application.
// Register windows with AddWindow, then call Run.
func NewApp(opts ...AppOption) *Application {
	cfg := DefaultAppConfig()
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}

	rootCtx := cfg.RootCtx
	if rootCtx == nil {
		rootCtx = context.Background()
	}
	ctx, cancel := context.WithCancel(rootCtx)

	a := &Application{
		ctx:              ctx,
		cancel:           cancel,
		quitOnLast:       cfg.QuitOnLast,
		windows:          make([]*windowConfig, 0),
		windowExit:       make(chan error, 1),
		allWindowsClosed: make(chan struct{}),
	}
	a.driver = (*Application).runWindow
	a.appMain = gioAppMain
	return a
}

// AddWindow registers a window to be opened when Run is called.
func (a *Application) AddWindow(opts WindowOptions) *Application {
	a.windows = append(a.windows, &windowConfig{
		Title:          opts.Title,
		Size:           opts.Size,
		Content:        opts.Content,
		useLocalWindow: true,
	})
	return a
}

// validate checks the configuration before starting windows or app.Main.
func (a *Application) validate() error {
	if len(a.windows) == 0 {
		return errors.New("runtime.App: no windows registered")
	}

	if len(a.windows) > 1 {
		goos := goRuntime.GOOS
		goarch := goRuntime.GOARCH
		if goos == "ios" || goos == "android" || goarch == "wasm" {
			return fmt.Errorf("runtime.App: multi-window not supported on %s/%s", goos, goarch)
		}
	}

	return nil
}

// Run blocks until the application closes. Returns nil when the last window
// closes normally, context.Canceled when the root context is cancelled, or
// the DestroyEvent.Err of a window that failed.
// Run must be called from the main goroutine.
func (a *Application) Run() error {
	if err := a.validate(); err != nil {
		return err
	}

	exitCh := make(chan error, 1)

	for _, w := range a.windows {
		a.wg.Add(1)
		go a.driver(a, w)
	}

	go func() {
		a.wg.Wait()
		close(a.allWindowsClosed)
	}()

	// Exit monitor: picks the first terminal event.
	allClosed := a.allWindowsClosed
	go func() {
		for {
			select {
			case err := <-a.windowExit:
				exitCh <- err
				return
			case <-a.ctx.Done():
				exitCh <- a.ctx.Err()
				return
			case <-allClosed:
				if a.quitOnLast {
					exitCh <- nil
					return
				}
				allClosed = nil
			}
		}
	}()

	if goRuntime.GOOS == "darwin" {
		go func() {
			err := <-exitCh
			osExit(exitCode(err))
		}()
		if a.appMain != nil {
			a.appMain()
		}
		return nil
	}

	return <-exitCh
}

func exitCode(err error) int {
	switch {
	case err == nil || errors.Is(err, context.Canceled):
		return 0
	default:
		return 1
	}
}

var osExit = os.Exit
