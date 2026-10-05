# foundation/form (fform)

A reactive, type-safe form engine for go-compose, inspired by Angular Reactive
Forms and Flutter's form architecture. The engine is pure Go and unit-testable;
a thin binding layer connects it to the framework's recomposition loop.

## Concepts

The form is a **tree of nodes**:

| Node | Role |
|------|------|
| `Control[T]` | A leaf holding one typed value, its validators, and lifecycle state. |
| `Group` | A container with named children (`map[string]FormNode`). |
| `Array` | A container with ordered children, supporting `Add`/`Insert`/`Remove`. |

Every node implements `FormNode`: `Status()`, `Validate()`, `Errors()`,
`IsTouched()`/`IsDirty()`/`IsPristine()`/`IsEnabled()`, `MarkAsTouched()`,
`Reset()`, `SetDisabled()`, `Parent()`/`Children()`/`Path()`/`Get()`, and
`RawValue()`.

- **Status** is computed recursively: a group is `INVALID` if any child is
  invalid, `DISABLED` if all children are disabled, `PENDING` while async work is
  outstanding, otherwise `VALID`. A disabled node reports `DISABLED`.
- **Disabled** distinguishes *own* from *inherited*: disabling a parent does not
  destroy children's own settings, so re-enabling restores them.
- **Validation** runs on a value snapshot (never under the node lock, so
  reentrant validators cannot deadlock) and aggregates failures.
- **Change notification** propagates: a child mutation notifies its ancestors'
  listeners, so a single root subscription covers the whole tree.

## Quick start

```go
func FormScreen() api.Composable {
    return func(c api.Composer) api.Composer {
        formState := fform.RememberFormState(c, buildForm)

        return fform.Form(formState, func(fs fform.FormScope) {
            // A full-width text field, resolved by path and typed.
            fform.ControlFieldOf[string](fs, "identity.name", func(b *fform.FormFieldBinding[string]) api.Composable {
                return components.TextFieldComponent(
                    b,
                    components.TextFieldWithLabel("Name"),
                    components.TextFieldWithHintText("Minimum 3 characters"),
                )
            })

            // A bool checkbox.
            fform.ControlFieldOf[bool](fs, "sameAsShipping", func(b *fform.FormFieldBinding[bool]) api.Composable {
                return components.CheckboxComponent(b, components.CheckboxWithLabel("Same as shipping"))
            })
        })(c)
    }
}

func buildForm(c api.Composer) *fform.Group {
    return fform.NewGroup(map[string]fform.FormNode{
        "identity": fform.NewGroup(map[string]fform.FormNode{
            "name": fform.NewControl(fform.NewPlainValueStore(""), "", fform.Required(""), fform.MinLength(3)),
        }),
        "sameAsShipping": fform.NewControl(fform.NewPlainValueStore(false), false),
    })
}
```

Values persist because the tree is remembered by `RememberFormState`; any tree
mutation bumps a version that drives recomposition, so a plain value store per
control is enough for tree-integrated fields.

## Form lifecycle

`*FormState` exposes the whole-form surface:

```go
formState.Validate()          // validates every node, returns bool
formState.Status()            // VALID | INVALID | PENDING | DISABLED
formState.Errors()            // map[dottedPath]message
formState.Value()             // map[string]any (disabled nodes omitted)
formState.IsTouched() / IsDirty()
formState.MarkAllTouched()    // surface all errors (attempted submit)
formState.Reset()
formState.SetDisabled(true)
formState.Submit(onValid)     // validate -> touch-all -> onValid(value) if valid
```

## Scope helpers

Inside `Form`'s content closure you receive a `FormScope`. Keys resolve by path
relative to the scope (`a.b.c`, `items[0].name`); an unresolvable key renders
nothing rather than panicking.

| Helper | Purpose |
|--------|---------|
| `fs.Field(key, content)` | Resolve a node and render content with it. |
| `fs.Form(key, header, children)` | Header row + nested scope resolving relative paths. |
| `fs.GroupScope(key, children)` | Nested scope with no header. |
| `fs.FormArray(key, row)` | One row per array child, with index + row-relative scope. |
| `fs.ArrayLength(key)` | Number of children of an array (for add/remove affordances). |
| `fs.Array(key)`, `fs.Group(key)`, `fs.Resolve(key)` | Typed/raw escape hatches. |
| `fform.ControlOf[T](resolver, key)` | Type-safe `*Control[T]` resolution. |
| `fform.ControlFieldOf[T](scope, key, content)` | Field as its own list item. |
| `fform.ControlViewOf[T](resolver, key, content)` | Field embedded inline (e.g. in an array row). |

