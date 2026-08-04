package apicheck

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/packages"
)

// checkPackageSignatures inspects every exported symbol of pkg and reports
// those whose signature references a gioui.org type.
//
// The checker combines go/types walking (params, returns, fields, methods,
// const/var types, generic type arguments) with an AST pass over type
// declarations. The AST pass is required to catch declared types whose RHS is
// a gioui named type — go/types resolves `type TextAlign gioText.Alignment`
// such that TextAlign's underlying type is int, losing the reference to
// gioText.Alignment entirely (spec: api-purity-analyzer — "Analyzer detects
// gioui underlying types").
func checkPackageSignatures(pkg *packages.Package) []Violation {
	var out []Violation
	scope := pkg.Types.Scope()
	for _, name := range scope.Names() {
		if !token.IsExported(name) {
			continue
		}
		obj := scope.Lookup(name)
		switch o := obj.(type) {
		case *types.TypeName:
			out = append(out, checkTypeName(pkg, o)...)
		case *types.Func:
			out = append(out, checkSignature(pkg, name, o.Type(), "function")...)
		case *types.Var:
			if gio, ok := findGio(o.Type()); ok {
				out = append(out, Violation{
					Kind:    KindSignature,
					Pkg:     pkg.PkgPath,
					Symbol:  name,
					GioType: gio,
					Detail:  "variable type",
				})
			}
		case *types.Const:
			if gio, ok := findGio(o.Type()); ok {
				out = append(out, Violation{
					Kind:    KindSignature,
					Pkg:     pkg.PkgPath,
					Symbol:  name,
					GioType: gio,
					Detail:  "constant type",
				})
			}
		}
	}
	return out
}

// checkTypeName inspects an exported type declaration: its declared type
// expression (AST), its struct fields and embedded types, and its method set.
func checkTypeName(pkg *packages.Package, o *types.TypeName) []Violation {
	name := o.Name()
	var out []Violation

	// 1. AST pass: the declared type expression. Catches both aliases
	// (`type GioImage = widget.Image`) and defined types whose RHS references
	// gioui (`type TextAlign gioText.Alignment`, `type Locale = system.Locale`).
	for _, file := range pkg.Syntax {
		for _, decl := range file.Decls {
			gd, ok := decl.(*ast.GenDecl)
			if !ok || gd.Tok != token.TYPE {
				continue
			}
			for _, spec := range gd.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok || ts.Name.Name != name {
					continue
				}
				tv := pkg.TypesInfo.Types[ts.Type]
				if tv.Type != nil {
					if gio, ok := findGio(tv.Type); ok {
						what := "defined type"
						if ts.Assign != token.NoPos {
							what = "type alias"
						}
						out = append(out, Violation{
							Kind:    KindSignature,
							Pkg:     pkg.PkgPath,
							Symbol:  name,
							GioType: gio,
							Detail:  what,
						})
					}
				}
			}
		}
	}

	// 2. types pass: struct fields (exported and embedded), methods, and
	// interface methods. Aliases inherit the target type's checks (done at the
	// target's own definition site, or exempt for seam types), so aliases and
	// types that resolve through an alias chain are skipped here.
	if types.Unalias(o.Type()) != o.Type() {
		return out
	}
	named, ok := o.Type().(*types.Named)
	if !ok {
		return out
	}

	// Struct fields: a field whose type is gioui leaks the engine into the
	// exported type's layout (e.g. TextFieldWidget embedding widget.Editor).
	if st, ok := named.Underlying().(*types.Struct); ok {
		for i := 0; i < st.NumFields(); i++ {
			f := st.Field(i)
			if gio, ok := findGio(f.Type()); ok {
				what := "struct field " + f.Name()
				if f.Anonymous() {
					what = "embedded type"
				}
				out = append(out, Violation{
					Kind:    KindSignature,
					Pkg:     pkg.PkgPath,
					Symbol:  name,
					GioType: gio,
					Detail:  what,
				})
			}
		}
	}

	// Methods on the type (value and pointer receivers). Exported methods with
	// gioui params/results leak through the type's callable surface.
	for i := 0; i < named.NumMethods(); i++ {
		m := named.Method(i)
		if !token.IsExported(m.Name()) {
			continue
		}
		out = append(out, checkSignature(pkg, name+"."+m.Name(), m.Type(), "method")...)
	}

	// Interface methods (including embedded ones flattened by go/types).
	if iface, ok := named.Underlying().(*types.Interface); ok {
		for i := 0; i < iface.NumMethods(); i++ {
			m := iface.Method(i)
			if !token.IsExported(m.Name()) {
				continue
			}
			out = append(out, checkSignature(pkg, name+"."+m.Name(), m.Type(), "interface method")...)
		}
	}

	return out
}

