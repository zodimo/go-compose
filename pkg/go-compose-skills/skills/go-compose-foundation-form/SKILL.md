---
name: go-compose-foundation-form
description: Use the go-compose reactive form engine (compose/foundation/form, imported as fform) — building a form tree of Control/Group/Array nodes, rendering it with Form + FormScope helpers, binding Material 3 field components, choosing validators from the rules catalog, and driving the validate/submit/reset lifecycle. Use when adding, editing, or debugging a form, form field, validation rule, or form-level action (submit/reset/status) in this repository.
---

# go-compose foundation form (`fform`)

**Requires go-compose v0.1.125 or newer.** The API described here (the
`fform.FormScope` helpers, `FormFieldBinding`, the `rules` catalog, and the
coded-error surface) assumes that minimum; on older versions some symbols are
missing or behave differently. Check the module version before relying on any
signature in this skill:

```bash
go list -m github.com/zodimo/go-compose    # must be >= v0.1.125
```

A reactive, type-safe form engine. The engine is **pure Go and unit-testable**;
a thin binding layer connects it to the framework's recomposition loop.

Import path:

```go
fform "github.com/zodimo/go-compose/compose/foundation/form"
"github.com/zodimo/go-compose/compose/foundation/form/components"
"github.com/zodimo/go-compose/compose/foundation/form/rules"
```

Canonical, working example: `cmd/demo/form/ui.go` (+ `ui_test.go`). Read it
before writing a new form; it exercises nested groups, a dynamic array,
cross-field disable, group validation, and the submit/reset lifecycle.

For copy-paste snippets of every component and pattern, see
`references/cookbook.md`. Read it when you need the exact call shape; this file
explains the model and the rules.

## The mental model

The form is a **tree of nodes**:

| Node | Role |
|------|------|
| `Control[T]` | Leaf holding one typed value, its validators, and lifecycle state. |
| `Group` | Container with named children (`map[string]FormNode`). |
| `Array` | Container with ordered children; `Add`/`Insert`/`Remove`. |

Three rules govern everything:

1. **The tree is the single source of truth.** You build it once in a
   `buildForm` function and remember it. UI binds *by path* into that tree.
2. **Values persist because the tree is remembered.** `Form` renders a
   `LazyColumn`, subscribes once to the root's change stream, and bumps a
   remembered version on any mutation. A plain `NewPlainValueStore` per control
   is enough — you do **not** need a `MutableValue` per field.
3. **Validation runs on read, never cached.** Status is recomputed whenever it
   is asked for. Validators (including group validators) must be side-effect
   free and cheap.

## Rendering: the `FormScope` helpers

`Form(formState, content, ...options)` renders a `LazyColumn` and hands
`content` a `FormScope`. Call `fform.RememberFormState(c, buildForm)` to build
and remember the tree, then pass the state to `Form`.

Keys resolve by path relative to the scope (`a.b.c`, `items[0].name`). **An
unresolvable key renders nothing rather than panicking** — the tree decides what
exists.

| Helper | Purpose |
|--------|---------|
| `fform.ControlFieldOf[T](fs, key, content)` | Field as its own full-width list item. **The default.** |
| `fform.ControlViewOf[T](resolver, key, content)` | Field embedded inline (e.g. in an array row with a remove button). |
| `fs.Field(key, content)` | Resolve any node, render content with the raw `FormNode`. |
| `fs.Form(key, header, children)` | Header row + nested scope over the subtree named by `key`. |
| `fs.GroupScope(key, children)` | Nested scope with no header. |
| `fs.FormArray(key, row)` | One list item per array child; row gets `(index, rowScope)`. |
| `fs.ArrayLength(key)` | Child count of an array (for add/remove affordances). |
| `fs.Array(key)` / `fs.Group(key)` / `fs.Resolve(key)` | Typed/raw escape hatches. |
| `fform.ControlOf[T](resolver, key)` | Type-safe `*Control[T]` resolution. |

