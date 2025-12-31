package input

import (
	"testing"

	"github.com/zodimo/go-compose/compose/ui/next/text"
)

func TestOutputTransformation_Prefix(t *testing.T) {
	transformation := PrefixTransformation("$")

	tfcs := NewTextFieldCharSequence("100", text.NewTextRange(3, 3))
	buffer := NewTextFieldBuffer(tfcs)

	transformation.TransformOutput(buffer)

	if buffer.String() != "$100" {
		t.Errorf("expected '$100', got '%s'", buffer.String())
	}
}

func TestOutputTransformation_Suffix(t *testing.T) {
	transformation := SuffixTransformation(" USD")

	tfcs := NewTextFieldCharSequence("100", text.NewTextRange(3, 3))
	buffer := NewTextFieldBuffer(tfcs)

	transformation.TransformOutput(buffer)

	if buffer.String() != "100 USD" {
		t.Errorf("expected '100 USD', got '%s'", buffer.String())
	}
}

func TestOutputTransformation_Chain(t *testing.T) {
	transformation := ChainOutputTransformations(
		PrefixTransformation("$"),
		SuffixTransformation(" USD"),
	)

	tfcs := NewTextFieldCharSequence("100", text.NewTextRange(3, 3))
	buffer := NewTextFieldBuffer(tfcs)

	transformation.TransformOutput(buffer)

	if buffer.String() != "$100 USD" {
		t.Errorf("expected '$100 USD', got '%s'", buffer.String())
	}
}

func TestOutputTransformation_PhoneNumber(t *testing.T) {
	transformation := PhoneNumberTransformation()

	tests := []struct {
		input string
		want  string
	}{
		{"123", "(123"},
		{"1234567", "(123) 456-7"},
		{"1234567890", "(123) 456-7890"},
		{"12345678901234", "(123) 456-7890"}, // Truncated
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			tfcs := NewTextFieldCharSequence(tt.input, text.NewTextRange(len(tt.input), len(tt.input)))
			buffer := NewTextFieldBuffer(tfcs)

			transformation.TransformOutput(buffer)

			if buffer.String() != tt.want {
				t.Errorf("expected '%s', got '%s'", tt.want, buffer.String())
			}
		})
	}
}

func TestOutputTransformation_Func(t *testing.T) {
	transformation := OutputTransformationFunc(func(buffer *TextFieldBuffer) {
		buffer.Append("!")
	})

	tfcs := NewTextFieldCharSequence("Hello", text.NewTextRange(5, 5))
	buffer := NewTextFieldBuffer(tfcs)

	transformation.TransformOutput(buffer)

	if buffer.String() != "Hello!" {
		t.Errorf("expected 'Hello!', got '%s'", buffer.String())
	}
}