// checkSignature inspects a function/method signature's parameters and results
// for gioui.org types.
func checkSignature(pkg *packages.Package, symbol string, typ types.Type, what string) []Violation {
	sig, ok := typ.(*types.Signature)
	if !ok {
		return nil
	}
	var out []Violation
	for i := 0; i < sig.Params().Len(); i++ {
		p := sig.Params().At(i)
		if gio, ok := findGio(p.Type()); ok {
			out = append(out, Violation{
				Kind:    KindSignature,
				Pkg:     pkg.PkgPath,
				Symbol:  symbol,
				GioType: gio,
				Detail:  what + " parameter " + p.Name(),
			})
		}
	}
	for i := 0; i < sig.Results().Len(); i++ {
		r := sig.Results().At(i)
		if gio, ok := findGio(r.Type()); ok {
			out = append(out, Violation{
				Kind:    KindSignature,
				Pkg:     pkg.PkgPath,
				Symbol:  symbol,
				GioType: gio,
				Detail:  what + " result",
			})
		}
	}
	return out
}

// findGio walks a type structure and returns the first gioui.org type found
// (as "gioui.org/<path>.<Name>"), or ("", false) if the type is gioui-free.
//
// The walk follows pointers, slices, arrays, maps, channels, signatures,
// struct fields, embedded interfaces, generic type arguments, and underlying
// types — so a gioui type hidden inside a composite type is still caught. The
// walk does NOT descend into interface method sets (those are checked
// separately by checkTypeName) to avoid double-reporting.
func findGio(t types.Type) (string, bool) {
	seen := map[types.Type]bool{}
	var walk func(t types.Type) (string, bool)
	walk = func(t types.Type) (string, bool) {
		if t == nil || seen[t] {
			return "", false
		}
		seen[t] = true

		switch tt := t.(type) {
		case *types.Alias:
			return walk(types.Unalias(tt))
		case *types.Named:
			if obj := tt.Obj(); obj != nil && obj.Pkg() != nil && strings.HasPrefix(obj.Pkg().Path(), gioPrefix) {
				return obj.Pkg().Path() + "." + obj.Name(), true
			}
			if ta := tt.TypeArgs(); ta != nil {
				for i := 0; i < ta.Len(); i++ {
					if p, ok := walk(ta.At(i)); ok {
						return p, true
					}
				}
			}
			// Do NOT recurse into the underlying type of a named type. A defined
			// type whose declared RHS is a go-compose type (e.g. the seam's
			// LayoutContext wrapping an engine pointer in an unexported field) is
			// clean at the signature level; its internals are checked at the type's
			// own definition site (exported structs are field-walked there) or are
			// exempt (seam packages). Recursing would false-positive on seam-defined
			// wrapper types re-exported by public packages.
			return "", false
		case *types.Pointer:
			return walk(tt.Elem())
		case *types.Slice:
			return walk(tt.Elem())
		case *types.Array:
			return walk(tt.Elem())
		case *types.Map:
			if p, ok := walk(tt.Key()); ok {
				return p, true
			}
			return walk(tt.Elem())
		case *types.Chan:
			return walk(tt.Elem())
		case *types.Signature:
			for i := 0; i < tt.Params().Len(); i++ {
				if p, ok := walk(tt.Params().At(i).Type()); ok {
					return p, true
				}
			}
			for i := 0; i < tt.Results().Len(); i++ {
				if p, ok := walk(tt.Results().At(i).Type()); ok {
					return p, true
				}
			}
		case *types.Struct:
			for i := 0; i < tt.NumFields(); i++ {
				if p, ok := walk(tt.Field(i).Type()); ok {
					return p, true
				}
			}
		case *types.Interface:
			for i := 0; i < tt.NumEmbeddeds(); i++ {
				if p, ok := walk(tt.EmbeddedType(i)); ok {
					return p, true
				}
			}
		case *types.TypeParam:
			if c := tt.Constraint(); c != nil {
				return walk(c)
			}
		}
		return "", false
	}
	return walk(t)
}
