module github.com/zodimo/go-compose

go 1.24.4

toolchain go1.24.11

require (
	gioui.org v0.9.0
	gioui.org/x v0.9.0
	git.sr.ht/~schnwalter/gio-mw v0.0.0-20250713180710-9d8d98474447
	github.com/go-text/typesetting v0.3.3
	github.com/zodimo/go-maybe v0.1.6
	github.com/zodimo/go-ternary v0.2.0
	github.com/zodimo/go-zero-hash v0.1.0
	golang.org/x/exp/shiny v0.0.0-20260112195511-716be5621a96
	golang.org/x/image v0.35.0
	golang.org/x/sync v0.19.0
	golang.org/x/text v0.33.0
	golang.org/x/tools v0.41.0
)

require (
	gioui.org/shader v1.0.8 // indirect
	git.wow.st/gmp/jni v0.0.0-20210610011705-34026c7e22d0 // indirect
	github.com/godbus/dbus/v5 v5.0.6 // indirect
	github.com/zodimo/go-lazy v0.1.1 // indirect
	golang.org/x/mod v0.32.0 // indirect
	golang.org/x/sys v0.41.0 // indirect
)

// required for android builds - contains the fix
// replace gioui.org => ../../development/zodimo-gio

// replace github.com/go-text/typesetting => ../../../GoProjects/clones/typesetting
replace github.com/go-text/typesetting => github.com/zodimo/typesetting v0.3.4-0.20260209162200-1565df70b998
