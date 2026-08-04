package apicheck

import (
	"golang.org/x/tools/go/packages"
)

// checkPackageImports enforces Rule B: a non-seam package may not directly
// import the engine packages (gioui.org/layout, gioui.org/op,
// gioui.org/widget, gioui.org/text, gioui.org/font). cmd/ packages and seam
// packages are excluded by the caller.
func checkPackageImports(pkg *packages.Package) []Violation {
	var out []Violation
	for _, imp := range pkg.Imports {
		for _, engine := range EnginePackages {
			if imp.PkgPath == engine {
				out = append(out, Violation{
					Kind:    KindImport,
					Pkg:     pkg.PkgPath,
					Symbol:  "",
					GioType: engine,
					Detail:  "engine import (confined to seam packages)",
				})
			}
		}
	}
	return out
}
