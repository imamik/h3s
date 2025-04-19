package str

import (
	"strings"
	"testing"
)

func TestNormalizeLength(t *testing.T) {
	cases := []struct {
		name   string
		s      string
		l      int
		expect string
	}{
		{"shorter than l", "foo", 6, "foo   "},
		{"exact length", "foobar", 6, "foobar"},
		{"longer than l", "foobarbaz", 6, "foo..."},
		{"empty string", "", 4, "    "},
		{"l zero", "foo", 0, "..."},
		{"l negative", "foo", -1, "..."},
		{"l less than 3", "foobar", 2, "..."},
		{"large l", "x", 100, "x" + strings.Repeat(" ", 99)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := NormalizeLength(tc.s, tc.l)
			if len(got) != tc.l && tc.l > 3 {
				t.Errorf("NormalizeLength(%q, %d) length = %d; want %d", tc.s, tc.l, len(got), tc.l)
			}
			if tc.l <= 3 && got != "..." {
				t.Errorf("NormalizeLength(%q, %d) = %q; want '...'", tc.s, tc.l, got)
			}
		})
	}
}
