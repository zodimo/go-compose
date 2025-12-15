package theme_test

import (
	"testing"

	"github.com/zodimo/go-compose/theme"
)

func TestThemeUsage(t *testing.T) {
	// Streamlined way to access a color descriptor
	primary := theme.Primary

	// Applying an update
	primaryWithOpacity := primary.SetOpacity(0.5)

	if !primaryWithOpacity.Compare(primaryWithOpacity) {
		t.Error("Compare should return true for same descriptor")
	}

	if primaryWithOpacity.Compare(primary) {
		t.Error("Expected primaryWithOpacity to be different from primary")
	}

	// Verify updates are preserved
	// Since we can't inspect private fields of ThemeColorDescriptor from this test package easily,
	// we rely on Compare behavior.

	// Ensure other globals are accessible
	_ = theme.OnPrimary
	_ = theme.Background
	_ = theme.Surface
}
