package handler

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestClampQuery(t *testing.T) {
	cases := []struct {
		name, in, want string
	}{
		{"short is untouched", "steins gate", "steins gate"},
		{"surrounding space goes", "  conan ", "conan"},
		// 199 ASCII bytes and a 2-byte é: 201 bytes but 200 characters
		{"200 characters fit, however many bytes", strings.Repeat("a", 199) + "é", strings.Repeat("a", 199) + "é"},
		{"cut at a character, never inside one", strings.Repeat("é", 201), strings.Repeat("é", 200)},
		{"cut CJK at a character", strings.Repeat("名", 250), strings.Repeat("名", 200)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := clampQuery(tc.in)
			if got != tc.want {
				t.Errorf("clampQuery = %q (%d bytes), want %q (%d bytes)", got, len(got), tc.want, len(tc.want))
			}
			if !utf8.ValidString(got) {
				t.Errorf("clampQuery returned invalid UTF-8")
			}
		})
	}
}
