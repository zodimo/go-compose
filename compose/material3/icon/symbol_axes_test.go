package icon

import (
	"testing"

	"github.com/go-text/typesetting/font"
	"github.com/go-text/typesetting/font/opentype"
)

func optionsWith(opts ...IconOption) IconOptions {
	o := DefaultIconOptions()
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

func TestAxisOptions_MarkProvided(t *testing.T) {
	cases := []struct {
		name       string
		option     IconOption
		wantValue  float32
		providedBy func(IconOptions) bool
	}{
		{"fill", WithFill(0), 0, func(o IconOptions) bool { return o.hasFill }},
		{"grade", WithGrade(0), 0, func(o IconOptions) bool { return o.hasGrade }},
		{"optical size", WithOpticalSize(0), 0, func(o IconOptions) bool { return o.hasOpticalSz }},
		{"weight", WithWeight(0), 0, func(o IconOptions) bool { return o.hasWeight }},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := optionsWith(tc.option)
			if !tc.providedBy(got) {
				t.Fatalf("%s option did not mark the axis as provided", tc.name)
			}
		})
	}
}

func TestAxisOptions_UnprovidedAxesReportUnset(t *testing.T) {
	got := DefaultIconOptions()
	if got.hasFill || got.hasGrade || got.hasOpticalSz || got.hasWeight {
		t.Fatalf("default options claim a provided axis: fill=%v grade=%v opsz=%v weight=%v",
			got.hasFill, got.hasGrade, got.hasOpticalSz, got.hasWeight)
	}
}

func TestAxisOptions_ExplicitZeroIsProvided(t *testing.T) {
	// The whole reason for the has-flags: WithFill(0) and WithGrade(0) must not
	// collapse into "unset", or a caller could never force the default value.
	opts := optionsWith(WithFill(0), WithGrade(0))
	axes := resolveAxes(opts, 24)

	if axes.Fill != 0 || axes.Grade != 0 {
		t.Fatalf("axes = %+v, want fill 0 grade 0", axes)
	}
	if !opts.hasFill || !opts.hasGrade {
		t.Fatal("explicit zero was not recorded as provided")
	}
}

func TestResolveAxes_ClampsOutOfRange(t *testing.T) {
	// Grade 500 clamps to the 200 maximum.
	if got := resolveAxes(optionsWith(WithGrade(500)), 24).Grade; got != GradeMax {
		t.Errorf("grade 500 resolved to %v, want %v", got, GradeMax)
	}
	// Fill -1 clamps to the 0 minimum.
	if got := resolveAxes(optionsWith(WithFill(-1)), 24).Fill; got != FillMin {
		t.Errorf("fill -1 resolved to %v, want %v", got, FillMin)
	}
	if got := resolveAxes(optionsWith(WithWeight(5000)), 24).Weight; got != WeightMax {
		t.Errorf("weight 5000 resolved to %v, want %v", got, WeightMax)
	}
	if got := resolveAxes(optionsWith(WithWeight(1)), 24).Weight; got != WeightMin {
		t.Errorf("weight 1 resolved to %v, want %v", got, WeightMin)
	}
	if got := resolveAxes(optionsWith(WithOpticalSize(1000)), 24).OpticalSize; got != OpticalMax {
		t.Errorf("opsz 1000 resolved to %v, want %v", got, OpticalMax)
	}
	if got := resolveAxes(optionsWith(WithGrade(-1000)), 24).Grade; got != GradeMin {
		t.Errorf("grade -1000 resolved to %v, want %v", got, GradeMin)
	}
}

func TestResolveAxes_NaNDoesNotPoisonKey(t *testing.T) {
	nan := float32(0)
	nan = nan / nan
	if got := resolveAxes(optionsWith(WithFill(nan)), 24).Fill; got != FillMin {
		t.Fatalf("NaN fill resolved to %v, want %v", got, FillMin)
	}
}

func TestOpticalSize_AdaptsToRenderedSize(t *testing.T) {
	cases := []struct {
		name   string
		opts   IconOptions
		sizePx int
		want   float32
	}{
		{name: "small icon takes axis minimum", opts: DefaultIconOptions(), sizePx: 20, want: OpticalMin},
		{name: "16dp takes axis minimum", opts: DefaultIconOptions(), sizePx: 16, want: OpticalMin},
		{name: "24dp tracks size", opts: DefaultIconOptions(), sizePx: 24, want: 24},
		{name: "32dp sits between bounds", opts: DefaultIconOptions(), sizePx: 32, want: 32},
		{name: "48dp takes axis maximum", opts: DefaultIconOptions(), sizePx: 48, want: OpticalMax},
		{name: "64dp takes axis maximum", opts: DefaultIconOptions(), sizePx: 64, want: OpticalMax},
		{name: "explicit opsz wins over size", opts: optionsWith(WithOpticalSize(24)), sizePx: 48, want: 24},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := resolveAxes(tc.opts, tc.sizePx).OpticalSize; got != tc.want {
				t.Fatalf("opsz = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestResolveAxes_DefaultInstanceWhenNothingProvided(t *testing.T) {
	got := resolveAxes(DefaultIconOptions(), 24)
	want := defaultAxisValues()
	if got != want {
		t.Fatalf("default axes = %+v, want %+v", got, want)
	}
}

func TestVariations_OmitsDefaults(t *testing.T) {
	// A plain 24dp icon must stay on the font's default instance so it renders
	// exactly as before this change.
	if v := resolveAxes(DefaultIconOptions(), 24).variations(); len(v) != 0 {
		t.Fatalf("default instance produced %d variations, want 0", len(v))
	}
}

func TestVariations_IncludesOnlyChangedAxes(t *testing.T) {
	v := resolveAxes(optionsWith(WithFill(1)), 24).variations()
	if len(v) != 1 {
		t.Fatalf("got %d variations, want 1: %+v", len(v), v)
	}
	if v[0].Tag != mustTag(t, axisFill) || v[0].Value != 1 {
		t.Fatalf("variation = %+v, want FILL=1", v[0])
	}
}

func mustTag(t *testing.T, s string) font.Tag {
	t.Helper()
	return opentype.MustNewTag(s)
}