`ControlOf`/`ControlViewOf` accept a `Resolver` — a `FormScope`, or any node
wrapped with `fform.NodeResolver(node)`.

## Components

`compose/foundation/form/components` provides bound Material 3 widgets. Each
takes a `*FormFieldBinding[T]`, reads value/errors/touched state, writes through
`SetValue` (which marks the control touched), and honors the control's enabled
state. Errors are touched-gated per the Material 3 convention.

- `TextFieldComponent(*FormFieldBinding[string], ...)`
- `NumberFieldComponent(*FormFieldBinding[int], ...)` — keeps raw text so partial
  input is typable; writes only parsed ints and reports unparsable text. An int
  control cannot be empty, so clearing it falls back to the current value.
- `OptionalNumberFieldComponent(*FormFieldBinding[Optional[int]], ...)` — the
  empty-capable number field; clearing writes `None`, see below.
- `CheckboxComponent` / `SwitchComponent(*FormFieldBinding[bool], ...)`
- `SelectComponent(*FormFieldBinding[string], []SelectOption, ...)`

## Empty vs. zero: `Optional[T]`

A `Control[int]` cannot represent "empty" — `0` is a real value, so an int field
can't be left blank and `Required(0)` can't tell "user typed 0" from "no input".
Use `Control[Optional[T]]` (an alias for `github.com/zodimo/go-maybe`) when the
field may legitimately be empty:

```go
"age": fform.NewControl[fform.Optional[int]](
    fform.NewPlainValueStore(fform.None[int]()),
    fform.None[int](),
    fform.RequiredOptional[int](), // fails on None, accepts Some(0)
),
```

```go
fform.ControlFieldOf[fform.Optional[int]](fs, "age", func(b *fform.FormFieldBinding[fform.Optional[int]]) api.Composable {
    return components.OptionalNumberFieldComponent(b, components.NumberFieldWithLabel("Age"))
})
```

Helpers: `Some(v)` / `None[T]()` / `OptionalOf(v, set)`, `RequiredOptional[T]()`,
`NewOptionalControl[T](store)`, and `OptionalValue[T](control)` for unwrapping.
`Required(None[T]())` also works (Optional is comparable), but `RequiredOptional`
reads more clearly. Because `Optional` is `maybe.Maybe`, values export through the
tree as `Optional[T]` and interoperate with json `omitzero`, `Match`, `Map`, etc.

## Dynamic arrays

```go
fs.FormArray("phones", func(index int, rowScope fform.FormScope) api.Composable {
    field := fform.ControlViewOf[string](rowScope, "number", func(b *fform.FormFieldBinding[string]) api.Composable {
        return components.TextFieldComponent(b, components.TextFieldWithLabel(fmt.Sprintf("Phone %d", index+1)))
    })
    return row.Row(c.Sequence(field, removeButton(func() { phones.Remove(index) })))
})
```

A row scope's path is positional (`phones[index]`), so array indices shift on
`Remove`/`Insert`. Resolve row paths from the index each composition; do not cache
path strings across mutations.

## Layout notes

`Form` renders a `LazyColumn`, so each `Field`/`Form`/`FormArray` row becomes a
list item keyed by the base path plus key. This means a field cannot be split
across list items, and inline composition (for example a field plus a remove
button in one row) uses `ControlViewOf` inside a `row.Row`.

## Standalone fields (no form tree)

A field does not have to live in a `Form`. `RememberFormFieldBinding` builds a
tree-free binding whose control is driven directly by a remembered
`state.MutableValueTyped[T]`:

```go
itemState := state.MustRemember(c, "age", func() int { return 0 })
binding := fform.RememberFormFieldBinding(c, itemState)
return components.NumberFieldComponent(binding, components.NumberFieldWithLabel("Age"))(c)
```

