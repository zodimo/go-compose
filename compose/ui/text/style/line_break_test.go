package style

import (
	"testing"
)

func TestLineBreak_Packing(t *testing.T) {
	strategy := LineBreakStrategyBalanced
	strictness := LineBreakStrictnessLoose
	wordBreak := LineBreakWordBreakPhrase

	lb := LineBreakOf(strategy, strictness, wordBreak)

	if lb.Strategy() != strategy {
		t.Errorf("Expected Strategy %v, got %v", strategy, lb.Strategy())
	}
	if lb.Strictness() != strictness {
		t.Errorf("Expected Strictness %v, got %v", strictness, lb.Strictness())
	}
	if lb.WordBreak() != wordBreak {
		t.Errorf("Expected WordBreak %v, got %v", wordBreak, lb.WordBreak())
	}
}

func TestLineBreak_Presets(t *testing.T) {
	// Simple
	if LineBreakSimple.Strategy() != LineBreakStrategySimple {
		t.Errorf("LineBreakSimple: Expected Strategy %v, got %v", LineBreakStrategySimple, LineBreakSimple.Strategy())
	}
	if LineBreakSimple.Strictness() != LineBreakStrictnessNormal {
		t.Errorf("LineBreakSimple: Expected Strictness %v, got %v", LineBreakStrictnessNormal, LineBreakSimple.Strictness())
	}
	if LineBreakSimple.WordBreak() != LineBreakWordBreakDefault {
		t.Errorf("LineBreakSimple: Expected WordBreak %v, got %v", LineBreakWordBreakDefault, LineBreakSimple.WordBreak())
	}

	// Heading
	if LineBreakHeading.Strategy() != LineBreakStrategyBalanced {
		t.Errorf("LineBreakHeading: Expected Strategy %v, got %v", LineBreakStrategyBalanced, LineBreakHeading.Strategy())
	}
	if LineBreakHeading.Strictness() != LineBreakStrictnessLoose {
		t.Errorf("LineBreakHeading: Expected Strictness %v, got %v", LineBreakStrictnessLoose, LineBreakHeading.Strictness())
	}
	if LineBreakHeading.WordBreak() != LineBreakWordBreakPhrase {
		t.Errorf("LineBreakHeading: Expected WordBreak %v, got %v", LineBreakWordBreakPhrase, LineBreakHeading.WordBreak())
	}

	// Paragraph
	if LineBreakParagraph.Strategy() != LineBreakStrategyHighQuality {
		t.Errorf("LineBreakParagraph: Expected Strategy %v, got %v", LineBreakStrategyHighQuality, LineBreakParagraph.Strategy())
	}
	if LineBreakParagraph.Strictness() != LineBreakStrictnessStrict {
		t.Errorf("LineBreakParagraph: Expected Strictness %v, got %v", LineBreakStrictnessStrict, LineBreakParagraph.Strictness())
	}
	if LineBreakParagraph.WordBreak() != LineBreakWordBreakDefault {
		t.Errorf("LineBreakParagraph: Expected WordBreak %v, got %v", LineBreakWordBreakDefault, LineBreakParagraph.WordBreak())
	}
}

func TestLineBreak_Unspecified(t *testing.T) {
	if LineBreakUnspecified.IsSpecified() {
		t.Error("LineBreakUnspecified should not be specified")
	}
	if LineBreakSimple.IsSpecified() == false {
		t.Error("LineBreakSimple should be specified")
	}
}

func TestLineBreak_Copy(t *testing.T) {
	orig := LineBreakSimple
	newLB := orig.Copy(LineBreakStrategyBalanced, orig.Strictness(), orig.WordBreak())

	if newLB.Strategy() != LineBreakStrategyBalanced {
		t.Errorf("Copy: Expected Strategy %v, got %v", LineBreakStrategyBalanced, newLB.Strategy())
	}
	if newLB.Strictness() != LineBreakStrictnessNormal {
		t.Errorf("Copy: Expected Strictness %v, got %v", LineBreakStrictnessNormal, newLB.Strictness())
	}
}

func TestLineBreak_String(t *testing.T) {
	s := LineBreakSimple.String()
	expected := "LineBreak(strategy=Strategy.Simple, strictness=Strictness.Normal, wordBreak=WordBreak.None)"
	if s != expected {
		t.Errorf("String(): Expected %q, got %q", expected, s)
	}
}
