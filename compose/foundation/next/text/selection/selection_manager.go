package selection

import (
	"sync"

	"github.com/zodimo/go-compose/compose/ui/geometry"
	"github.com/zodimo/go-compose/compose/ui/layout"
	"github.com/zodimo/go-compose/compose/ui/next/text"
	"github.com/zodimo/go-compose/internal/modifier"
)

// Handle represents a selection handle type.
type Handle int

const (
	// HandleNone represents no handle.
	HandleNone Handle = iota
	// HandleSelectionStart represents the start handle.
	HandleSelectionStart
	// HandleSelectionEnd represents the end handle.
	HandleSelectionEnd
)

// HapticFeedback provides haptic feedback functionality.
// This is a platform-specific interface that should be implemented per platform.
type HapticFeedback interface {
	PerformHapticFeedback(feedbackType HapticFeedbackType)
}

// HapticFeedbackType represents the type of haptic feedback.
type HapticFeedbackType int

const (
	// HapticFeedbackTextHandleMove is the feedback for text handle movement.
	HapticFeedbackTextHandleMove HapticFeedbackType = iota
)

// TextToolbar provides text toolbar functionality (copy, paste, etc).
// This is a platform-specific interface that should be implemented per platform.
type TextToolbar interface {
	ShowMenu(rect geometry.Rect, onCopy func(), onSelectAll func())
	Hide()
	Status() TextToolbarStatus
}

// TextToolbarStatus represents the status of the text toolbar.
type TextToolbarStatus int

const (
	// TextToolbarStatusHidden means the toolbar is hidden.
	TextToolbarStatusHidden TextToolbarStatus = iota
	// TextToolbarStatusShown means the toolbar is shown.
	TextToolbarStatusShown
)

// SelectionManager is a bridge class between user interaction and text composables for selection.
//
// https://cs.android.com/androidx/platform/frameworks/support/+/androidx-main:compose/foundation/foundation/src/commonMain/kotlin/androidx/compose/foundation/text/selection/SelectionManager.kt
type SelectionManager struct {
	mu sync.RWMutex

	selectionRegistrar *SelectionRegistrarImpl

	// selection is the current selection state.
	selection *Selection

	// isInTouchMode indicates if touch mode is active.
	isInTouchMode bool

	// onSelectionChange is invoked when the selection should change.
	onSelectionChange func(*Selection)

	// hapticFeedback provides haptic feedback.
	HapticFeedback HapticFeedback

	// onCopyHandler handles copy operations.
	OnCopyHandler func(text.AnnotatedString)

	// textToolbar shows floating toolbar.
	TextToolbar TextToolbar

	// hasFocus indicates if the container has focus.
	hasFocus bool

	// containerLayoutCoordinates stores the layout coordinates.
	containerLayoutCoordinates layout.LayoutCoordinates
	previousPosition           *geometry.Offset

	// Handle positions
	startHandlePosition *geometry.Offset
	endHandlePosition   *geometry.Offset

	// Dragging state
	draggingHandle      Handle
	currentDragPosition *geometry.Offset
	dragBeginPosition   geometry.Offset
	dragTotalDistance   geometry.Offset

	// Handle line heights
	startHandleLineHeight float32
	endHandleLineHeight   float32

	// previousSelectionLayout for comparison
	previousSelectionLayout SelectionLayout

	// showToolbar indicates if toolbar should be shown.
	showToolbar bool

	// isLongPressOrClickSelection indicates selection was initiated by long press/click.
	isLongPressOrClickSelection bool

	// modifierCache stores the computed modifier.
	modifierCache modifier.Modifier
}