Validation, touched-gated errors, and recomposition work as usual. The trade-off:
standalone fields are invisible to `FormState` — form-level `Validate`, `Reset`,
`Value`, and `CodedErrors` only see controls in the tree. Prefer
`ControlFieldOf`/`ControlViewOf` for anything that belongs to a form.

## Validation rules and error codes

The base package ships `Required` and `MinLength`. A fuller rule catalog lives in
`compose/foundation/form/rules` — every rule is a plain `fform.ValidatorFunc[T]`,
so it drops straight into `NewControl`:

```go
fform.NewControl(fform.NewPlainValueStore(""), "",
    rules.StringNotEmpty(),
    rules.MinLength(3),
    rules.AlphaUnicode(),
)

fform.NewControl(fform.NewPlainValueStore(""), "",
    rules.StringNotEmpty(),
    rules.Email(),
)

// Cross-field: confirmation must equal the password control's current value.
fform.NewControl(fform.NewPlainValueStore(""), "",
    rules.EqualTo(func() string { return password.Value() }),
)
```

Catalog: required/not-empty, length (runes and items), string shape (regex,
contains/excludes, starts/ends-with, alpha/numeric/alphanumeric), formats (email,
url, uuid, ip/ipv4/ipv6, dns label/subdomain, base64, hex, semver), comparables
(one-of/none-of, eq/neq, gt/gte/lt/lte, in-range), collections (unique, unique-by,
slice contains, map key/value), and cross-field (equal-to, greater/less than
field, required-when, forbidden-when, at-least-one-set, mutually-exclusive).

Length rules come in two flavors: `MinLength` fails on empty input (0 < min);
`MinLengthIfNotEmpty` and `Length` pass on empty and enforce a size only once the
user has typed something — the right choice for an optional field with a minimum
size when present.

### Group-level validation

Some cross-field rules are about the *group*, not any single child: "provide at
least one of email/phone", "these are mutually exclusive", "end must be after
start". Attach them to the group with `WithGroupValidator`; the failure is
reported against the group's own path.

```go
fform.NewGroup(
    map[string]fform.FormNode{
        "email": email,
        "phone": phone,
    },
    fform.WithGroupValidator(rules.GroupAtLeastOneSet[string]("email", "phone")),
)
```

Group rules read named children (dotted paths allowed), ignore disabled and
unresolvable children, and compose with the same coded-error machinery, so
`formState.Errors()["contact"]` and `formState.CodedErrors()` both see the
failure. Catalog: `GroupAtLeastOneSet`, `GroupMutuallyExclusive`,
`GroupRequiredTogether`, `GroupOrdered`, and `GroupCustom`. Any
`fform.GroupValidatorFunc` (`func(*fform.Group) error`) works, so hand-written
group validators are first-class too.

Group validators run on read (the engine never caches status), so they must be
side-effect free and cheap, and they never run for an effectively disabled group.
For a rule that concerns one field but compares it to a sibling, prefer the
per-control cross-field rules (`EqualTo`, `GreaterThanField`); use a group
validator when the subject is the group.

### Machine-readable failures

Every rule raises a `*fform.ValidationError` carrying a stable `ErrorCode`, so you
can branch on *why* something failed without parsing messages:

```go
formState.Validate()
for _, ve := range formState.CodedErrors() { // []*ValidationError{Path, Code, Message}
    switch ve.Code {
    case rules.CodeEmail:
        // localize / highlight
    }
}

codes := formState.Codes()            // distinct codes across the tree, sorted
path, ok := formState.FirstInvalidPath() // document-order first invalid field
```

`CodeOf(err)` / `CodesOf(err)` work on any error, unwrapping through `errors.Join`
and `%w` chains. To attach a code to a hand-written validator, wrap it:

```go
fform.WithCode("phone_format", func(s string) error { ... })
```

Existing behavior is unchanged: `Errors() map[string]string` still works, and
`Required`/`MinLength` in the base package now emit codes (`CodeRequired`,
`CodeMinLength`).

### Provenance

The `rules` catalog is inspired by [govy](https://github.com/nobl9/govy) and other
Go validation libraries, but every rule here is reimplemented from scratch in this
repository. No third-party source is copied, so no external license obligations
attach to the `rules` package.
