module github.com/zodimo/go-compose

go 1.27.0

require (
	gioui.org v0.10.0
	gioui.org/x v0.10.0
	git.sr.ht/~schnwalter/gio-mw v0.0.0-20260221053317-be0445f63b48
	github.com/go-text/typesetting v0.3.4
	github.com/zodimo/go-maybe v0.1.9
	github.com/zodimo/go-ternary v0.2.0
	github.com/zodimo/go-zero-hash v0.1.0
	golang.org/x/exp/shiny v0.0.0-20260611194520-c48552f49976
	golang.org/x/image v0.42.0
	golang.org/x/sync v0.21.0
	golang.org/x/text v0.38.0
	golang.org/x/tools v0.46.0
)

require (
	gioui.org/shader v1.0.8 // indirect
	git.wow.st/gmp/jni v0.0.0-20260127013417-d142949d346a // indirect
	github.com/godbus/dbus/v5 v5.2.2 // indirect
	github.com/zodimo/go-lazy v0.1.1 // indirect
	golang.org/x/mod v0.37.0 // indirect
	golang.org/x/net v0.56.0 // indirect
	golang.org/x/sys v0.46.0 // indirect
)

// required for android builds - contains the fix
// replace gioui.org => ../../development/zodimo-gio

// replace github.com/go-text/typesetting => ../../../GoProjects/clones/typesetting

// replace github.com/go-text/typesetting => github.com/zodimo/typesetting v0.3.4-0.20260209162200-1565df70b998

// replace github.com/go-text/typesetting => ../typesetting