// NewSelectionManager creates a new SelectionManager with the given registrar.
func NewSelectionManager(registrar *SelectionRegistrarImpl) *SelectionManager {
	m := &SelectionManager{
		selectionRegistrar: registrar,
		isInTouchMode:      true,
		hasFocus:           false,
		draggingHandle:     HandleNone,
		showToolbar:        false,
	}

	// Set up callbacks
	registrar.OnPositionChangeCallback = func(selectableId int64) {
		if _, exists := registrar.Subselections()[selectableId]; exists {
			m.updateHandleOffsets()
			m.updateSelectionToolbar()
		}
	}

	registrar.OnSelectionUpdateStartCallback = func(isInTouchMode bool, layoutCoordinates layout.LayoutCoordinates, position geometry.Offset, adjustment SelectionAdjustment) {
		positionInContainer := m.convertToContainerCoordinates(layoutCoordinates, position)
		if positionInContainer.IsOffset() {
			m.mu.Lock()
			m.isInTouchMode = isInTouchMode
			m.mu.Unlock()

			m.startSelection(positionInContainer, false, adjustment)
			m.setShowToolbar(false)
		}
	}

	registrar.OnSelectionUpdateSelectAllCallback = func(isInTouchMode bool, selectableId int64) {
		m.mu.Lock()
		m.isInTouchMode = isInTouchMode
		m.mu.Unlock()
		m.selectAllInSelectable(selectableId)
		m.setShowToolbar(false)
	}

	registrar.OnSelectionUpdateCallback = func(isInTouchMode bool, layoutCoordinates layout.LayoutCoordinates, newPosition, previousPosition geometry.Offset, isStartHandle bool, adjustment SelectionAdjustment) bool {
		newPosInContainer := m.convertToContainerCoordinates(layoutCoordinates, newPosition)
		prevPosInContainer := m.convertToContainerCoordinates(layoutCoordinates, previousPosition)

		m.mu.Lock()
		m.isInTouchMode = isInTouchMode
		m.mu.Unlock()

		return m.updateSelection(newPosInContainer, prevPosInContainer, isStartHandle, adjustment)
	}

	registrar.OnSelectionUpdateEndCallback = func() {
		m.setShowToolbar(true)
		m.mu.Lock()
		m.draggingHandle = HandleNone
		m.currentDragPosition = nil
		m.isLongPressOrClickSelection = false
		m.mu.Unlock()
	}

	registrar.OnSelectableChangeCallback = func(selectableId int64) {
		if _, exists := registrar.Subselections()[selectableId]; exists {
			m.OnRelease()
			m.SetSelection(nil)
		}
	}

	registrar.AfterSelectableUnsubscribe = func(selectableId int64) {
		m.mu.Lock()
		selection := m.selection
		m.mu.Unlock()

		if selection != nil {
			if selectableId == selection.Start.SelectableId {
				m.mu.Lock()
				m.startHandlePosition = nil
				m.mu.Unlock()
			}
			if selectableId == selection.End.SelectableId {
				m.mu.Lock()
				m.endHandlePosition = nil
				m.mu.Unlock()
			}
		}

		if _, exists := registrar.Subselections()[selectableId]; exists {
			m.updateSelectionToolbar()
		}
	}

	return m
}

// Selection returns the current selection.
func (m *SelectionManager) Selection() *Selection {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.selection
}

// SetSelection sets the current selection.
func (m *SelectionManager) SetSelection(selection *Selection) {
	m.mu.Lock()
	m.selection = selection
	m.mu.Unlock()

	if selection != nil {
		m.updateHandleOffsets()
	}
}

// IsInTouchMode returns whether touch mode is active.
func (m *SelectionManager) IsInTouchMode() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.isInTouchMode
}

// SetIsInTouchMode sets the touch mode state.
func (m *SelectionManager) SetIsInTouchMode(value bool) {
	m.mu.Lock()
	changed := m.isInTouchMode != value
	m.isInTouchMode = value
	m.mu.Unlock()

	if changed {
		m.updateSelectionToolbar()
	}
}

// HasFocus returns whether the container has focus.
func (m *SelectionManager) HasFocus() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.hasFocus
}

// SetHasFocus sets the focus state.
func (m *SelectionManager) SetHasFocus(value bool) {
	m.mu.Lock()
	m.hasFocus = value
	m.mu.Unlock()
}