Two things about these signatures that are easy to get wrong:

- `ControlOf`/`ControlViewOf` take a `Resolver` — a `FormScope`, or any node
  wrapped with `fform.NodeResolver(node)`.
- `ControlOf`/`ControlFieldOf` are **package-level generic functions, not scope
  methods** (Go forbids type parameters on interface methods).
- The `children`/`row` callbacks receive a **narrowed scope**, not the outer
  `fs`. Inside `fs.Form("identity", ..., func(fs fform.FormScope) { ... })` you
  shadow `fs` with a scope rooted at `identity`, so `"password"` resolves to
  `identity.password`. Bind paths relative to the scope you were handed.

**Layout constraint:** `Form` renders a `LazyColumn`, so each
`Field`/`Form`/`FormArray` call appends a list item keyed by base path + key. A
field cannot be split across list items. To place a field *beside* something
(button, label) inside one row, use `ControlViewOf` inside a `row.Row`.

## The binding: `FormFieldBinding[T]`

Every component takes a `*FormFieldBinding[T]`, which wraps a `*Control[T]`:

- Reads: `Value()`, `Errors()`, `HasErrors()`, `IsTouched()`, `IsDirty()`,
  `IsEnabled()`, `Status()`, `ErrorMessage()`.
- Writes: `SetValue(v)` — sets the value **and marks the control touched**.
- `Control() *Control[T]` for direct access.

Errors are **touched-gated** per Material 3: a component only shows an error
message once the control is touched, so invalid-but-unvisited fields render
cleanly. This rule lives in `components/field.go` and applies to every widget.

## Components (`.../form/components`)

Each reads value/errors/touched and honors the control's enabled state:

- `TextFieldComponent(*FormFieldBinding[string], ...)`
- `NumberFieldComponent(*FormFieldBinding[int], ...)` — keeps raw text so partial
  input like `-` or `1.` is typable; writes only parsed ints and reports
  unparsable text. **An int control cannot be empty**, so clearing falls back to
  the current value.
- `OptionalNumberFieldComponent(*FormFieldBinding[Optional[int]], ...)` — the
  empty-capable number field.
- `CheckboxComponent` / `SwitchComponent(*FormFieldBinding[bool], ...)`
- `SelectComponent(*FormFieldBinding[string], []SelectOption, ...)` — note the
  option slice is a required positional argument, before component options.

Option constructors follow `WidgetWithX`, e.g. `TextFieldWithLabel`,
`TextFieldWithHintText`, `TextFieldWithInline`, `NumberFieldWithLabel`,
`NumberFieldWithHintText`, `NumberFieldWithTextStateKey`, `CheckboxWithLabel`,
`CheckboxWithErrorIndent`, `SelectWithLabel`, `SelectWithPlaceholder`,
`SelectWithHintText`. See the cookbook for a full call per widget.

## Minimal complete form

