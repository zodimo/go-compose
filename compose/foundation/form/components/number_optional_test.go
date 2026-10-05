package components

import (
	"strconv"
	"testing"

	fform "github.com/zodimo/go-compose/compose/foundation/form"
)

func TestOptionalText_MapsNoneToEmpty(t *testing.T) {
	if got := optionalText(fform.None[int]()); got != "" {
		t.Fatalf("None should render as empty string, got %q", got)
	}
	if got := optionalText(fform.Some(0)); got != "0" {
		t.Fatalf("Some(0) should render as \"0\", got %q", got)
	}
	if got := optionalText(fform.Some(-42)); got != "-42" {
		t.Fatalf("Some(-42) should render as \"-42\", got %q", got)
	}
}

func TestOptionalKey_DistinguishesNoneFromSomeZero(t *testing.T) {
	if optionalKey(fform.None[int]()) == optionalKey(fform.Some(0)) {
		t.Fatal("None and Some(0) must have distinct keys so the buffer re-syncs")
	}
	if optionalKey(fform.Some(1)) == optionalKey(fform.Some(2)) {
		t.Fatal("distinct values must have distinct keys")
	}
	if optionalKey(fform.Some(5)) != optionalKey(fform.Some(5)) {
		t.Fatal("equal values must have equal keys")
	}
}

// TestOptionalNumberRoundTrip exercises the parse/empty mapping the component
// applies on every keystroke: empty -> None, parseable -> Some, unparsable ->
// unchanged control (the component leaves the last value and shows an error).
func TestOptionalNumberRoundTrip(t *testing.T) {
	write := func(control *fform.Control[fform.Optional[int]], s string) {
		switch {
		case s == "":
			control.Set(fform.None[int]())
		default:
			if n, err := strconv.Atoi(s); err == nil {
				control.Set(fform.Some(n))
			}
		}
	}

	control := fform.NewOptionalControl(fform.NewPlainValueStore(fform.None[int]()))

	write(control, "12")
	if got, ok := fform.OptionalValue(control); !ok || got != 12 {
		t.Fatalf("expected (12, true), got (%d, %t)", got, ok)
	}

	write(control, "")
	if _, ok := fform.OptionalValue(control); ok {
		t.Fatal("clearing the field should set None")
	}

	write(control, "0")
	if got, ok := fform.OptionalValue(control); !ok || got != 0 {
		t.Fatalf("expected (0, true) for \"0\", got (%d, %t)", got, ok)
	}

	// Unparsable input leaves the last good value in place.
	write(control, "abc")
	if got, ok := fform.OptionalValue(control); !ok || got != 0 {
		t.Fatalf("unparsable input must keep the last value (0), got (%d, %t)", got, ok)
	}
}
