package str

import "testing"

func TestRepeat(t *testing.T) {
	cases := []struct {
		name   string
		s      string
		count  int
		expect string
	}{
		{"repeat 3 times", "a", 3, "aaa"},
		{"repeat 0 times", "abc", 0, ""},
		{"repeat negative", "x", -2, ""},
		{"repeat empty string", "", 5, ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Repeat(tc.s, tc.count)
			if got != tc.expect {
				t.Errorf("Repeat(%q, %d) = %q; want %q", tc.s, tc.count, got, tc.expect)
			}
		})
	}
}
