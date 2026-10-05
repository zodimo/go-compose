package main

import (
	"github.com/zodimo/go-compose/pkg/api"
	"github.com/zodimo/go-compose/state"
)

// submittedState holds the most recently submitted form value for display. It is
// UI-only state (never part of the form tree) remembered across recompositions
// so the submission summary survives redraws.
type submittedState struct {
	value map[string]any
	has   bool
}

// rememberSubmitted remembers a submittedState for the demo.
func rememberSubmitted(c api.Composer) *submittedState {
	return state.MustRemember[*submittedState](c, "form-demo-submitted", func() *submittedState {
		return &submittedState{}
	}).Get()
}

// Set records a successful submission.
func (s *submittedState) Set(value map[string]any) {
	s.value = value
	s.has = true
}

// Clear drops the recorded submission (used on Reset or a failed submit).
func (s *submittedState) Clear() {
	s.value = nil
	s.has = false
}

// Get returns the recorded submission and whether one exists.
func (s *submittedState) Get() (map[string]any, bool) {
	return s.value, s.has
}
