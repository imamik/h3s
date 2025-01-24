// Package str provides string manipulation utilities
package str

import "strings"

// Repeat returns a string consisting of count copies of the input string s
func Repeat(s string, count int) string {
	return strings.Repeat(s, count)
}