```go
func FormScreen() api.Composable {
    return func(c api.Composer) api.Composer {
        formState := fform.RememberFormState(c, buildForm)

        return column.Column(
            c.Sequence(
                fform.Form(formState, formContent),
                submitButton(formState),
            ),
        )(c)
    }
}

// Build the tree ONCE. Receives the composer so cross-field wiring can be set
// up at build time.
func buildForm(c api.Composer) *fform.Group {
    password := fform.NewControl(fform.NewPlainValueStore(""), "",
        rules.StringNotEmpty(), rules.MinLength(8))
    sameAsShipping := fform.NewControl(fform.NewPlainValueStore(false), false)

    billing := fform.NewGroup(map[string]fform.FormNode{
        "street": fform.NewControl(fform.NewPlainValueStore(""), "", rules.StringNotEmpty()),
        "city":   fform.NewControl(fform.NewPlainValueStore(""), "", rules.StringNotEmpty()),
    })
    // Cross-field: checking the box disables the billing subtree; unchecking
    // restores each child's own setting.
    sameAsShipping.OnValueChange(func(on bool) { billing.SetDisabled(on) })

    return fform.NewGroup(map[string]fform.FormNode{
        "identity": fform.NewGroup(map[string]fform.FormNode{
            "password": password,
            "confirm": fform.NewControl(fform.NewPlainValueStore(""), "",
                rules.EqualTo(func() string { return password.Value() })),
        }),
        "sameAsShipping": sameAsShipping,
        "billing":        billing,
    })
}

// Bind tree nodes to rendered components. Nested Form/GroupScope scopes make
// keys resolve relative to a subtree.
func formContent(fs fform.FormScope) {
    fs.Form("identity", text.TitleLarge("Identity"), func(fs fform.FormScope) {
        fform.ControlFieldOf[string](fs, "password", func(b *fform.FormFieldBinding[string]) api.Composable {
            return components.TextFieldComponent(b,
                components.TextFieldWithLabel("Password"),
                components.TextFieldWithHintText("Minimum 8 characters"))
        })
        fform.ControlFieldOf[string](fs, "confirm", func(b *fform.FormFieldBinding[string]) api.Composable {
            return components.TextFieldComponent(b, components.TextFieldWithLabel("Confirm password"))
        })
    })
    fs.Form("billing", text.TitleLarge("Billing Address"), func(fs fform.FormScope) {
        fform.ControlFieldOf[string](fs, "street", func(b *fform.FormFieldBinding[string]) api.Composable {
            return components.TextFieldComponent(b, components.TextFieldWithLabel("Street"))
        })
        fform.ControlFieldOf[string](fs, "city", func(b *fform.FormFieldBinding[string]) api.Composable {
            return components.TextFieldComponent(b, components.TextFieldWithLabel("City"))
        })
    })
}

func submitButton(formState *fform.FormState) api.Composable {
    return func(c api.Composer) api.Composer {
        return button.Filled(func() {
            // Submit validates, marks all touched so errors render, and only
            // calls onValid when the form is valid.
            formState.Submit(func(value map[string]any) { /* use value */ })
        }, "Submit")(c)
    }
}
```

## Empty vs. zero: `Optional[T]`

A `Control[int]` cannot represent "empty" — `0` is a real value, so an int field
can't be blank and `Required(0)` can't tell "user typed 0" from "no input". Use
`Control[Optional[T]]` (`Optional` aliases `github.com/zodimo/go-maybe.Maybe`)
when the field may legitimately be empty:

```go
"age": fform.NewControl[fform.Optional[int]](
    fform.NewPlainValueStore(fform.None[int]()),
    fform.None[int](),
    fform.RequiredOptional[int](), // fails on None, accepts Some(0)
),
```

Helpers: `Some(v)` / `None[T]()` / `OptionalOf(v, set)`, `RequiredOptional[T]()`,
`NewOptionalControl[T](store)`, `OptionalValue[T](control)` for unwrapping.

## Dynamic arrays

```go
phones, _ := fs.Array("phones")

fs.Field("phones", func(fform.FormNode) api.Composable {
    return row.Row(c.Sequence(
        text.TitleLarge("Phone Numbers"),
        button.Text(func() { phones.Add(phoneRow("")) }, "Add"),
    ))(c)
})

fs.FormArray("phones", func(index int, rowScope fform.FormScope) api.Composable {
    field := fform.ControlViewOf[string](rowScope, "number",
        func(b *fform.FormFieldBinding[string]) api.Composable {
            return components.TextFieldComponent(b,
                components.TextFieldWithLabel(fmt.Sprintf("Phone %d", index+1)))
        })
    return row.Row(c.Sequence(
        field,
        button.Text(func() { phones.Remove(index) }, "Remove"),
    ))(c)
})
```

**A row scope's path is positional** (`phones[index]`), so array indices shift
on `Remove`/`Insert`. Resolve row paths from the index on each composition —
never cache path strings across mutations. Row stores are plain holders owned by
the tree (`NewPlainValueStore`), not remembered composition state.

