package text

import (
	"testing"

	"github.com/zodimo/go-compose/compose/ui/text/style"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

func TestNormalizedParagraphStyles(t *testing.T) {
	// Test case: 0-5 (para1), 5-10 (para2)
	// Expect: 0-5, 5-10
	para1 := ParagraphStyle{TextAlign: style.TextAlignLeft}
	para2 := ParagraphStyle{TextAlign: style.TextAlignRight}
	defaultStyle := ParagraphStyle{TextAlign: style.TextAlignCenter}

	as := NewAnnotatedString("0123456789", nil, []Range[ParagraphStyle]{
		{Item: para1, Start: 0, End: 5},
		{Item: para2, Start: 5, End: 10},
	})

	normalized := as.NormalizedParagraphStyles(defaultStyle)
	if len(normalized) != 2 {
		t.Fatalf("Expected 2 paragraphs, got %d", len(normalized))
	}
	if normalized[0].Start != 0 || normalized[0].End != 5 {
		t.Errorf("Para 1 range mismatch: %v", normalized[0])
	}
	if normalized[1].Start != 5 || normalized[1].End != 10 {
		t.Errorf("Para 2 range mismatch: %v", normalized[1])
	}
}

func TestNormalizedParagraphStyles_Gaps(t *testing.T) {
	// Test case: 0-3 (para1), 7-10 (para2)
	// Gap 3-7 should be default
	para1 := ParagraphStyle{TextAlign: style.TextAlignLeft}
	para2 := ParagraphStyle{TextAlign: style.TextAlignRight}
	defaultStyle := ParagraphStyle{TextAlign: style.TextAlignCenter}

	as := NewAnnotatedString("0123456789", nil, []Range[ParagraphStyle]{
		{Item: para1, Start: 0, End: 3},
		{Item: para2, Start: 7, End: 10},
	})

	normalized := as.NormalizedParagraphStyles(defaultStyle)
	if len(normalized) != 3 {
		t.Fatalf("Expected 3 paragraphs, got %d", len(normalized))
	}
	// 0-3
	if normalized[0].Start != 0 || normalized[0].End != 3 {
		t.Errorf("Para 1 range mismatch: %v", normalized[0])
	}
	// 3-7 (gap)
	if normalized[1].Start != 3 || normalized[1].End != 7 {
		t.Errorf("Para 2 range mismatch: %v", normalized[1])
	}
	if normalized[1].Item.TextAlign != defaultStyle.TextAlign {
		t.Errorf("Gap style mismatch, expected default")
	}
	// 7-10
	if normalized[2].Start != 7 || normalized[2].End != 10 {
		t.Errorf("Para 3 range mismatch: %v", normalized[2])
	}
}

func TestBuilder_WithHelpers(t *testing.T) {
	b := NewBuilder()
	b.Append("Hello")
	b.WithStyle(SpanStyle{FontSize: unit.Sp(10)}, func() {
		b.Append(" World")
	})

	as := b.ToAnnotatedString()
	if as.Text() != "Hello World" {
		t.Errorf("Expected 'Hello World', got '%s'", as.Text())
	}
	styles := as.SpanStyles()
	if len(styles) != 1 {
		t.Fatalf("Expected 1 span style, got %d", len(styles))
	}
	if styles[0].Start != 5 || styles[0].End != 11 {
		t.Errorf("Style range mismatch: %v", styles[0])
	}
}

func TestMapAnnotations(t *testing.T) {
	as := NewAnnotatedString("Test", nil, nil)
	// Add annotation manually since NewAnnotatedString generic is strict on SpanStyles/ParaStyles args separate
	// We use builder to add generic annotation
	b := NewBuilderFromAnnotatedString(as)
	b.AddStringAnnotation("tag", "val", 0, 4)
	as = b.ToAnnotatedString()

	mapped := as.MapAnnotations(func(r Range[Annotation]) Range[Annotation] {
		if sa, ok := r.Item.(StringAnnotation); ok {
			return Range[Annotation]{
				Item:  StringAnnotation(string(sa) + "_mapped"),
				Start: r.Start,
				End:   r.End,
				Tag:   r.Tag,
			}
		}
		return r
	})

	anns := mapped.GetStringAnnotations("tag", 0, 4)
	if len(anns) != 1 {
		t.Fatalf("Expected 1 annotation, got %d", len(anns))
	}
	if anns[0].Item != "val_mapped" {
		t.Errorf("Expected 'val_mapped', got '%s'", anns[0].Item)
	}
}

func TestTransformText(t *testing.T) {
	as := NewAnnotatedString("Test", nil, nil)
	upper := as.ToUpper()
	if upper.Text() != "TEST" {
		t.Errorf("Expected TEST, got %s", upper.Text())
	}

	cap := as.Capitalize() // "Test" -> "Test" (already capped)
	if cap.Text() != "Test" {
		t.Errorf("Expected Test, got %s", cap.Text())
	}

	lower := as.ToLower()
	if lower.Text() != "test" {
		t.Errorf("Expected test, got %s", lower.Text())
	}

	decap := lower.Capitalize()
	if decap.Text() != "Test" {
		t.Errorf("Expected Test, got %s", decap.Text())
	}
}
