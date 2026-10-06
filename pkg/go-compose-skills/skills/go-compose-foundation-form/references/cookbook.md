# foundation-form cookbook

Copy-paste call shapes for every component and the common patterns. All snippets
assume:

```go
import (
    fform "github.com/zodimo/go-compose/compose/foundation/form"
    "github.com/zodimo/go-compose/compose/foundation/form/components"
    "github.com/zodimo/go-compose/compose/foundation/form/rules"
    "github.com/zodimo/go-compose/compose/foundation/layout/column"
    "github.com/zodimo/go-compose/compose/foundation/layout/row"
    "github.com/zodimo/go-compose/compose/material3/button"
    "github.com/zodimo/go-compose/compose/material3/text"
    "github.com/zodimo/go-compose/pkg/api"
)
```

## Text field

```go
fform.ControlFieldOf[string](fs, "name", func(b *fform.FormFieldBinding[string]) api.Composable {
    return components.TextFieldComponent(
        b,
        components.TextFieldWithLabel("Name"),
        components.TextFieldWithHintText("Minimum 3 characters"),
    )
})
```

Options: `TextFieldWithLabel`, `TextFieldWithHintText`, `TextFieldWithInline(bool)`,
`TextFieldWithModifier`.

## Number field (int, cannot be empty)

```go
fform.ControlFieldOf[int](fs, "quantity", func(b *fform.FormFieldBinding[int]) api.Composable {
    return components.NumberFieldComponent(
        b,
        components.NumberFieldWithLabel("Quantity"),
        components.NumberFieldWithHintText("Whole number"),
        components.NumberFieldWithTextStateKey("quantity-text"), // optional, stabilizes the raw editor text
    )
})
```

Options: `NumberFieldWithLabel`, `NumberFieldWithHintText`,
`NumberFieldWithTextStateKey`, `NumberFieldWithModifier`. An int field always
holds a value; clearing falls back to the current value.

## Optional number field (clearable)

Control type is `fform.Optional[int]`:

```go
fform.ControlFieldOf[fform.Optional[int]](fs, "age", func(b *fform.FormFieldBinding[fform.Optional[int]]) api.Composable {
    return components.OptionalNumberFieldComponent(
        b,
        components.NumberFieldWithLabel("Age"),
        components.NumberFieldWithHintText("Leave blank if unknown"),
    )
})
```

Build the matching control with `None[int]()` and validate with
`fform.RequiredOptional[int]()` when it must be filled.

## Checkbox

```go
fform.ControlFieldOf[bool](fs, "sameAsShipping", func(b *fform.FormFieldBinding[bool]) api.Composable {
    return components.CheckboxComponent(
        b,
        components.CheckboxWithLabel("Same as shipping"),
        components.CheckboxWithErrorIndent(40), // align error text under the label
    )
})
```

Options: `CheckboxWithLabel`, `CheckboxWithHintText`, `CheckboxWithErrorIndent(int)`,
`CheckboxWithModifier`.

## Switch

```go
fform.ControlFieldOf[bool](fs, "notifications", func(b *fform.FormFieldBinding[bool]) api.Composable {
    return components.SwitchComponent(
        b,
        components.SwitchWithLabel("Enable notifications"),
    )
})
```

Options: `SwitchWithLabel`, `SwitchWithHintText`, `SwitchWithErrorIndent(int)`,
`SwitchWithModifier`.

## Select

`SelectComponent` takes the option slice as a **required positional argument**
before any component options. `SelectOption` is a struct literal:

```go
fform.ControlFieldOf[string](fs, "tier", func(b *fform.FormFieldBinding[string]) api.Composable {
    return components.SelectComponent(
        b,
        []components.SelectOption{
            {Value: "free", Label: "Free"},
            {Value: "pro", Label: "Pro"},
            {Value: "enterprise", Label: "Enterprise"},
        },
        components.SelectWithLabel("Tier"),
        components.SelectWithPlaceholder("Choose a tier"),
    )
})
```

Options: `SelectWithLabel`, `SelectWithPlaceholder`, `SelectWithHintText`,
`SelectWithModifier`. The bound control is a `string` holding the chosen
`Value`.

## Array of rows with add / remove

```go
func phoneSection(fs fform.FormScope) {
    phones, ok := fs.Array("phones")
    if !ok {
        return
    }

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
}
```

