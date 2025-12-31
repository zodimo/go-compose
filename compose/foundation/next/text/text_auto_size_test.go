package text

// import (
// 	"testing"

// 	"github.com/zodimo/go-compose/compose/ui/geometry"
// 	"github.com/zodimo/go-compose/compose/ui/graphics"
// 	uitext "github.com/zodimo/go-compose/compose/ui/next/text"
// 	"github.com/zodimo/go-compose/compose/ui/next/text/style"
// 	"github.com/zodimo/go-compose/compose/ui/unit"
// )

// var _ uitext.MultiParagraph = (*mockMultiParagraph)(nil)

// // Mock MultiParagraph
// type mockMultiParagraph struct {
// 	didOverflowWidth      bool
// 	didOverflowHeight     bool
// 	didExceedMaxLines     bool
// 	lineCount             int
// 	ellipsizedLines       map[int]bool
// 	resolvedTextDirection style.ResolvedTextDirection
// 	boundingBox           geometry.Rect
// 	cursorRect            geometry.Rect
// 	lineBaseline          float32
// 	pathForRange          graphics.Path
// }

// func (m *mockMultiParagraph) LineCount() int {
// 	return m.lineCount
// }

// func (m *mockMultiParagraph) IsLineEllipsized(lineIndex int) bool {
// 	return m.ellipsizedLines[lineIndex]
// }

// func (m *mockMultiParagraph) DidOverflowWidth() bool {
// 	return m.didOverflowWidth
// }

// func (m *mockMultiParagraph) DidOverflowHeight() bool {
// 	return m.didOverflowHeight
// }
// func (m *mockMultiParagraph) DidExceedMaxLines() bool {
// 	return m.didExceedMaxLines
// }

// // FirstBaseline
// func (m *mockMultiParagraph) FirstBaseline() float32 {
// 	return 0
// }

// // GetBidiRunDirection
// func (m *mockMultiParagraph) GetBidiRunDirection(offset int) style.ResolvedTextDirection {
// 	return m.resolvedTextDirection
// }
// func (m *mockMultiParagraph) GetBoundingBox(offset int) geometry.Rect {
// 	return m.boundingBox
// }

// func (m *mockMultiParagraph) GetCursorRect(offset int) geometry.Rect {
// 	return m.cursorRect
// }

// func (m *mockMultiParagraph) GetHorizontalPosition(offset int, usePrimaryDirection bool) float32 {
// 	return 0
// }

// func (m *mockMultiParagraph) GetLineBaseline(lineIndex int) float32 {
// 	return m.lineBaseline
// }

// func (m *mockMultiParagraph) GetLineBottom(lineIndex int) float32 {
// 	return 0
// }

// func (m *mockMultiParagraph) GetLineEnd(lineIndex int, visibleEnd bool) int {
// 	return 0
// }

// func (m *mockMultiParagraph) GetLineForOffset(offset int) int {
// 	return 0
// }
// func (m *mockMultiParagraph) GetLineForVerticalPosition(vertical float32) int {
// 	return 0
// }
// func (m *mockMultiParagraph) GetLineLeft(lineIndex int) float32 {
// 	return 0
// }
// func (m *mockMultiParagraph) GetLineRight(lineIndex int) float32 {
// 	return 0
// }
// func (m *mockMultiParagraph) GetLineTop(lineIndex int) float32 {
// 	return 0
// }
// func (m *mockMultiParagraph) GetLineWidth(lineIndex int) float32 {
// 	return 0
// }
// func (m *mockMultiParagraph) GetLineVisibleEnd(lineIndex int) int {
// 	return 0
// }
// func (m *mockMultiParagraph) GetLineStart(lineIndex int) int {
// 	return 0
// }
// func (m *mockMultiParagraph) GetOffsetForPosition(position geometry.Offset) int {
// 	return 0
// }
// func (m *mockMultiParagraph) GetParagraphDirection(offset int) style.ResolvedTextDirection {
// 	return m.resolvedTextDirection
// }
// func (m *mockMultiParagraph) GetPathForRange(start, end int) graphics.Path {
// 	return m.pathForRange
// }

// // LastBaseline
// // func (m *mockMultiParagraph) LastBaseline() float32 {
// // 	return 0
// // }

// // Mock Scope
// type mockTextAutoSizeLayoutScope struct {
// 	maxSize float32
// }

// func (m *mockTextAutoSizeLayoutScope) PerformLayout(constraints uiConstraints, text uiAnnotatedString, fontSize uiTextUnit) uiTextLayoutResult {
// 	// Simple mock logic: if memory size (value) > maxSize, it overflows width
// 	overflow := fontSize.Value() > m.maxSize

// 	mp := &mockMultiParagraph{
// 		didOverflowWidth: overflow,
// 		lineCount:        1,
// 	}

// 	input := uitext.TextLayoutInput{
// 		Overflow: style.OverFlowClip,
// 	}

// 	return uitext.NewTextLayoutResult(input, mp, unit.IntSize{})
// }

// func TestStepBasedTextAutoSize(t *testing.T) {
// 	// Setup
// 	minSize := uiTextUnitSp(10)
// 	maxSize := uiTextUnitSp(100)
// 	step := uiTextUnitSp(10)

// 	autoSize := NewStepBasedTextAutoSize(minSize, maxSize, step)

// 	// Case 1: Max size fits
// 	// Mock layout that fits everything <= 100
// 	scope := &mockTextAutoSizeLayoutScope{maxSize: 100}
// 	result := autoSize.GetFontSize(scope, 0, uiAnnotatedString{})
// 	if result.Value() != 100 {
// 		t.Errorf("Expected 100, got %f", result.Value())
// 	}

// 	// Case 2: Fits at 50
// 	scope = &mockTextAutoSizeLayoutScope{maxSize: 55} // 60 overflows, 50 fits
// 	result = autoSize.GetFontSize(scope, 0, uiAnnotatedString{})
// 	if result.Value() != 50 {
// 		t.Errorf("Expected 50, got %f", result.Value())
// 	}

// 	// Case 3: Does not fit even at min (10), should return min
// 	scope = &mockTextAutoSizeLayoutScope{maxSize: 5}
// 	result = autoSize.GetFontSize(scope, 0, uiAnnotatedString{})
// 	if result.Value() != 10 {
// 		t.Errorf("Expected 10, got %f", result.Value())
// 	}
// }
