package diff

import (
	"reflect"
	"testing"
)

func TestCalculateDiff(t *testing.T) {
	tests := []struct {
		name     string
		original string
		modified string
		want     []DiffSection
	}{
		{
			name:     "identical strings",
			original: "line1\nline2\n",
			modified: "line1\nline2\n",
			want: []DiffSection{
				{Type: ChangeEqual, Original: "line1\nline2\n", Modified: "line1\nline2\n"},
			},
		},
		{
			name:     "one line changed",
			original: "line1\nline2\nline3\n",
			modified: "line1\nline2 modified\nline3\n",
			want: []DiffSection{
				{Type: ChangeEqual, Original: "line1\n", Modified: "line1\n"},
				{Type: ChangeDiff, Original: "line2\n", Modified: "line2 modified\n"},
				{Type: ChangeEqual, Original: "line3\n", Modified: "line3\n"},
			},
		},
		{
			name:     "line added",
			original: "line1\n",
			modified: "line1\nline2\n",
			want: []DiffSection{
				{Type: ChangeEqual, Original: "line1\n", Modified: "line1\n"},
				{Type: ChangeDiff, Original: "", Modified: "line2\n"},
			},
		},
		{
			name:     "line removed",
			original: "line1\nline2\n",
			modified: "line1\n",
			want: []DiffSection{
				{Type: ChangeEqual, Original: "line1\n", Modified: "line1\n"},
				{Type: ChangeDiff, Original: "line2\n", Modified: ""},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalculateDiff(tt.original, tt.modified); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CalculateDiff() = %v, want %v", got, tt.want)
			}
		})
	}
}