`Add`/`Insert`/`Remove` mutate the `Array` node directly; the tree notifies and
the form recomposes. Row paths are positional (`phones[index].number`), so
re-derive them from `index` each composition.

## Inline field beside a button (inside one row)

Fields must be their own list item; use `ControlViewOf` (not `ControlFieldOf`)
to embed one inside other content:

```go
fs.Field("search", func(fform.FormNode) api.Composable {
    return row.Row(c.Sequence(
        fform.ControlViewOf[string](fs, "query", func(b *fform.FormFieldBinding[string]) api.Composable {
            return components.TextFieldComponent(b,
                components.TextFieldWithLabel("Search"),
                components.TextFieldWithInline(true))
        }),
        button.Filled(func() { /* ... */ }, "Go"),
    ))(c)
})
```

## Cross-field wiring at build time

`buildForm` receives the composer, so wire sibling dependencies there:

```go
sameAsShipping := fform.NewControl(fform.NewPlainValueStore(false), false)
billing := fform.NewGroup(map[string]fform.FormNode{
    "street": fform.NewControl(fform.NewPlainValueStore(""), "", rules.StringNotEmpty()),
    "city":   fform.NewControl(fform.NewPlainValueStore(""), "", rules.StringNotEmpty()),
})
sameAsShipping.OnValueChange(func(on bool) { billing.SetDisabled(on) })
```

Cross-field *validation* is different — it reads the sibling's current value at
validation time:

```go
password := fform.NewControl(fform.NewPlainValueStore(""), "", rules.StringNotEmpty())
confirm := fform.NewControl(fform.NewPlainValueStore(""), "",
    rules.EqualTo(func() string { return password.Value() }))
```

## Group-level validator (at-least-one-of)

```go
contact := fform.NewGroup(
    map[string]fform.FormNode{
        "email": fform.NewControl(fform.NewPlainValueStore(""), "", rules.Email()),
        "phone": fform.NewControl(fform.NewPlainValueStore(""), "",
            rules.Numeric(), rules.MinLengthIfNotEmpty(10)),
    },
    fform.WithGroupValidator(rules.GroupAtLeastOneSet[string]("email", "phone")),
)
```

The failure is reported under the group's own path (`contact`), visible in
`formState.Errors()["contact"]`.

## Submit / Reset buttons

```go
func formActions(formState *fform.FormState) api.Composable {
    return func(c api.Composer) api.Composer {
        return row.Row(c.Sequence(
            button.Filled(func() {
                formState.Submit(func(value map[string]any) { /* use value */ })
            }, "Submit"),
            button.Outlined(func() { formState.Reset() }, "Reset"),
        ))(c)
    }
}
```

## Reading coded failures

```go
formState.Validate()
for _, ve := range formState.CodedErrors() {
    switch ve.Code {
    case rules.CodeEmail:
        // ve.Path, ve.Message
    }
}
```

Codes for the base package: `fform.CodeRequired`, `fform.CodeMinLength`. The
rules catalog defines `rules.Code*` constants (e.g. `rules.CodeEmail`,
`rules.CodeLength`). Wrap a hand-written validator with
`fform.WithCode("my_code", myValidator)`.

## Standalone field (not part of a form)

```go
itemState := state.MustRemember(c, "age", func() int { return 0 })
binding := fform.RememberFormFieldBinding(c, itemState)
return components.NumberFieldComponent(binding, components.NumberFieldWithLabel("Age"))(c)
```

Such fields are invisible to `FormState` (no form-level validate/reset/value).

## Testing bound paths

A cheap test that every path the UI binds actually resolves:

```go
func TestBuildFormPathsResolve(t *testing.T) {
    c := compose.NewComposer()
    root := buildForm(c)
    resolver := fform.NodeResolver(root)

    for _, path := range []string{"identity.name", "billing.street", "phones[0].number"} {
        if _, ok := fform.ControlOf[string](resolver, path); !ok {
            t.Errorf("expected %q to resolve to Control[string]", path)
        }
    }
}
```

Then drive the tree directly to test lifecycle behavior: resolve a control with
`fform.ControlOf[T](fform.NodeResolver(formState.Root()), path)`, call `Set`,
and assert on `formState.Submit` / `formState.Errors()`.
