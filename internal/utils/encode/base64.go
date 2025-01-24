// Package encode contains the functionality for en- & decoding strings
package encode

import "encoding/base64"

// ToBase64 encodes a string to base64
func ToBase64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}
