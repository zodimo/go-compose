package modifiers

import (
	"testing"

	"github.com/zodimo/go-compose/compose/ui/graphics"
	"github.com/zodimo/go-compose/compose/ui/text"
	"github.com/zodimo/go-compose/compose/ui/text/style"
	"github.com/zodimo/go-compose/compose/ui/unit"
)

// Mock implementations for testing

type mockLayoutCoordinates struct {
	attached bool
}

func (m *mockLayoutCoordinates) IsAttached() bool {
	return m.attached
}

type mockSelectable struct {
	lastVisibleOffset int
}

func (m *mockSelectable) GetLastVisibleOffset() int {
	return m.lastVisibleOffset
}

type mockSelectionRegistrar struct {
	subscribed      bool
	unsubscribed    bool
	selectableId    int64
	positionChanged bool
	contentChanged  bool
	subselections   map[int64]*Selection
}

func newMockSelectionRegistrar() *mockSelectionRegistrar {
	return &mockSelectionRegistrar{
		subselections: make(map[int64]*Selection),
	}
}

func (m *mockSelectionRegistrar) Subscribe(delegate MultiWidgetSelectionDelegate) Selectable {
	m.subscribed = true
	m.selectableId = delegate.SelectableId
	return &mockSelectable{lastVisibleOffset: 100}
}

func (m *mockSelectionRegistrar) Unsubscribe(selectable Selectable) {
	m.unsubscribed = true
}

func (m *mockSelectionRegistrar) NotifySelectableChange(selectableId int64) {
	m.contentChanged = true
}

func (m *mockSelectionRegistrar) NotifyPositionChange(selectableId int64) {
	m.positionChanged = true
}

func (m *mockSelectionRegistrar) Subselections() map[int64]*Selection {
	return m.subselections
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

	newCoords := &mockLayoutCoordinates{attached: false}
	copied := params.Copy(newCoords, nil)

	if copied.LayoutCoordinatesValue() != newCoords {
		t.Error("Expected copied params to have new coordinates")
	}
	if copied.TextLayoutResultValue() != nil {
		t.Error("Expected copied params to have nil TextLayoutResult")
	}
}

func TestStaticTextSelectionParams_CopyWithLayoutCoordinates(t *testing.T) {
	params := EmptyStaticTextSelectionParams()
	coords := &mockLayoutCoordinates{attached: true}

	copied := params.CopyWithLayoutCoordinates(coords)

	if copied.LayoutCoordinatesValue() != coords {
		t.Error("Expected copied params to have new coordinates")
	}
}

// Tests for SelectionController

func TestNewSelectionController(t *testing.T) {
	registrar := newMockSelectionRegistrar()
	color := graphics.ColorBlack

	sc := NewSelectionController(123, registrar, color)

	if sc.SelectableId() != 123 {
		t.Errorf("Expected selectableId 123, got %d", sc.SelectableId())
	}
	if sc.Modifier() == nil {
		t.Error("Expected non-nil modifier")
	}
}

func TestSelectionController_OnRemembered(t *testing.T) {
	registrar := newMockSelectionRegistrar()
	sc := NewSelectionController(42, registrar, graphics.ColorBlack)

	sc.OnRemembered()

	if !registrar.subscribed {
		t.Error("Expected registrar.Subscribe to be called")
	}
	if registrar.selectableId != 42 {
		t.Errorf("Expected selectableId 42, got %d", registrar.selectableId)
	}
}

func TestSelectionController_OnForgotten(t *testing.T) {
	registrar := newMockSelectionRegistrar()
	sc := NewSelectionController(42, registrar, graphics.ColorBlack)

	sc.OnRemembered()
	sc.OnForgotten()

	if !registrar.unsubscribed {
		t.Error("Expected registrar.Unsubscribe to be called")
	}
}

func TestSelectionController_OnAbandoned(t *testing.T) {
	registrar := newMockSelectionRegistrar()
	sc := NewSelectionController(42, registrar, graphics.ColorBlack)

	sc.OnRemembered()
	sc.OnAbandoned()

	if !registrar.unsubscribed {
		t.Error("Expected registrar.Unsubscribe to be called")
	}
}

func TestSelectionController_UpdateTextLayout_NotifiesOnTextChange(t *testing.T) {
	registrar := newMockSelectionRegistrar()
	sc := NewSelectionController(42, registrar, graphics.ColorBlack)

	// Create first layout result
	input1 := text.NewTextLayoutInput(
		text.NewAnnotatedString("Hello", nil, nil),
		text.TextStyle{},
		nil,
		1,
		true,
		style.OverFlowClip,
		unit.NewDensity(1.0, 1.0),
		unit.LayoutDirectionLtr,
		nil,
		unit.NewConstraints(0, 100, 0, 100),
	)
	result1 := text.NewTextLayoutResult(input1, nil, unit.IntSize{})
	sc.UpdateTextLayout(&result1)

	// Create second layout result with different text
	input2 := text.NewTextLayoutInput(
		text.NewAnnotatedString("World", nil, nil),
		text.TextStyle{},
		nil,
		1,
		true,
		style.OverFlowClip,
		unit.NewDensity(1.0, 1.0),
		unit.LayoutDirectionLtr,
		nil,
		unit.NewConstraints(0, 100, 0, 100),
	)
	result2 := text.NewTextLayoutResult(input2, nil, unit.IntSize{})
	sc.UpdateTextLayout(&result2)

	if !registrar.contentChanged {
		t.Error("Expected NotifySelectableChange to be called when text changes")
	}
}

func TestSelectionController_UpdateGlobalPosition(t *testing.T) {
	registrar := newMockSelectionRegistrar()
	sc := NewSelectionController(42, registrar, graphics.ColorBlack)

	coords := &mockLayoutCoordinates{attached: true}
	sc.UpdateGlobalPosition(coords)

	if !registrar.positionChanged {
		t.Error("Expected NotifyPositionChange to be called")
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
	registrar.subselections[42] = &Selection{
		Start:          SelectionAnchorInfo{Offset: 5},
		End:            SelectionAnchorInfo{Offset: 10},
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
	registrar.subselections[42] = &Selection{
		Start:          SelectionAnchorInfo{Offset: 5},
		End:            SelectionAnchorInfo{Offset: 5},
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
