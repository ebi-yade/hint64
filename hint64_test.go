package hint64

import (
	"testing"

	"github.com/ebi-yade/gotest/cases"
	"github.com/stretchr/testify/require"
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
				pe := &ParseError{}
				require.ErrorAs(t, err, &pe)
				require.Equal(t, "too many decimal places for suffix", pe.Message)
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
				pe := &ParseError{}
				require.ErrorAs(t, err, &pe)
				require.Equal(t, "too many decimal places for suffix", pe.Message)
			},
		},
		{
			name:  "with decimal point but no decimals",
			input: "11.k",
			want:  0,
			errCheck: func(t *testing.T, err error) {
				pe := &ParseError{}
				require.ErrorAs(t, err, &pe)
				require.Equal(t, "no digits after decimal point", pe.Message)
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
				pe := &ParseError{}
				require.ErrorAs(t, err, &pe)
				require.Equal(t, "group after underscore must be exactly 3 digits", pe.Message)
			},
		},
		{
			name:  "with trailing underscore",
			input: "11540_",
			want:  0,
			errCheck: func(t *testing.T, err error) {
				pe := &ParseError{}
				require.ErrorAs(t, err, &pe)
				require.Equal(t, "invalid underscore position", pe.Message)
			},
		},
		{
			name:  "empty string",
			input: "",
			want:  0,
			errCheck: func(t *testing.T, err error) {
				pe := &ParseError{}
				require.ErrorAs(t, err, &pe)
				require.Equal(t, "empty input", pe.Message)
			},
		},
		{
			name:  "only sign",
			input: "+",
			want:  0,
			errCheck: func(t *testing.T, err error) {
				pe := &ParseError{}
				require.ErrorAs(t, err, &pe)
				require.Equal(t, "no digits after sign", pe.Message)
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
