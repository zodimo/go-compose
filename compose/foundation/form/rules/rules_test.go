package rules

import (
	"errors"
	"regexp"
	"testing"

	fform "github.com/zodimo/go-compose/compose/foundation/form"
)

// expectFails asserts validator(v) fails and carries want code.
func expectFails(t *testing.T, name string, err error, want Code) {
	t.Helper()
	if err == nil {
		t.Fatalf("%s: expected failure, got nil", name)
	}
	got, ok := fform.CodeOf(err)
	if !ok {
		t.Fatalf("%s: expected coded error, got %v", name, err)
	}
	if got != want {
		t.Fatalf("%s: expected code %q, got %q", name, want, got)
	}
}

func expectPasses(t *testing.T, name string, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("%s: expected pass, got %v", name, err)
	}
}

func TestRequiredAndNotEmpty(t *testing.T) {
	expectFails(t, "required empty", Required("")(""), CodeRequired)
	expectPasses(t, "required set", Required("")("x"))
	// Required("") accepts whitespace; StringNotEmpty does not.
	expectPasses(t, "required whitespace", Required("")(" "))
	expectFails(t, "notempty whitespace", StringNotEmpty()("   "), CodeRequired)
	expectPasses(t, "notempty text", StringNotEmpty()("hi"))
}

func TestLengthRunes(t *testing.T) {
	// Rune counting: 2 runes but 6 bytes.
	expectPasses(t, "min runes", MinLength(2)("héllo"))
	expectFails(t, "min too short", MinLength(2)("h"), CodeMinLength)
	expectFails(t, "max too long", MaxLength(3)("abcd"), CodeMaxLength)
	expectPasses(t, "length ok", Length(2, 4)("abc"))
	expectFails(t, "length low", Length(2, 4)("a"), CodeLength)
	expectFails(t, "length high", Length(2, 4)("abcde"), CodeLength)
	// Empty passes for Length and MinLengthIfNotEmpty: optional fields.
	expectPasses(t, "length empty ok", Length(2, 4)(""))
	expectPasses(t, "min-if-not-empty empty", MinLengthIfNotEmpty(3)(""))
	expectFails(t, "min-if-not-empty short", MinLengthIfNotEmpty(3)("ab"), CodeMinLength)
	expectPasses(t, "min-if-not-empty ok", MinLengthIfNotEmpty(3)("abc"))
}

func TestSliceLength(t *testing.T) {
	expectFails(t, "min items", MinItems[string](2)([]string{"a"}), CodeMinLength)
	expectPasses(t, "min items ok", MinItems[string](1)([]string{"a"}))
	expectFails(t, "max items", MaxItems[string](1)([]string{"a", "b"}), CodeMaxLength)
}

func TestRuneClasses(t *testing.T) {
	expectPasses(t, "alpha", Alpha()("Abc"))
	expectFails(t, "alpha digit", Alpha()("Ab1"), CodeAlpha)
	expectFails(t, "alpha ascii", Alpha()("café"), CodeAlpha)
	expectPasses(t, "alpha unicode", AlphaUnicode()("café"))
	expectPasses(t, "numeric", Numeric()("123"))
	expectFails(t, "numeric letter", Numeric()("12a"), CodeNumeric)
	expectPasses(t, "alnum", Alphanumeric()("a1B2"))
	expectFails(t, "alnum symbol", Alphanumeric()("a-1"), CodeAlphanumeric)
	// Empty passes: pair with StringNotEmpty.
	expectPasses(t, "alpha empty", Alpha()(""))
}

func TestStringShape(t *testing.T) {
	re := regexp.MustCompile(`^abc$`)
	expectPasses(t, "match", MatchRegexp(re)("abc"))
	expectFails(t, "match miss", MatchRegexp(re)("abd"), CodeMatchRegexp)
	expectFails(t, "deny hit", DenyRegexp(re)("abc"), CodeDenyRegexp)
	expectPasses(t, "deny miss", DenyRegexp(re)("xyz"))

	expectPasses(t, "contains", Contains("a", "b")("ab"))
	expectFails(t, "contains miss", Contains("a", "z")("ab"), CodeContains)
	expectFails(t, "excludes hit", Excludes("z")("xyz"), CodeExcludes)
	expectPasses(t, "excludes miss", Excludes("q")("xyz"))
	expectPasses(t, "startswith", StartsWith("http://", "https://")("https://x"))
	expectFails(t, "startswith miss", StartsWith("http://")("ftp://x"), CodeStartsWith)
	expectPasses(t, "endswith", EndsWith(".go", ".md")("a.go"))
	expectFails(t, "endswith miss", EndsWith(".go")("a.txt"), CodeEndsWith)
}

