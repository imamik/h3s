// Package str contains utility functions for strings
package str

import "strings"

// NormalizeLength returns a string with a fixed length
func NormalizeLength(s string, l int) string {
	if l <= 3 {
		return "..."
	}
	if len(s) > l-3 {
		return s[:l-3] + "..."
	}
	return s + strings.Repeat(" ", l-len(s))
}
