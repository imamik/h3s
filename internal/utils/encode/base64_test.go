package encode

import (
	"encoding/base64"
	"testing"
)

func TestToBase64(t *testing.T) {
	cases := []struct {
		name   string
		input  string
		expect string
	}{
		{"empty string", "", ""},
		{"ascii", "hello", base64.StdEncoding.EncodeToString([]byte("hello"))},
		{"unicode", "你好", base64.StdEncoding.EncodeToString([]byte("你好"))},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := ToBase64(tc.input)
			if got != tc.expect {
				t.Errorf("ToBase64(%q) = %q; want %q", tc.input, got, tc.expect)
			}
		})
	}
}