func TestFormatRules(t *testing.T) {
	expectPasses(t, "email", Email()("jane@example.com"))
	expectFails(t, "email bad", Email()("not-an-email"), CodeEmail)

	expectPasses(t, "url", URL()("https://example.com/path"))
	expectFails(t, "url relative", URL()("/just/a/path"), CodeURL)
	expectFails(t, "url schemeless", URL()("example.com"), CodeURL)

	expectPasses(t, "uuid", UUID()("123e4567-e89b-12d3-a456-426614174000"))
	expectFails(t, "uuid bad", UUID()("123"), CodeUUID)

	expectPasses(t, "ip4", IP()("192.168.0.1"))
	expectPasses(t, "ip6", IP()("::1"))
	expectFails(t, "ip bad", IP()("999.1.1.1"), CodeIP)
	expectPasses(t, "ipv4", IPv4()("10.0.0.1"))
	expectFails(t, "ipv4 with v6", IPv4()("::1"), CodeIPv4)
	expectPasses(t, "ipv6", IPv6()("2001:db8::1"))
	expectFails(t, "ipv6 with v4", IPv6()("10.0.0.1"), CodeIPv6)

	expectPasses(t, "dns label", DNSLabel()("my-name-1"))
	expectFails(t, "dns label upper", DNSLabel()("MyName"), CodeDNSLabel)
	expectPasses(t, "dns sub", DNSSubdomain()("a.b.example.com"))
	expectFails(t, "dns sub empty", DNSSubdomain()(""), CodeDNSSubdomain)

	expectPasses(t, "base64", Base64()("aGVsbG8="))
	expectFails(t, "base64 bad", Base64()("!not base64!"), CodeBase64)

	expectPasses(t, "hex", Hexadecimal()("deadBEEF"))
	expectFails(t, "hex odd", Hexadecimal()("abc"), CodeHex)

	expectPasses(t, "semver", Semver()("1.2.3"))
	expectPasses(t, "semver pre", Semver()("1.2.3-rc.1+build.5"))
	expectFails(t, "semver bad", Semver()("1.2"), CodeSemver)
	expectFails(t, "semver leading zero", Semver()("01.2.3"), CodeSemver)
}

func TestComparable(t *testing.T) {
	expectPasses(t, "oneof hit", OneOf("a", "b")("b"))
	expectFails(t, "oneof miss", OneOf("a", "b")("c"), CodeOneOf)
	expectPasses(t, "oneof empty set", OneOf[string]()("anything"))
	expectFails(t, "noneof hit", NoneOf("x")("x"), CodeNoneOf)

	expectFails(t, "eq", EQ(5)(4), CodeEQ)
	expectPasses(t, "eq ok", EQ(5)(5))
	expectFails(t, "neq", NEQ(5)(5), CodeNEQ)

	expectFails(t, "gt", GT(5)(5), CodeGT)
	expectPasses(t, "gt ok", GT(5)(6))
	expectFails(t, "gte", GTE(5)(4), CodeGTE)
	expectPasses(t, "gte eq", GTE(5)(5))
	expectFails(t, "lt", LT(5)(5), CodeLT)
	expectFails(t, "lte", LTE(5)(6), CodeLTE)
	expectFails(t, "range low", InRange(1, 10)(0), CodeGTE)
	expectFails(t, "range high", InRange(1, 10)(11), CodeLTE)
	expectPasses(t, "range ok", InRange(1, 10)(10))
}

