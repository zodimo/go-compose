package input

import (
	"testing"

	"github.com/zodimo/go-compose/compose/ui/next/text"
)

func TestInputTransformation_MaxLength(t *testing.T) {
	transformation := MaxLengthTransformation(5)

	tests := []struct {
		name    string
		input   string
		add     string
		want    string
		reverts bool
	}{
		{"under limit", "Hi", "!", "Hi!", false},
		{"at limit", "Hello", "", "Hello", false},
		{"over limit", "Hello", " World", "Hello", true}, // Should revert
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tfcs := NewTextFieldCharSequence(tt.input, text.NewTextRange(len(tt.input), len(tt.input)))
			buffer := NewTextFieldBuffer(tfcs)
			buffer.Append(tt.add)

			transformation.TransformInput(buffer)

			if buffer.String() != tt.want {
				t.Errorf("expected '%s', got '%s'", tt.want, buffer.String())
			}
		})
	}
}

func TestInputTransformation_AllCaps(t *testing.T) {
	transformation := AllCapsTransformation()

	tfcs := NewTextFieldCharSequence("", text.NewTextRange(0, 0))
	buffer := NewTextFieldBuffer(tfcs)
	buffer.Append("hello world")

	transformation.TransformInput(buffer)

	if buffer.String() != "HELLO WORLD" {
		t.Errorf("expected 'HELLO WORLD', got '%s'", buffer.String())
	}
}

func TestInputTransformation_DigitsOnly(t *testing.T) {
	transformation := DigitsOnlyTransformation()

	tfcs := NewTextFieldCharSequence("", text.NewTextRange(0, 0))
	buffer := NewTextFieldBuffer(tfcs)
	buffer.Append("abc123def456")

	transformation.TransformInput(buffer)

	if buffer.String() != "123456" {
		t.Errorf("expected '123456', got '%s'", buffer.String())
	}
}

func TestInputTransformation_AlphanumericOnly(t *testing.T) {
	transformation := AlphanumericOnlyTransformation()

	tfcs := NewTextFieldCharSequence("", text.NewTextRange(0, 0))
	buffer := NewTextFieldBuffer(tfcs)
	buffer.Append("Hello, World! 123")

	transformation.TransformInput(buffer)

	if buffer.String() != "HelloWorld123" {
		t.Errorf("expected 'HelloWorld123', got '%s'", buffer.String())
	}
}

func TestInputTransformation_Chain(t *testing.T) {
	// Max 10 chars AND uppercase
	transformation := ChainTransformations(
		MaxLengthTransformation(10),
		AllCapsTransformation(),
	)

	tfcs := NewTextFieldCharSequence("", text.NewTextRange(0, 0))
	buffer := NewTextFieldBuffer(tfcs)
	buffer.Append("hello")

	transformation.TransformInput(buffer)

	if buffer.String() != "HELLO" {
		t.Errorf("expected 'HELLO', got '%s'", buffer.String())
	}
}

func TestInputTransformation_ByValue(t *testing.T) {
	// Only allow even-length strings
	transformation := InputTransformationByValue(func(current, proposed string) string {
		if len(proposed)%2 == 0 {
			return proposed
		}
		return current
	})

	tfcs := NewTextFieldCharSequence("Hi", text.NewTextRange(2, 2))
	buffer := NewTextFieldBuffer(tfcs)
	buffer.Append("!") // Would make length 3 (odd)

	transformation.TransformInput(buffer)

	if buffer.String() != "Hi" {
		t.Errorf("expected 'Hi', got '%s'", buffer.String())
	}
}
