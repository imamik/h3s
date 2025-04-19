// Package str provides string manipulation utilities
package str

import "strings"

// Repeat returns a string consisting of count copies of the input string s
// If count is negative, returns an empty string.
func Repeat(s string, count int) string {
	if count < 0 {
		return ""
	}
	return strings.Repeat(s, count)
}