## Validators

Two validator kinds:

- `fform.ValidatorFunc[T]` — `func(T) error`, attached to a `Control[T]` at
  `NewControl`.
- `fform.GroupValidatorFunc` — `func(*fform.Group) error`, attached to a `Group`
  with `fform.WithGroupValidator(...)`.

The base package ships `fform.Required(zero)`, `fform.MinLength(n)`, and
`fform.RequiredOptional[T]()`. Everything else lives in the `rules` catalog —
every rule is a plain `fform.ValidatorFunc[T]`, so it drops straight into
`NewControl`:

```go
fform.NewControl(fform.NewPlainValueStore(""), "",
    rules.StringNotEmpty(), rules.MinLength(3), rules.AlphaUnicode())

// Cross-field: confirmation must equal the password's current value.
rules.EqualTo(func() string { return password.Value() })
```

**Catalog** (`compose/foundation/form/rules`):
required/not-empty; length (runes and items); string shape (regex,
contains/excludes, starts/ends-with, alpha/numeric/alphanumeric, unicode
variants); formats (email, url, uuid, ip/ipv4/ipv6, dns label/subdomain, base64,
hex, semver); comparables (one-of/none-of, eq/neq/gt/gte/lt/lte, `InRange`);
collections (`Unique`, `UniqueBy`, `SliceContainsAll`, map key/value); cross-field
(`EqualTo`, `NotEqualTo`, `GreaterThanField`, `LessThanField`, `RequiredWhen`,
`ForbiddenWhen`, `AtLeastOneSet`, `MutuallyExclusive`); `Custom`.

### Empty-input gotcha in the length rules

`MinLength` fails on empty input because `0 < min` — so a `MinLength` control is
also implicitly required. The rules that **let empty input pass** are
`MinLengthIfNotEmpty`, `Length`, and `MaxLength`:

- Use `MinLengthIfNotEmpty(n)` for an optional field that must meet a minimum
  size *only once something is typed* (it guards with `value == ""`).
- `Length(min, max)` likewise **accepts empty**, so it enforces a size band only
  for non-empty input. To make a field both required and bounded, combine it
  with `StringNotEmpty()` (or pair `MinLength` with an explicit optional story).
- `MaxLength(n)` never rejects empty on its own (an empty string trivially
  satisfies an upper bound); add `StringNotEmpty()` if presence matters.

### Group-level validation

Some cross-field rules concern the *group*, not any single child: "provide at
least one of email/phone", "these are mutually exclusive", "end must be after
start". Attach them to the group; the failure is reported against the group's
own path.

```go
fform.NewGroup(
    map[string]fform.FormNode{"email": email, "phone": phone},
    fform.WithGroupValidator(rules.GroupAtLeastOneSet[string]("email", "phone")),
)
```

Group rules read named children (dotted paths allowed), ignore disabled and
unresolvable children, and run on read — so they must be side-effect free and
cheap, and never run for an effectively disabled group. Rule catalog:
`GroupAtLeastOneSet`, `GroupMutuallyExclusive`, `GroupRequiredTogether`,
`GroupOrdered`, `GroupCustom`. Any hand-written `fform.GroupValidatorFunc` works
too. For a rule about one field compared to a sibling, prefer the per-control
cross-field rules (`EqualTo`, `GreaterThanField`); use a group validator when the
subject is the group.

## Lifecycle: `*FormState`

```go
formState.Validate()          // validate every node, returns bool
formState.Status()            // StatusValid | StatusInvalid | StatusPending | StatusDisabled
formState.Errors()            // map[dottedPath]message
formState.Value()             // map[string]any (disabled nodes omitted)
formState.IsTouched() / IsDirty() / IsPristine() / IsEnabled()
formState.MarkAllTouched()    // surface all errors (attempted submit)
formState.MarkAllUntouched() / MarkAllPristine()
formState.Reset()             // initial values, cleared errors, pristine+untouched
formState.SetDisabled(true)
formState.Submit(onValid)     // validate -> touch-all -> onValid(value) only if valid; returns bool
formState.Root()              // *Group for advanced wiring
```

