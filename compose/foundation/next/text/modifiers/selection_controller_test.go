package modifiers

import (
	"testing"

	"github.com/zodimo/go-compose/compose/foundation/next/text/selection"
	"github.com/zodimo/go-compose/compose/ui/geometry"
	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/layout"
	"github.com/zodimo/go-compose/compose/ui/next/text"
	"github.com/zodimo/go-compose/compose/ui/next/text/style"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

// Mock implementations for testing

type mockLayoutCoordinates struct {
	attached bool
}

func (m *mockLayoutCoordinates) IsAttached() bool {
	return m.attached
}

func (m *mockLayoutCoordinates) Size() unit.IntSize {
	return unit.IntSize{}
}

func (m *mockLayoutCoordinates) PositionInRoot() geometry.Offset {
	return geometry.OffsetZero
}

func (m *mockLayoutCoordinates) PositionInWindow() geometry.Offset {
	return geometry.OffsetZero
}

func (m *mockLayoutCoordinates) LocalPositionOf(source layout.LayoutCoordinates, relativeToSource geometry.Offset) geometry.Offset {
	return geometry.OffsetZero
}

func (m *mockLayoutCoordinates) VisibleBounds() geometry.Rect {
	return geometry.RectZero
}

func (m *mockLayoutCoordinates) BoundsInWindow() geometry.Rect {
	return geometry.RectZero
}

// mockSelectable implements the full selection.Selectable interface
type mockSelectable struct {
	selectableId      int64
	lastVisibleOffset int
}

func (m *mockSelectable) SelectableId() int64 {
	return m.selectableId
}

func (m *mockSelectable) AppendSelectableInfoToBuilder(builder selection.SelectionLayoutBuilder) {}

func (m *mockSelectable) GetSelectAllSelection() *selection.Selection {
	return nil
}

func (m *mockSelectable) GetHandlePosition(sel selection.Selection, isStartHandle bool) geometry.Offset {
	return geometry.OffsetZero
}

func (m *mockSelectable) GetLayoutCoordinates() layout.LayoutCoordinates {
	return nil
}

func (m *mockSelectable) TextLayoutResult() *text.TextLayoutResult {
	return nil
}

func (m *mockSelectable) GetText() text.AnnotatedString {
	return text.NewAnnotatedString("", nil, nil)
}

func (m *mockSelectable) GetBoundingBox(offset int) geometry.Rect {
	return geometry.RectZero
}

func (m *mockSelectable) GetLineLeft(offset int) float32 {
	return 0
}

func (m *mockSelectable) GetLineRight(offset int) float32 {
	return 0
}

func (m *mockSelectable) GetCenterYForOffset(offset int) float32 {
	return 0
}

func (m *mockSelectable) GetRangeOfLineContaining(offset int) text.TextRange {
	return text.TextRangeZero
}

func (m *mockSelectable) GetLastVisibleOffset() int {
	return m.lastVisibleOffset
}

func (m *mockSelectable) GetLineHeight(offset int) float32 {
	return 0
}

type mockSelectionRegistrar struct {
	subscribed      bool
	unsubscribed    bool
	selectableId    int64
	positionChanged bool
	contentChanged  bool
	subselections   map[int64]*selection.Selection
	nextId          int64
}

func newMockSelectionRegistrar() *mockSelectionRegistrar {
	return &mockSelectionRegistrar{
		subselections: make(map[int64]*selection.Selection),
		nextId:        1,
	}
}

func (m *mockSelectionRegistrar) Subscribe(selectable selection.Selectable) selection.Selectable {
	m.subscribed = true
	m.selectableId = selectable.SelectableId()
	return &mockSelectable{selectableId: selectable.SelectableId(), lastVisibleOffset: 100}
}

func (m *mockSelectionRegistrar) Unsubscribe(selectable selection.Selectable) {
	m.unsubscribed = true
}

func (m *mockSelectionRegistrar) NextSelectableId() int64 {
	id := m.nextId
	m.nextId++
	return id
}

func (m *mockSelectionRegistrar) NotifySelectableChange(selectableId int64) {
	m.contentChanged = true
}

func (m *mockSelectionRegistrar) NotifyPositionChange(selectableId int64) {
	m.positionChanged = true
}

func (m *mockSelectionRegistrar) Subselections() map[int64]*selection.Selection {
	return m.subselections
}

func (m *mockSelectionRegistrar) NotifySelectionUpdateStart(
	layoutCoordinates layout.LayoutCoordinates,
	startPosition geometry.Offset,
	adjustment selection.SelectionAdjustment,
	isInTouchMode bool,
) {
}

func (m *mockSelectionRegistrar) NotifySelectionUpdateSelectAll(selectableId int64, isInTouchMode bool) {
}

func (m *mockSelectionRegistrar) NotifySelectionUpdate(
	layoutCoordinates layout.LayoutCoordinates,
	newPosition geometry.Offset,
	previousPosition geometry.Offset,
	isStartHandle bool,
	adjustment selection.SelectionAdjustment,
	isInTouchMode bool,
) bool {
	return true
}

func (m *mockSelectionRegistrar) NotifySelectionUpdateEnd() {
}

// Tests for StaticTextSelectionParams

func TestEmptyStaticTextSelectionParams(t *testing.T) {
	params := EmptyStaticTextSelectionParams()

	if params.LayoutCoordinatesValue() != nil {
		t.Error("Expected nil LayoutCoordinates")
	}
	if params.TextLayoutResultValue() != nil {
		t.Error("Expected nil TextLayoutResult")
	}
}

func TestStaticTextSelectionParams_GetPathForRange_NilResult(t *testing.T) {
	params := EmptyStaticTextSelectionParams()

	path := params.GetPathForRange(0, 10)
	if path != nil {
		t.Error("Expected nil path when TextLayoutResult is nil")
	}
}

func TestStaticTextSelectionParams_ShouldClip_NilResult(t *testing.T) {
	params := EmptyStaticTextSelectionParams()

	if params.ShouldClip() {
		t.Error("Expected ShouldClip to return false when TextLayoutResult is nil")
	}
}

func TestStaticTextSelectionParams_Copy(t *testing.T) {
	coords := &mockLayoutCoordinates{attached: true}
	params := NewStaticTextSelectionParams(coords, nil)

	if params.LayoutCoordinatesValue() != coords {
		t.Error("Expected LayoutCoordinates to be set")
	}

	newCoords := &mockLayoutCoordinates{attached: false}
	copied := params.CopyWithLayoutCoordinates(newCoords)

	if copied.LayoutCoordinatesValue() != newCoords {
		t.Error("Expected copied LayoutCoordinates to be updated")
	}
	// Original should be unchanged
	if params.LayoutCoordinatesValue() != coords {
		t.Error("Expected original LayoutCoordinates to be unchanged")
	}
}

// Tests for SelectionController

func TestNewSelectionController(t *testing.T) {
	registrar := newMockSelectionRegistrar()
	sc := NewSelectionController(42, registrar, graphics.ColorBlue)

	if sc.SelectableId() != 42 {
		t.Errorf("Expected SelectableId 42, got %d", sc.SelectableId())
	}
	if sc.Modifier() == nil {
		t.Error("Expected non-nil Modifier")
	}
}

func TestSelectionController_OnRemembered(t *testing.T) {
	registrar := newMockSelectionRegistrar()
	sc := NewSelectionController(42, registrar, graphics.ColorBlue)

	sc.OnRemembered()

	if !registrar.subscribed {
		t.Error("Expected Subscribe to be called")
	}
	if registrar.selectableId != 42 {
		t.Errorf("Expected selectableId 42, got %d", registrar.selectableId)
	}
}

func TestSelectionController_OnForgotten(t *testing.T) {
	registrar := newMockSelectionRegistrar()
	sc := NewSelectionController(42, registrar, graphics.ColorBlue)

	sc.OnRemembered()
	sc.OnForgotten()

	if !registrar.unsubscribed {
		t.Error("Expected Unsubscribe to be called")
	}
}

func TestSelectionController_OnAbandoned(t *testing.T) {
	registrar := newMockSelectionRegistrar()
	sc := NewSelectionController(42, registrar, graphics.ColorBlue)

	sc.OnRemembered()
	sc.OnAbandoned()

	if !registrar.unsubscribed {
		t.Error("Expected Unsubscribe to be called on abandon")
	}
}

func TestSelectionController_UpdateGlobalPosition(t *testing.T) {
	registrar := newMockSelectionRegistrar()
	sc := NewSelectionController(42, registrar, graphics.ColorBlue)

	coords := &mockLayoutCoordinates{attached: true}
	sc.UpdateGlobalPosition(coords)

	if !registrar.positionChanged {
		t.Error("Expected NotifyPositionChange to be called")
	}
}

func createMockTextLayoutResult() *text.TextLayoutResult {
	annotatedString := text.NewAnnotatedString("Hello World", nil, nil)
	layoutInput := text.NewTextLayoutInput(
		annotatedString,
		text.TextStyle{},
		nil,  // placeholders
		1,    // maxLines
		true, // softWrap
		style.OverFlowClip,
		nil, // density
		0,   // layoutDirection
		nil, // fontFamilyResolver
		unit.NewConstraints(0, 1000, 0, 1000),
	)
	result := text.NewTextLayoutResult(
		layoutInput,
		nil, // multiParagraph
		unit.IntSize{Width: 100, Height: 20},
	)
	return &result
}

func TestSelectionController_UpdateTextLayout_FirstUpdate(t *testing.T) {
	registrar := newMockSelectionRegistrar()
	sc := NewSelectionController(42, registrar, graphics.ColorBlue)

	result := createMockTextLayoutResult()
	sc.UpdateTextLayout(result)

	// First update should not notify change
	if registrar.contentChanged {
		t.Error("Expected no content change notification on first update")
	}
}

func TestSelectionController_UpdateTextLayout_SameText(t *testing.T) {
	registrar := newMockSelectionRegistrar()
	sc := NewSelectionController(42, registrar, graphics.ColorBlue)

	result1 := createMockTextLayoutResult()
	sc.UpdateTextLayout(result1)

	result2 := createMockTextLayoutResult() // Same text
	sc.UpdateTextLayout(result2)

	// Same text should not notify change
	if registrar.contentChanged {
		t.Error("Expected no content change notification when text is same")
	}
}

func TestSelectionController_Draw_NoSelection(t *testing.T) {
	registrar := newMockSelectionRegistrar()
	sc := NewSelectionController(42, registrar, graphics.ColorBlack)

	drawCalled := false
	sc.Draw(func(path graphics.Path, color graphics.Color, shouldClip bool) {
		drawCalled = true
	})

	if drawCalled {
		t.Error("Expected draw not to be called when no selection exists")
	}
}

func TestSelectionController_Draw_WithSelection(t *testing.T) {
	registrar := newMockSelectionRegistrar()
	registrar.subselections[42] = &selection.Selection{
		Start:          selection.AnchorInfo{Offset: 5},
		End:            selection.AnchorInfo{Offset: 10},
		HandlesCrossed: false,
	}

	sc := NewSelectionController(42, registrar, graphics.ColorBlue)
	sc.OnRemembered()

	// Without a TextLayoutResult, GetPathForRange returns nil, so draw won't be called
	drawCalled := false
	sc.Draw(func(path graphics.Path, color graphics.Color, shouldClip bool) {
		drawCalled = true
	})

	// Draw shouldn't be called because we don't have a path
	if drawCalled {
		t.Error("Expected draw not to be called when path is nil")
	}
}

func TestSelectionController_Draw_SameStartEnd(t *testing.T) {
	registrar := newMockSelectionRegistrar()
	registrar.subselections[42] = &selection.Selection{
		Start:          selection.AnchorInfo{Offset: 5},
		End:            selection.AnchorInfo{Offset: 5},
		HandlesCrossed: false,
	}

	sc := NewSelectionController(42, registrar, graphics.ColorBlue)

	drawCalled := false
	sc.Draw(func(path graphics.Path, color graphics.Color, shouldClip bool) {
		drawCalled = true
	})

	if drawCalled {
		t.Error("Expected draw not to be called when start == end")
	}
}

func TestSelectionController_ImplementsRememberObserver(t *testing.T) {
	var _ RememberObserver = (*SelectionController)(nil)
}

func TestMin(t *testing.T) {
	if min(5, 10) != 5 {
		t.Error("Expected min(5, 10) = 5")
	}
	if min(10, 5) != 5 {
		t.Error("Expected min(10, 5) = 5")
	}
	if min(5, 5) != 5 {
		t.Error("Expected min(5, 5) = 5")
	}
}
