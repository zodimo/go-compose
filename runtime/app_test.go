package runtime

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/zodimo/go-compose/compose/ui/unit"
	"github.com/zodimo/go-compose/pkg/api"
)

func TestAppNoWindowsError(t *testing.T) {
	app := NewApp()
	err := app.Run()
	if err == nil {
		t.Fatal("expected error for no windows registered, got nil")
	}
	if err.Error() != "runtime.App: no windows registered" {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestAppExitOnLastWindowClose(t *testing.T) {
	app := NewApp()
	app.AddWindow(WindowOptions{Title: "Test Window"})

	// Inject a mock driver that simulates normal window close after a short delay
	app.driver = func(a *Application, cfg *windowConfig) {
		defer a.wg.Done()
		time.Sleep(10 * time.Millisecond)
	}

	err := app.Run()
	if err != nil {
		t.Fatalf("expected nil error on normal close, got: %v", err)
	}
}

func TestAppExitOnContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	app := NewApp(WithRootContext(ctx))
	app.AddWindow(WindowOptions{Title: "Test Window"})

	// Inject a mock driver that blocks until context is cancelled
	app.driver = func(a *Application, cfg *windowConfig) {
		defer a.wg.Done()
		<-a.ctx.Done()
	}

	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()

	err := app.Run()
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got: %v", err)
	}
}

func TestAppWindowErrorWins(t *testing.T) {
	app := NewApp()
	app.AddWindow(WindowOptions{Title: "Failing Window 1"})
	app.AddWindow(WindowOptions{Title: "Failing Window 2"})

	expectedErr := errors.New("gpu failure")

	var once sync.Once
	app.driver = func(a *Application, cfg *windowConfig) {
		defer a.wg.Done()
		time.Sleep(10 * time.Millisecond)
		once.Do(func() {
			select {
			case a.windowExit <- expectedErr:
			default:
			}
		})
	}

	err := app.Run()
	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected %v, got: %v", expectedErr, err)
	}
}

func TestAppQuitWhenLastWindowClosesFalse(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	app := NewApp(
		WithRootContext(ctx),
		QuitWhenLastWindowCloses(false),
	)
	app.AddWindow(WindowOptions{Title: "Test Window"})

	app.driver = func(a *Application, cfg *windowConfig) {
		defer a.wg.Done()
		// Closes immediately
	}

	// Should not exit immediately when window closes because quitOnLast is false
	done := make(chan error, 1)
	go func() {
		done <- app.Run()
	}()

	select {
	case <-done:
		t.Fatal("app.Run() returned prematurely when QuitWhenLastWindowCloses is false")
	case <-time.After(50 * time.Millisecond):
		// Expected: still waiting
	}

	// Now cancel context; app should exit with context.Canceled
	cancel()

	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("expected context.Canceled, got: %v", err)
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("app.Run() did not return after context cancellation")
	}
}

func TestAppConvenienceConstructor(t *testing.T) {
	dummyUI := func(c api.Composer) api.Composer { return c }
	app := App(
		Title("Single Window Demo"),
		Size(unit.NewDpSize(800, 600)),
		Content(dummyUI),
	)

	if len(app.windows) != 1 {
		t.Fatalf("expected 1 window, got %d", len(app.windows))
	}
	if app.windows[0].Title != "Single Window Demo" {
		t.Fatalf("expected Title 'Single Window Demo', got %s", app.windows[0].Title)
	}
	if app.windows[0].Size.Width != 800 || app.windows[0].Size.Height != 600 {
		t.Fatalf("unexpected size: %v", app.windows[0].Size)
	}
	if app.windows[0].Content == nil {
		t.Fatal("expected non-nil Content")
	}
}