`Submit` is the conventional button handler and encodes the
validate → touch-all → proceed-only-if-valid flow. Call it, don't reorder those
steps by hand.

**Status semantics:** a group is `INVALID` if any child is invalid, `DISABLED`
if all children are disabled, `PENDING` while async work is outstanding,
otherwise `VALID`. Disabling a parent does **not** destroy children's own
disabled settings, so re-enabling restores them.

## Machine-readable failures

Every rule raises a `*fform.ValidationError` carrying a stable `ErrorCode`, so
you can branch on *why* something failed without parsing messages:

```go
formState.Validate()
for _, ve := range formState.CodedErrors() { // []*ValidationError{Path, Code, Message}
    switch ve.Code {
    case rules.CodeEmail:
        // localize / highlight
    }
}
codes := formState.Codes()                 // distinct codes across the tree, sorted
path, ok := formState.FirstInvalidPath()   // lexically-first invalid field (deterministic)
```

Validators themselves leave `Path` empty; the binding layer stamps the control's
dotted path. `fform.CodeOf(err)` / `fform.CodesOf(err)` unwrap through
`errors.Join` and `%w` chains. To code a hand-written validator, wrap it:
`fform.WithCode("phone_format", func(s string) error { ... })`. On a binding, use
`binding.Codes()` / `binding.ValidationErrors()`.

## Standalone fields (no form tree)

A field can live outside a `Form`. `RememberFormFieldBinding` builds a tree-free
binding driven by a remembered `state.MutableValueTyped[T]`:

```go
itemState := state.MustRemember(c, "age", func() int { return 0 })
binding := fform.RememberFormFieldBinding(c, itemState)
return components.NumberFieldComponent(binding, components.NumberFieldWithLabel("Age"))(c)
```

Validation, touched-gated errors, and recomposition work as usual. **Trade-off:**
standalone fields are invisible to `FormState` — form-level `Validate`, `Reset`,
`Value`, and `CodedErrors` only see controls in the tree. Prefer
`ControlFieldOf`/`ControlViewOf` for anything that belongs to a form.

## Conventions and footguns

- **Build the tree in one `buildForm(c)` function** and remember it with
  `RememberFormState(c, buildForm)`. Never rebuild nodes during composition.
- **Keys that don't resolve render nothing** (no panic). If a field is missing
  from the UI, check the path and that the node exists in the tree.
- **The scope callback is narrowed.** Inside a nested `fs.Form`/`fs.GroupScope`,
  bind paths relative to the inner scope — don't repeat the parent prefix.
- **Resolve-and-assert escape hatch:** `fform.ControlOf[T](resolver, key)` returns
  `(*Control[T], bool)`; `NodeResolver(node)` adapts a bare node (e.g. the group
  passed to a group validator) to a `Resolver`.
- **Direct `Control.Set` is a no-op when disabled**, so a disabled subtree silently
  drops writes — that is how cross-field disable works.
- **`NumberFieldComponent` loses empty state.** Use `OptionalNumberFieldComponent`
  for a field that must be clearable.
- **Row paths are positional.** Re-derive `phones[i].number` from the index each
  composition; don't cache it.
- **Validators run on read** — keep them pure and fast; no I/O, no mutation.

## Verify your work

```bash
go build ./...
go test ./compose/foundation/form/...
go run ./cmd/demo/form/          # run the reference demo
```

`cmd/demo/form/ui_test.go` shows the intended test style: build the tree, wrap
the root with `fform.NodeResolver`, then assert on `fform.ControlOf[T]` and array
lengths. A test that every bound path resolves is cheap insurance against the
silent "renders nothing" miss.