func TestCollections(t *testing.T) {
	expectPasses(t, "unique", Unique[string]()([]string{"a", "b"}))
	expectFails(t, "unique dup", Unique[string]()([]string{"a", "a"}), CodeUnique)
	expectPasses(t, "uniqueby", UniqueBy(func(s string) int { return len(s) })([]string{"a", "bb"}))
	expectFails(t, "uniqueby dup", UniqueBy(func(s string) int { return len(s) })([]string{"a", "b"}), CodeUnique)

	expectPasses(t, "contains all", SliceContainsAll("a", "b")([]string{"a", "b", "c"}))
	expectFails(t, "contains all miss", SliceContainsAll("a", "z")([]string{"a"}), CodeSliceContains)
	expectFails(t, "contains none hit", SliceContainsNone("z")([]string{"z"}), CodeSliceContains)

	m := map[string]int{"a": 1}
	expectPasses(t, "map key", MapKeyPresent[string, int]("a")(m))
	expectFails(t, "map key miss", MapKeyPresent[string, int]("b")(m), CodeMapKeyPresent)
	expectFails(t, "map value zero", MapValueRequired(0, "b")(m), CodeRequired)
	expectPasses(t, "map value set", MapValueRequired(0, "a")(m))
}

func TestCrossField(t *testing.T) {
	password := "s3cret"
	expectPasses(t, "equalto", EqualTo(func() string { return password })(password))
	expectFails(t, "equalto miss", EqualTo(func() string { return password })("other"), CodeCompareFields)
	expectFails(t, "notequalto hit", NotEqualTo(func() string { return password })(password), CodeCompareFields)

	start := 10
	expectPasses(t, "gtfield", GreaterThanField(func() int { return start })(11))
	expectFails(t, "gtfield miss", GreaterThanField(func() int { return start })(10), CodeCompareFields)
	expectFails(t, "ltfield miss", LessThanField(func() int { return start })(10), CodeCompareFields)

	cond := true
	expectFails(t, "requiredwhen", RequiredWhen("", func() bool { return cond })(""), CodeRequired)
	cond = false
	expectPasses(t, "requiredwhen off", RequiredWhen("", func() bool { return cond })(""))

	empty := ""
	set := "x"
	expectFails(t, "atleastone none", AtLeastOneSet(func() string { return empty }, func() string { return empty })(empty), CodeCrossField)
	expectPasses(t, "atleastone some", AtLeastOneSet(func() string { return empty }, func() string { return set })(empty))
	expectFails(t, "mutex both", MutuallyExclusive(func() string { return set }, func() string { return set })(set), CodeCrossField)
	expectPasses(t, "mutex one", MutuallyExclusive(func() string { return set }, func() string { return empty })(set))

	expectPasses(t, "deepequal", DeepEqual([]int{1, 2})([]int{1, 2}))
	expectFails(t, "deepequal miss", DeepEqual([]int{1, 2})([]int{1, 3}), CodeEQ)
}

func TestRulesComposeInControl(t *testing.T) {
	// Rules are plain ValidatorFuncs: they drop into NewControl and their codes
	// survive the engine's errors.Join aggregation.
	control := fform.NewControl(
		fform.NewPlainValueStore(""),
		"",
		StringNotEmpty(),
		MinLength(3),
		Alphanumeric(),
	)

	if control.Validate() {
		t.Fatal("expected invalid")
	}
	// Empty input fails both StringNotEmpty and MinLength; both codes survive
	// the engine's errors.Join aggregation, sorted and deduped.
	codes := fform.CodesOf(errors.Join(control.ValidationErrors()...))
	if len(codes) != 2 || codes[0] != CodeMinLength || codes[1] != CodeRequired {
		t.Fatalf("expected [min_length required] for empty value, got %v", codes)
	}

	control.Set("ab")
	control.Validate()
	codes = fform.CodesOf(errors.Join(control.ValidationErrors()...))
	if len(codes) != 1 || codes[0] != CodeMinLength {
		t.Fatalf("expected [min_length] for \"ab\", got %v", codes)
	}

	control.Set("abc")
	if !control.Validate() {
		t.Fatalf("expected valid for \"abc\", errors=%v", control.Errors())
	}
}

