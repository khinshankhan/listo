package views

import "testing"

func TestResolveLinkKind(t *testing.T) {
	cases := []struct {
		href string
		want LinkKind
	}{
		{"/entries/7/steins-gate-anime", LinkInternal},
		{"/", LinkInternal},
		{"#relations", LinkHash},
		{"mailto:hi@example.com", LinkMailto},
		{"tel:+15555550100", LinkTel},
		{"https://anilist.co/anime/9253", LinkExternal},
		{"http://example.com", LinkExternal},
		// as quicksilver has it: no leading slash is not root-relative
		{"entries/7", LinkExternal},
		// as quicksilver has it, though a browser would leave the site
		{"//example.com/x", LinkInternal},
	}
	for _, tc := range cases {
		t.Run(tc.href, func(t *testing.T) {
			if got := ResolveLinkKind(tc.href); got != tc.want {
				t.Errorf("ResolveLinkKind(%q) = %q, want %q", tc.href, got, tc.want)
			}
		})
	}
}

func TestIsFileLink(t *testing.T) {
	cases := []struct {
		href string
		want bool
	}{
		{"/cv.pdf", true},
		{"/cover.PNG", true},
		{"https://example.com/logo.svg", true},
		{"/readme.txt", false},
		{"/entries/7/steins-gate-anime", false},
		{"https://anilist.co/anime/9253", false},
		// as quicksilver has it: the extension has to end the href
		{"/cv.pdf?download=1", false},
	}
	for _, tc := range cases {
		t.Run(tc.href, func(t *testing.T) {
			if got := IsFileLink(tc.href); got != tc.want {
				t.Errorf("IsFileLink(%q) = %v, want %v", tc.href, got, tc.want)
			}
		})
	}
}