// Modifier returns the modifier for the selection container.
func (m *SelectionManager) Modifier() modifier.Modifier {
	// TODO: Build full modifier chain including:
	// - onGloballyPositioned
	// - focusRequester
	// - onFocusChanged
	// - focusable
	// - updateSelectionTouchMode
	// - onKeyEvent
	// - selectionMagnifier
	// - contextMenuComponents
	return modifier.EmptyModifier
}

// StartHandlePosition returns the start handle position.
func (m *SelectionManager) StartHandlePosition() *geometry.Offset {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.startHandlePosition
}

// EndHandlePosition returns the end handle position.
func (m *SelectionManager) EndHandlePosition() *geometry.Offset {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.endHandlePosition
}

// StartHandleLineHeight returns the line height at the start handle.
func (m *SelectionManager) StartHandleLineHeight() float32 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.startHandleLineHeight
}

// EndHandleLineHeight returns the line height at the end handle.
func (m *SelectionManager) EndHandleLineHeight() float32 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.endHandleLineHeight
}

// DraggingHandle returns the currently dragging handle.
func (m *SelectionManager) DraggingHandle() Handle {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.draggingHandle
}

// IsTriviallyCollapsedSelection returns whether start and end anchors are equal.
func (m *SelectionManager) IsTriviallyCollapsedSelection() bool {
	m.mu.RLock()
	selection := m.selection
	m.mu.RUnlock()

	if selection == nil {
		return true
	}
	return selection.Start == selection.End
}

// IsNonEmptySelection returns whether the selection selects any characters.
func (m *SelectionManager) IsNonEmptySelection() bool {
	m.mu.RLock()
	selection := m.selection
	m.mu.RUnlock()

	if selection == nil {
		return false
	}
	if selection.Start == selection.End {
		return false
	}

	if selection.Start.SelectableId == selection.End.SelectableId {
		return true
	}

	// Check subselections for non-empty ranges
	for _, subSelection := range m.selectionRegistrar.Subselections() {
		if subSelection.Start.Offset != subSelection.End.Offset {
			return true
		}
	}
	return false
}

// OnRelease clears the selection.
func (m *SelectionManager) OnRelease() {
	m.selectionRegistrar.SetSubselections(make(map[int64]*Selection))
	m.setShowToolbar(false)

	m.mu.RLock()
	selection := m.selection
	isInTouchMode := m.isInTouchMode
	m.mu.RUnlock()

	if selection != nil {
		if m.onSelectionChange != nil {
			m.onSelectionChange(nil)
		}
		if isInTouchMode && m.HapticFeedback != nil {
			m.HapticFeedback.PerformHapticFeedback(HapticFeedbackTextHandleMove)
		}
	}
}

// GetSelectedText returns the selected text as an AnnotatedString.
func (m *SelectionManager) GetSelectedText() *text.AnnotatedString {
	m.mu.RLock()
	selection := m.selection
	m.mu.RUnlock()

	if selection == nil {
		return nil
	}

	subselections := m.selectionRegistrar.Subselections()
	if len(subselections) == 0 {
		return nil
	}

	// TODO: Implement full text extraction across selectables
	return nil
}

// Copy copies the selected text.
func (m *SelectionManager) Copy() {
	selectedText := m.GetSelectedText()
	if selectedText != nil && len(selectedText.Text()) > 0 {
		if m.OnCopyHandler != nil {
			m.OnCopyHandler(*selectedText)
		}
	}
}

// SelectAll selects all text in the container.
func (m *SelectionManager) SelectAll() {
	// TODO: Implement selectAll logic
}

// SetOnSelectionChange sets the selection change callback.
func (m *SelectionManager) SetOnSelectionChange(callback func(*Selection)) {
	m.mu.Lock()
	m.onSelectionChange = func(newSelection *Selection) {
		m.selection = newSelection
		if callback != nil {
			callback(newSelection)
		}
	}
	m.mu.Unlock()
}

