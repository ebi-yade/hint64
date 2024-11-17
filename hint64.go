// Package hnum provides functionality to parse human-readable number strings
package hint64

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// ParseError represents an error that occurs during parsing.
type ParseError struct {
	Input   string
	Pos     int
	Message string
}

func (e *ParseError) Error() string {
	marker := strings.Repeat(" ", e.Pos) + "^"
	return fmt.Sprintf("%s\n%s\n%s", e.Message, e.Input, marker)
}

// multipliers defines the scale factors for metric suffixes
var multipliers = map[rune]int64{
	'k': 1000,
	'K': 1000,
	'M': 1000000,
	'G': 1000000000,
	'T': 1000000000000,
	'B': 1000000000, // Billion (alternative to G)
}

// Parse parses a human-readable number string into an int64.
func Parse(s string) (int64, error) {
	if len(s) == 0 {
		return 0, &ParseError{
			Input:   s,
			Pos:     0,
			Message: "empty input",
		}
	}

	// Handle sign
	var sign int64 = 1
	pos := 0
	if s[0] == '+' || s[0] == '-' {
		if s[0] == '-' {
			sign = -1
		}
		pos = 1
		if pos >= len(s) {
			return 0, &ParseError{
				Input:   s,
				Pos:     pos,
				Message: "no digits after sign",
			}
		}
	}

	// Extract the potentially suffixed part
	numPart := s[pos:]
	if len(numPart) == 0 {
		return 0, &ParseError{
			Input:   s,
			Pos:     pos,
			Message: "no digits found",
		}
	}

	// Check for and handle metric suffix
	var scale int64 = 1
	if len(numPart) > 1 {
		lastChar := rune(numPart[len(numPart)-1])
		if multiplier, hasSuffix := multipliers[lastChar]; hasSuffix {
			scale = multiplier
			numPart = numPart[:len(numPart)-1]
			if len(numPart) == 0 {
				return 0, &ParseError{
					Input:   s,
					Pos:     len(s) - 1,
					Message: "no digits before suffix",
				}
			}
		}
	}

	// Handle decimal point
	parts := strings.Split(numPart, ".")
	if len(parts) > 2 {
		return 0, &ParseError{
			Input:   s,
			Pos:     strings.Index(numPart, ".") + pos,
			Message: "multiple decimal points",
		}
	}
	if len(parts) == 2 {
		if scale == 1 {
			return 0, &ParseError{
				Input:   s,
				Pos:     strings.Index(numPart, ".") + pos,
				Message: "decimal point only allowed with suffix",
			}
		}
		if len(parts[1]) == 0 {
			return 0, &ParseError{
				Input:   s,
				Pos:     strings.Index(numPart, ".") + pos,
				Message: "no digits after decimal point",
			}
		}
		// Check if decimal places and scale would result in a fractional number
		decimalPlaces := len(parts[1])
		remainingScale := scale
		for i := 0; i < decimalPlaces; i++ {
			remainingScale /= 10
		}
		if remainingScale == 0 {
			return 0, &ParseError{
				Input:   s,
				Pos:     strings.Index(numPart, ".") + 1 + len(parts[1]) + pos,
				Message: "too many decimal places for suffix",
			}
		}
	}

	// Validate underscores and build clean number
	var cleanNum strings.Builder
	cleanNum.Grow(len(parts[0]))
	groups := strings.Split(parts[0], "_")
	for i, group := range groups {
		if len(group) == 0 {
			errorPos := 0
			for j := 0; j < i; j++ {
				errorPos += len(groups[j]) + 1
			}
			return 0, &ParseError{
				Input:   s,
				Pos:     errorPos + pos,
				Message: "invalid underscore position",
			}
		}
		if i > 0 && len(group) != 3 {
			errorPos := 0
			for j := 0; j < i; j++ {
				errorPos += len(groups[j]) + 1
			}
			return 0, &ParseError{
				Input:   s,
				Pos:     errorPos + pos - 1,
				Message: "group after underscore must be exactly 3 digits",
			}
		}
		for _, r := range group {
			if !unicode.IsDigit(r) {
				return 0, &ParseError{
					Input:   s,
					Pos:     strings.Index(parts[0], string(r)) + pos,
					Message: "invalid character in number",
				}
			}
			cleanNum.WriteRune(r)
		}
	}

	// Build clean decimal part if exists
	if len(parts) == 2 {
		for _, r := range parts[1] {
			if !unicode.IsDigit(r) {
				return 0, &ParseError{
					Input:   s,
					Pos:     strings.Index(numPart, string(r)) + pos,
					Message: "invalid character in decimal part",
				}
			}
		}
		cleanNum.WriteString(parts[1])
	}

	// Parse the clean number
	base, err := strconv.ParseInt(cleanNum.String(), 10, 64)
	if err != nil {
		return 0, &ParseError{
			Input:   s,
			Pos:     0,
			Message: "number too large for int64",
		}
	}

	// If we had a decimal point, adjust the scale
	if len(parts) == 2 {
		for i := 0; i < len(parts[1]); i++ {
			scale /= 10
		}
	}

	// Calculate final result
	result := base * scale * sign
	if scale != 0 && result/scale != base*sign {
		return 0, &ParseError{
			Input:   s,
			Pos:     0,
			Message: "number too large for int64",
		}
	}

	return result, nil
}
