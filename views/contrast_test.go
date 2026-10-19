package views

import (
	"context"
	"regexp"
	"strings"
	"testing"

	"github.com/a-h/templ"

	"github.com/khinshankhan/listo/web"
)

// foreground-subtle is 2.5:1 on the ground, too faint for small text, so it
// may only colour hidden decoration.
func TestSubtleIsOnlyDecoration(t *testing.T) {
	entry := Entry{
		Title: "Steins;Gate", Medium: "anime", MediumHref: "/entries?medium=anime", Year: "2011",
		Eyebrow: []string{"anime", "spring 2011"},
		Titles:  []Title{{Kind: "native", Title: "シュタインズ・ゲート", Lang: "ja"}},
		Tags:    []string{"sci-fi"},
		Facts:   []Fact{{Label: "episodes", Value: "24"}},
		Groups:  []GroupRef{{Name: "Main timeline", Href: "/groups/3/x", Position: "#1"}, {Name: "Steins;Gate", Href: "/groups/2/x"}},
		Relations: []Related{{Reading: "prequel of", Title: "Steins;Gate 0", Href: "/entries/11/x", Meta: "anime · 2018",
			Anchor: "here at Ep. 23", Note: "branches off episode 23"}},
		Links: []SiteLink{{Site: "AniList", Href: "https://anilist.co/anime/9253", ID: "9253"}},
	}
	chrome := Chrome{Sidebar: &Sidebar{
		Media:      []MediumLink{{Href: "/entries", Count: 2, Active: true}, {Label: "anime", Href: "/entries?medium=anime", Count: 2}},
		Franchises: []GroupCard{{Name: "Steins;Gate", Href: "/groups/2/x", Entries: 8}},
	}}

	var b strings.Builder
	page := Layout(Head{Title: "t"}, web.Tags{}, chrome)
	if err := page.Render(templ.WithChildren(context.Background(), EntryPage(entry)), &b); err != nil {
		t.Fatal(err)
	}

	tag := regexp.MustCompile(`<[a-z]+[^>]*foreground-subtle[^>]*>`)
	for _, el := range tag.FindAllString(b.String(), -1) {
		// an svg's strokes inherit its aria-hidden
		if strings.Contains(el, "<circle") || strings.Contains(el, "<path") {
			continue
		}
		if !strings.Contains(el, `aria-hidden="true"`) {
			t.Errorf("readable element in foreground-subtle: %s", el)
		}
	}
}
