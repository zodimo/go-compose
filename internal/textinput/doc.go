// Package textinput provides internal wrapper types for gioui.org text input
// widgets. These wrappers hide gioui types from the public API surface while
// allowing the input controller packages to delegate to the underlying widgets.
//
// The types in this package are defined types (not aliases) so the api-check
// analyzer sees them as internal types, not gioui types.
package textinput