func TestGroupRules(t *testing.T) {
	// A group with two string controls and an int pair for ordering.
	build := func() *fform.Group {
		return fform.NewGroup(map[string]fform.FormNode{
			"email": fform.NewControl(fform.NewPlainValueStore(""), ""),
			"phone": fform.NewControl(fform.NewPlainValueStore(""), ""),
			"sms":   fform.NewControl(fform.NewPlainValueStore(""), ""),
			"start": fform.NewControl(fform.NewPlainValueStore(0), 0),
			"end":   fform.NewControl(fform.NewPlainValueStore(0), 0),
		})
	}
	set := func(g *fform.Group, key, value string) {
		ctl, ok := fform.ControlOf[string](fform.NodeResolver(g), key)
		if !ok {
			t.Fatalf("could not resolve %q", key)
		}
		ctl.Set(value)
	}
	setInt := func(g *fform.Group, key string, value int) {
		ctl, ok := fform.ControlOf[int](fform.NodeResolver(g), key)
		if !ok {
			t.Fatalf("could not resolve %q", key)
		}
		ctl.Set(value)
	}

	t.Run("AtLeastOneSet", func(t *testing.T) {
		g := build()
		v := GroupAtLeastOneSet[string]("email", "phone")
		if err := v(g); err == nil {
			t.Fatal("expected failure with both empty")
		} else if code, _ := fform.CodeOf(err); code != CodeCrossField {
			t.Fatalf("expected %q, got %q", CodeCrossField, code)
		}
		set(g, "phone", "555")
		if err := v(g); err != nil {
			t.Fatalf("expected pass once set, got %v", err)
		}
	})

	t.Run("MutuallyExclusive", func(t *testing.T) {
		g := build()
		v := GroupMutuallyExclusive[string]("email", "phone", "sms")
		if err := v(g); err != nil {
			t.Fatalf("expected pass when none set, got %v", err)
		}
		set(g, "email", "a@b.c")
		if err := v(g); err != nil {
			t.Fatalf("expected pass with one set, got %v", err)
		}
		set(g, "sms", "555")
		if err := v(g); err == nil {
			t.Fatal("expected failure with two set")
		}
	})

	t.Run("RequiredTogether", func(t *testing.T) {
		g := build()
		v := GroupRequiredTogether[string]("email", "phone")
		if err := v(g); err != nil {
			t.Fatalf("expected pass when neither set, got %v", err)
		}
		set(g, "email", "a@b.c")
		if err := v(g); err == nil {
			t.Fatal("expected failure with only one of the pair set")
		}
		set(g, "phone", "555")
		if err := v(g); err != nil {
			t.Fatalf("expected pass with both set, got %v", err)
		}
	})

	t.Run("Ordered", func(t *testing.T) {
		g := build()
		v := GroupOrdered[int]("end", "start")
		setInt(g, "start", 5)
		setInt(g, "end", 3)
		if err := v(g); err == nil {
			t.Fatal("expected failure when end <= start")
		}
		setInt(g, "end", 6)
		if err := v(g); err != nil {
			t.Fatalf("expected pass when end > start, got %v", err)
		}
	})

	t.Run("SkipsDisabledAndMissing", func(t *testing.T) {
		g := build()
		set(g, "email", "a@b.c")
		email, _ := fform.ControlOf[string](fform.NodeResolver(g), "email")
		email.SetDisabled(true)
		// Disabled email is ignored; nothing else set -> AtLeastOneSet fails.
		v := GroupAtLeastOneSet[string]("email", "phone", "missing")
		if err := v(g); err == nil {
			t.Fatal("expected failure when the only set control is disabled")
		}
	})

	t.Run("GroupCustom", func(t *testing.T) {
		g := build()
		v := GroupCustom("custom", func(g *fform.Group) bool {
			email, _ := fform.ControlOf[string](fform.NodeResolver(g), "email")
			return email.Value() != ""
		}, "email is required")
		if err := v(g); err == nil {
			t.Fatal("expected failure")
		}
		set(g, "email", "a@b.c")
		if err := v(g); err != nil {
			t.Fatalf("expected pass, got %v", err)
		}
	})
}
