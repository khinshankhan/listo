package views

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

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

func TestLinkMarkup(t *testing.T) {
	cases := []struct {
		name     string
		props    LinkProps
		want     []string
		dontWant []string
	}{
		{"internal stays in the tab, no icon", LinkProps{Href: "/entries/7/steins-gate-anime"},
			[]string{`href="/entries/7/steins-gate-anime"`, ">Steins;Gate</a>"},
			[]string{"target=", "rel=", "<svg", "sr-only"}},
		{"external opens a tab and says so", LinkProps{Href: "https://anilist.co/anime/9253"},
			[]string{`target="_blank"`, `rel="noopener noreferrer external"`, "<svg", `<span class="sr-only"> (external link)</span>`},
			nil},
		{"a file says it downloads", LinkProps{Href: "/cv.pdf"},
			[]string{`<span class="sr-only"> (file link)</span>`},
			[]string{"target="}},
		{"the icon can be turned off", LinkProps{Href: "https://anilist.co/anime/9253", NoIcon: true},
			[]string{`target="_blank"`},
			[]string{"<svg", "sr-only"}},
		{"the current page is marked", LinkProps{Href: "/entries", Current: true},
			[]string{`aria-current="page"`, `data-active="true"`},
			nil},
		{"a script href is sanitized away", LinkProps{Href: "javascript:alert(1)"},
			[]string{`href="about:invalid#TemplFailedSanitizationURL"`},
			[]string{"javascript:"}},
		{"the variant's classes come with the caller's", LinkProps{Href: "/", Variant: LinkNone, Monochrome: true, Class: "font-semibold"},
			[]string{"focus-visible:outline-2", "font-semibold"},
			[]string{"hover:text-accent-11", "subtle-underline"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var b strings.Builder
			children := templ.ComponentFunc(func(_ context.Context, w io.Writer) error {
				_, err := io.WriteString(w, "Steins;Gate")
				return err
			})
			if err := Link(tc.props).Render(templ.WithChildren(context.Background(), children), &b); err != nil {
				t.Fatal(err)
			}
			got := b.String()
			for _, w := range tc.want {
				if !strings.Contains(got, w) {
					t.Errorf("missing %q in %s", w, got)
				}
			}
			for _, w := range tc.dontWant {
				if strings.Contains(got, w) {
					t.Errorf("has %q in %s", w, got)
				}
			}
		})
	}
}
