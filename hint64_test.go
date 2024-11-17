package hint64

import (
	"strings"
	"testing"

	"github.com/ebi-yade/gotest/cases"
)

func Test_Parse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		want     int64
		errCheck cases.ErrorCheck
	}{
		{
			name:     "simple positive number",
			input:    "12345",
			want:     12345,
			errCheck: cases.NoError,
		},
		{
			name:     "negative number",
			input:    "-12345",
			want:     -12345,
			errCheck: cases.NoError,
		},
		{
			name:     "explicit positive number",
			input:    "+12345",
			want:     12345,
			errCheck: cases.NoError,
		},
		{
			name:  "decimal with kilo resulting in fraction",
			input: "1.2345k",
			want:  0,
			errCheck: func(t *testing.T, err error) {
				if err == nil {
					t.Error("expected error for decimal resulting in fraction")
				}
				if parseErr, ok := err.(*ParseError); ok {
					if !strings.Contains(parseErr.Message, "too many decimal places") {
						t.Errorf("expected error about decimal places, got: %v", parseErr.Message)
					}
				}
			},
		},
		{
			name:     "decimal with mega valid integer result",
			input:    "1.2345M",
			want:     1234500,
			errCheck: cases.NoError,
		},
		{
			name:  "decimal with mega resulting in fraction",
			input: "1.23456789M",
			want:  0,
			errCheck: func(t *testing.T, err error) {
				if err == nil {
					t.Error("expected error for decimal resulting in fraction")
				}
			},
		},
		{
			name:  "with decimal point but no decimals",
			input: "11.k",
			want:  0,
			errCheck: func(t *testing.T, err error) {
				if err == nil {
					t.Error("expected error for decimal point without decimals")
				}
			},
		},
		{
			name:     "with proper underscore grouping",
			input:    "11_540",
			want:     11540,
			errCheck: cases.NoError,
		},
		{
			name:     "with proper mega value",
			input:    "-2_758_000",
			want:     -2758000,
			errCheck: cases.NoError,
		},
		{
			name:  "with invalid underscore position",
			input: "115_40",
			want:  0,
			errCheck: func(t *testing.T, err error) {
				if err == nil {
					t.Error("expected error for invalid underscore position")
				}
				if parseErr, ok := err.(*ParseError); ok {
					if parseErr.Pos != 3 {
						t.Errorf("expected error at position 3, got %d", parseErr.Pos)
					}
				}
			},
		},
		{
			name:  "with trailing underscore",
			input: "11540_",
			want:  0,
			errCheck: func(t *testing.T, err error) {
				if err == nil {
					t.Error("expected error for trailing underscore")
				}
			},
		},
		{
			name:  "empty string",
			input: "",
			want:  0,
			errCheck: func(t *testing.T, err error) {
				if err == nil {
					t.Error("expected error for empty string")
				}
			},
		},
		{
			name:  "only sign",
			input: "+",
			want:  0,
			errCheck: func(t *testing.T, err error) {
				if err == nil {
					t.Error("expected error for only sign")
				}
			},
		},
		{
			name:     "very large number with T suffix",
			input:    "999T",
			want:     999000000000000,
			errCheck: cases.NoError,
		},
		{
			name:     "alternative B(illion) suffix",
			input:    "1B",
			want:     1000000000,
			errCheck: cases.NoError,
		},
		{
			name:     "lowercase k suffix",
			input:    "123k",
			want:     123000,
			errCheck: cases.NoError,
		},
		{
			name:     "uppercase K suffix",
			input:    "123K",
			want:     123000,
			errCheck: cases.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			tt.errCheck(t, err)
			if err == nil && got != tt.want {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
