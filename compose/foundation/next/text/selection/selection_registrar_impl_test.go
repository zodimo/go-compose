package selection

import (
	"testing"

	"github.com/zodimo/go-compose/compose/ui/geometry"
	"github.com/zodimo/go-compose/compose/ui/layout"
	"github.com/zodimo/go-compose/compose/ui/next/text"
)

func TestNewSelectionRegistrarImpl(t *testing.T) {
	registrar := NewSelectionRegistrarImpl()

	if registrar == nil {
		t.Fatal("Expected non-nil registrar")
	}
	if len(registrar.Selectables()) != 0 {
		t.Error("Expected empty selectables list")
	}
	if len(registrar.SelectableMap()) != 0 {
		t.Error("Expected empty selectable map")
	}
	if len(registrar.Subselections()) != 0 {
		t.Error("Expected empty subselections map")
	}
}

func TestSelectionRegistrarImpl_NextSelectableId(t *testing.T) {
	registrar := NewSelectionRegistrarImpl()

	id1 := registrar.NextSelectableId()
	id2 := registrar.NextSelectableId()
	id3 := registrar.NextSelectableId()

	if id1 == InvalidSelectableId {
		t.Error("First ID should not be invalid")
	}
	if id2 == InvalidSelectableId {
		t.Error("Second ID should not be invalid")
	}
	if id3 == InvalidSelectableId {
		t.Error("Third ID should not be invalid")
	}

	// IDs should be unique
	ids := map[int64]bool{id1: true, id2: true, id3: true}
	if len(ids) != 3 {
		t.Error("IDs should be unique")
	}
}

// mockSelectable for testing
type mockSelectable struct {
	selectableId int64
}

func (m *mockSelectable) SelectableId() int64 {
	return m.selectableId
}

func (m *mockSelectable) AppendSelectableInfoToBuilder(builder SelectionLayoutBuilder) {}
func (m *mockSelectable) GetSelectAllSelection() *Selection                            { return nil }
func (m *mockSelectable) GetHandlePosition(selection Selection, isStartHandle bool) geometry.Offset {
	return geometry.OffsetZero
}
func (m *mockSelectable) GetLayoutCoordinates() layout.LayoutCoordinates { return nil }
func (m *mockSelectable) TextLayoutResult() *text.TextLayoutResult       { return nil }
func (m *mockSelectable) GetText() text.AnnotatedString {
	return text.NewAnnotatedString("", nil, nil)
}
func (m *mockSelectable) GetBoundingBox(offset int) geometry.Rect { return geometry.RectZero }
func (m *mockSelectable) GetLineLeft(offset int) float32          { return 0 }
func (m *mockSelectable) GetLineRight(offset int) float32         { return 0 }
func (m *mockSelectable) GetCenterYForOffset(offset int) float32  { return 0 }
func (m *mockSelectable) GetRangeOfLineContaining(offset int) text.TextRange {
	return text.TextRangeZero
}
func (m *mockSelectable) GetLastVisibleOffset() int        { return 0 }
func (m *mockSelectable) GetLineHeight(offset int) float32 { return 0 }

func TestSelectionRegistrarImpl_SubscribeUnsubscribe(t *testing.T) {
	registrar := NewSelectionRegistrarImpl()
	selectable := &mockSelectable{selectableId: 42}

	// Subscribe
	result := registrar.Subscribe(selectable)
	if result != selectable {
		t.Error("Subscribe should return the same selectable")
	}

	selectables := registrar.Selectables()
	if len(selectables) != 1 {
		t.Errorf("Expected 1 selectable, got %d", len(selectables))
	}

	selectableMap := registrar.SelectableMap()
	if _, exists := selectableMap[42]; !exists {
		t.Error("Expected selectable to be in map")
	}

	// Unsubscribe
	registrar.Unsubscribe(selectable)

	selectables = registrar.Selectables()
	if len(selectables) != 0 {
		t.Errorf("Expected 0 selectables after unsubscribe, got %d", len(selectables))
	}

	selectableMap = registrar.SelectableMap()
	if _, exists := selectableMap[42]; exists {
		t.Error("Expected selectable to be removed from map")
	}
}

func TestSelectionRegistrarImpl_Subscribe_DuplicatePanics(t *testing.T) {
	registrar := NewSelectionRegistrarImpl()
	selectable := &mockSelectable{selectableId: 42}

	registrar.Subscribe(selectable)

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when subscribing duplicate")
		}
	}()

	// This should panic
	registrar.Subscribe(selectable)
}

func TestSelectionRegistrarImpl_Subscribe_InvalidIdPanics(t *testing.T) {
	registrar := NewSelectionRegistrarImpl()
	selectable := &mockSelectable{selectableId: InvalidSelectableId}

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when subscribing with invalid ID")
		}
	}()

	// This should panic
	registrar.Subscribe(selectable)
}

func TestSelectionRegistrarImpl_SetSubselections(t *testing.T) {
	registrar := NewSelectionRegistrarImpl()

	subs := map[int64]*Selection{
		1: {Start: AnchorInfo{Offset: 5}, End: AnchorInfo{Offset: 10}},
	}
	registrar.SetSubselections(subs)

	result := registrar.Subselections()
	if len(result) != 1 {
		t.Errorf("Expected 1 subselection, got %d", len(result))
	}
	if _, exists := result[1]; !exists {
		t.Error("Expected subselection with key 1")
	}
}

func TestSelectionRegistrarImpl_NotifyCallbacks(t *testing.T) {
	registrar := NewSelectionRegistrarImpl()

	positionChanged := false
	selectableChanged := false
	updateEnded := false

	registrar.OnPositionChangeCallback = func(id int64) {
		positionChanged = true
	}
	registrar.OnSelectableChangeCallback = func(id int64) {
		selectableChanged = true
	}
	registrar.OnSelectionUpdateEndCallback = func() {
		updateEnded = true
	}

	registrar.NotifyPositionChange(42)
	if !positionChanged {
		t.Error("Expected position change callback to be called")
	}

	registrar.NotifySelectableChange(42)
	if !selectableChanged {
		t.Error("Expected selectable change callback to be called")
	}

	registrar.NotifySelectionUpdateEnd()
	if !updateEnded {
		t.Error("Expected selection update end callback to be called")
	}
}

func TestSelectionRegistrarImpl_Unsubscribe_AfterCallback(t *testing.T) {
	registrar := NewSelectionRegistrarImpl()
	selectable := &mockSelectable{selectableId: 42}

	unsubscribedId := int64(-1)
	registrar.AfterSelectableUnsubscribe = func(id int64) {
		unsubscribedId = id
	}

	registrar.Subscribe(selectable)
	registrar.Unsubscribe(selectable)

	if unsubscribedId != 42 {
		t.Errorf("Expected AfterSelectableUnsubscribe to be called with 42, got %d", unsubscribedId)
	}
}

func TestHasSelection(t *testing.T) {
	// Test with nil registrar
	if HasSelection(nil, 42) {
		t.Error("Expected HasSelection to return false for nil registrar")
	}

	// Test with no selection
	registrar := NewSelectionRegistrarImpl()
	if HasSelection(registrar, 42) {
		t.Error("Expected HasSelection to return false when no selection exists")
	}

	// Test with selection
	registrar.SetSubselections(map[int64]*Selection{
		42: {Start: AnchorInfo{Offset: 5}, End: AnchorInfo{Offset: 10}},
	})
	if !HasSelection(registrar, 42) {
		t.Error("Expected HasSelection to return true when selection exists")
	}
	if HasSelection(registrar, 99) {
		t.Error("Expected HasSelection to return false for different ID")
	}
}
