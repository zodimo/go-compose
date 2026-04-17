package state

import (
	"testing"
)

func TestMustRememberWithNilInterface(t *testing.T) {
	mv := NewMutableValue(nil, nil, nil)

	// if T is an interface we can support nil
	mvTyped, err := MutableValueToTyped[error](mv)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// and it should return nil
	if val := mvTyped.Get(); val != nil {
		t.Errorf("expected <nil>, got %v", val)
	}
}