// SetContainerLayoutCoordinates sets the container's layout coordinates.
func (m *SelectionManager) SetContainerLayoutCoordinates(coords layout.LayoutCoordinates) {
	m.mu.Lock()
	m.containerLayoutCoordinates = coords
	m.mu.Unlock()

	m.mu.RLock()
	hasFocus := m.hasFocus
	selection := m.selection
	m.mu.RUnlock()

	if hasFocus && selection != nil {
		positionInWindow := coords.PositionInWindow()
		m.mu.Lock()
		previousPosition := m.previousPosition
		if previousPosition == nil || *previousPosition != positionInWindow {
			m.previousPosition = &positionInWindow
			m.mu.Unlock()
			m.updateHandleOffsets()
			m.updateSelectionToolbar()
		} else {
			m.mu.Unlock()
		}
	}
}

// Private helper methods

func (m *SelectionManager) setShowToolbar(show bool) {
	m.mu.Lock()
	m.showToolbar = show
	m.mu.Unlock()
	m.updateSelectionToolbar()
}

func (m *SelectionManager) updateHandleOffsets() {
	m.mu.RLock()
	selection := m.selection
	containerCoords := m.containerLayoutCoordinates
	m.mu.RUnlock()

	if selection == nil || containerCoords == nil || !containerCoords.IsAttached() {
		m.mu.Lock()
		m.startHandlePosition = nil
		m.endHandlePosition = nil
		m.mu.Unlock()
		return
	}

	// TODO: Implement full handle position calculation
}

func (m *SelectionManager) updateSelectionToolbar() {
	m.mu.RLock()
	hasFocus := m.hasFocus
	showToolbar := m.showToolbar
	isInTouchMode := m.isInTouchMode
	m.mu.RUnlock()

	if !hasFocus {
		return
	}

	if showToolbar && isInTouchMode {
		// TODO: Show toolbar
	} else if m.TextToolbar != nil && m.TextToolbar.Status() == TextToolbarStatusShown {
		m.TextToolbar.Hide()
	}
}

func (m *SelectionManager) convertToContainerCoordinates(layoutCoordinates layout.LayoutCoordinates, offset geometry.Offset) geometry.Offset {
	m.mu.RLock()
	containerCoords := m.containerLayoutCoordinates
	m.mu.RUnlock()

	if containerCoords == nil || !containerCoords.IsAttached() {
		return geometry.OffsetUnspecified
	}
	return containerCoords.LocalPositionOf(layoutCoordinates, offset)
}

func (m *SelectionManager) startSelection(position geometry.Offset, isStartHandle bool, adjustment SelectionAdjustment) {
	m.mu.Lock()
	m.previousSelectionLayout = nil
	m.mu.Unlock()

	m.updateSelectionInternal(position, geometry.OffsetUnspecified, isStartHandle, adjustment)
}

func (m *SelectionManager) updateSelection(newPosition, previousPosition geometry.Offset, isStartHandle bool, adjustment SelectionAdjustment) bool {
	if !newPosition.IsOffset() {
		return false
	}
	return m.updateSelectionInternal(newPosition, previousPosition, isStartHandle, adjustment)
}

func (m *SelectionManager) updateSelectionInternal(position, previousPosition geometry.Offset, isStartHandle bool, adjustment SelectionAdjustment) bool {
	m.mu.Lock()
	if isStartHandle {
		m.draggingHandle = HandleSelectionStart
	} else {
		m.draggingHandle = HandleSelectionEnd
	}
	m.currentDragPosition = &position
	m.mu.Unlock()

	// TODO: Implement full selection update logic with SelectionLayout
	return false
}

func (m *SelectionManager) selectAllInSelectable(selectableId int64) {
	// TODO: Implement selectAll for a specific selectable
}

// getAnchorSelectable returns the Selectable for the given anchor.
func (m *SelectionManager) getAnchorSelectable(anchor AnchorInfo) Selectable {
	selectableMap := m.selectionRegistrar.SelectableMap()
	if selectable, exists := selectableMap[anchor.SelectableId]; exists {
		return selectable
	}
	return nil
}
