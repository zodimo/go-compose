// Package apicheck implements an API-purity analyzer for go-compose.
//
// It enforces two contracts from the "abstract-render-backend" change:
//
//   Rule A (signature purity, spec: public-api-purity): no exported symbol in
//   the public package trees (compose/, modifiers/, theme/, runtime/, pkg/)
//   references a type from any gioui.org package in its signature. A signature
//   includes function parameters, return values, struct fields, embedded
//   types, type aliases, declared/underlying types of defined types, method
//   signatures, and const/var types.
//
//   Rule B (import restriction, spec: render-backend-seam): only the seam
//   packages (internal/layoutnode, internal/render, and their subpackages) may
//   directly import the engine packages gioui.org/layout, gioui.org/op,
//   gioui.org/widget, gioui.org/text, and gioui.org/font. cmd/ packages are
//   excluded because they are application shells (app glue), not framework
//   API.
//
// Rule A is always enforced (this is what `make check-api` gates on: it exits
// non-zero when the public API is not gioui-free). Rule B is enforced when
// Options.All is set; the Makefile enables it once the framework no longer
// imports engine packages (Phase 3 of the change).
package apicheck

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"
)

// gioPrefix is the import path prefix shared by every gioui.org package.
const gioPrefix = "gioui.org"

// modulePath is the go module this analyzer guards.
const modulePath = "github.com/zodimo/go-compose"

// EnginePackages are the engine packages whose direct import is confined to
// the seam packages (spec: render-backend-seam — "Seam packages are the only
// engine-type importers").
var EnginePackages = []string{
	"gioui.org/layout",
	"gioui.org/op",
	"gioui.org/widget",
	"gioui.org/text",
	"gioui.org/font",
}

// PublicTrees are the public package trees checked for signature purity
// (spec: api-purity-analyzer, public-api-purity).
var PublicTrees = []string{
	"./compose/...",
	"./modifiers/...",
	"./theme/...",
	"./runtime/...",
	"./pkg/...",
}

// IsSeam reports whether pkgPath is a seam package (internal/layoutnode,
// internal/render, or a subpackage thereof). Seam packages are exempt from the
// import restriction.
func IsSeam(pkgPath string) bool {
	return isWithin(pkgPath, modulePath+"/internal/layoutnode") ||
		isWithin(pkgPath, modulePath+"/internal/render")
}

// isPublicTree reports whether pkgPath belongs to one of the public package
// trees.
func isPublicTree(pkgPath string) bool {
	for _, tree := range []string{
		modulePath + "/compose",
		modulePath + "/modifiers",
		modulePath + "/theme",
		modulePath + "/runtime",
		modulePath + "/pkg",
	} {
		if isWithin(pkgPath, tree) {
			return true
		}
	}
	return false
}

// isWithin reports whether pkgPath equals base or is a subpackage of base.
func isWithin(pkgPath, base string) bool {
	return pkgPath == base || strings.HasPrefix(pkgPath, base+"/")
}

// Kind classifies a violation.
type Kind string

const (
	// KindSignature is a leaked gioui.org type in an exported signature.
	KindSignature Kind = "signature"
	// KindImport is a direct import of an engine package by a non-seam package.
	KindImport Kind = "import"
)

// Violation is a single purity violation. It names the package, the exported
// symbol, and the gioui.org type that leaks, so an implementer can locate and
// fix it without further investigation (spec: api-purity-analyzer —
// "Analyzer reports offending symbols").
type Violation struct {
	Kind    Kind
	Pkg     string
	Symbol  string
	GioType string
	Detail  string
}

// String renders the violation in an actionable, greppable form.
func (v Violation) String() string {
	switch v.Kind {
	case KindSignature:
		return fmt.Sprintf("%s: %s: exported %q references %s (%s)",
			v.Pkg, v.Kind, v.Symbol, v.GioType, v.Detail)
	case KindImport:
		return fmt.Sprintf("%s: %s: imports engine package %s (confined to seam)",
			v.Pkg, v.Kind, v.GioType)
	default:
		return fmt.Sprintf("%s: %s: %s (%s)", v.Pkg, v.Kind, v.Symbol, v.Detail)
	}
}

// Options controls the analyzer run.
type Options struct {
	// All additionally enforces Rule B (the engine-import restriction) and
	// extends the signature rule to every non-seam, non-cmd package. The
	// Makefile enables it once Phase 3 removes engine imports from non-seam
	// packages.
	All bool
}

// loadMode is the go/packages mode needed for both rules.
const loadMode = packages.NeedName |
	packages.NeedFiles |
	packages.NeedCompiledGoFiles |
	packages.NeedImports |
	packages.NeedDeps |
	packages.NeedTypes |
	packages.NeedSyntax |
	packages.NeedTypesInfo

// Check loads the public package trees (or the whole module when opts.All is
// set) and returns every purity violation found, sorted for stable output.
func Check(opts Options) ([]Violation, error) {
	root, err := moduleRoot()
	if err != nil {
		return nil, err
	}

	patterns := PublicTrees
	if opts.All {
		patterns = []string{"./..."}
	}

	cfg := &packages.Config{
		Mode: loadMode,
		Dir:  root,
	}
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		return nil, fmt.Errorf("load packages: %w", err)
	}
	if n := packages.PrintErrors(pkgs); n > 0 {
		return nil, fmt.Errorf("loading public packages produced %d error(s); refusing to report on a broken tree", n)
	}

	var violations []Violation
	for _, pkg := range pkgs {
		if IsSeam(pkg.PkgPath) {
			continue
		}
		if pkg.Types == nil || pkg.TypesInfo == nil {
			continue
		}
		if isPublicTree(pkg.PkgPath) || opts.All {
			if !opts.All || !isCmdPackage(pkg.PkgPath) {
				violations = append(violations, checkPackageSignatures(pkg)...)
			}
		}
		if opts.All && !isCmdPackage(pkg.PkgPath) {
			violations = append(violations, checkPackageImports(pkg)...)
		}
	}

	// Deduplicate and sort for deterministic output.
	seen := map[Violation]bool{}
	unique := violations[:0]
	for _, v := range violations {
		if !seen[v] {
			seen[v] = true
			unique = append(unique, v)
		}
	}
	sort.Slice(unique, func(i, j int) bool {
		a, b := unique[i], unique[j]
		if a.Pkg != b.Pkg {
			return a.Pkg < b.Pkg
		}
		if a.Kind != b.Kind {
			return a.Kind < b.Kind
		}
		if a.Symbol != b.Symbol {
			return a.Symbol < b.Symbol
		}
		return a.GioType < b.GioType
	})
	return unique, nil
}

// isCmdPackage reports whether pkgPath is an application shell under cmd/.
// These are app glue (demos, CLI) and are excluded from both rules; the design
// explicitly permits gioui app imports there.
func isCmdPackage(pkgPath string) bool {
	return isWithin(pkgPath, modulePath+"/cmd")
}

// moduleRoot walks up from the working directory to find the module root
// (the directory containing go.mod). Relative load patterns are resolved
// against it so the analyzer works from any working directory (e.g. under
// `go test ./internal/api-check/...`).
func moduleRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found above %s", dir)
		}
		dir = parent
	}
}
