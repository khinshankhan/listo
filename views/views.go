// Package views is every page, in templ.
package views

import "fmt"

// Head is what a page puts in <head>.
type Head struct {
	Title, Description, Canonical string
	// NoIndex keeps a page out of search engines, eg search results.
	NoIndex bool
}

// Entry is an entry page, in display form. Format is empty when it would
// repeat the medium.
type Entry struct {
	Title, Medium, Format, Year string
	// MediumHref browses the entry's medium, for the breadcrumb.
	MediumHref string
	// Eyebrow is the line over the title: medium, format, when.
	Eyebrow []string
	// Titles are every other title the entry goes by.
	Titles    []Title
	Tags      []string
	Facts     []Fact
	Relations []Related
	Groups    []GroupRef
	Links     []SiteLink
}

// Fact is one row of an entry's facts card: "episodes", "24".
type Fact struct {
	Label, Value string
}

// Title is one of a node's titles. Lang is set when known.
type Title struct {
	Kind, Title, Lang string
}

// Related is one relation from this node's side. Anchor is where it attaches,
// eg "at Ep. 23".
type Related struct {
	Reading, Title, Href, Meta, Anchor, Note string
}

// GroupRef is a group a node sits in, nearest first.
type GroupRef struct {
	Name, Href, Position string
}

// Ordering is a curated path through a group. Names render as typed.
type Ordering struct {
	// Kind is "recommended", "release" or "chronological".
	Kind, Name, Rationale string
	Origin, Href          string
	Positions             []Position
}

// Position is one stop on the path; more than one choice is a fork. Section
// is set where a new section starts.
type Position struct {
	Ordinal         int
	Number, Section string
	Choices         []Choice
}

// Skip is a position the first time through can pass over.
func (p Position) Skip() bool {
	for _, c := range p.Choices {
		for _, s := range c.Steps {
			if !s.Skip {
				return false
			}
		}
	}
	return true
}

func (p Position) PickOne() bool {
	return len(p.Choices) > 1
}

// Choice is one way through a position: its entries, in order.
type Choice struct {
	Steps []Step
}

// Step is one entry in a choice. Necessity is empty when it's required.
// Medium and Detail describe the entry, eg "anime", "spring 2011, 24
// episodes".
type Step struct {
	Title, Href, Units, Necessity, Note string
	Medium, Detail                      string
	Skip                                bool
}

// SiteLink is a jump-out link to another site's page for the same work.
type SiteLink struct {
	Site, Href, ID string
}

// Chrome is the page shell. A nil Sidebar means top bar only.
type Chrome struct {
	// Query fills the top bar's search box, on a results page.
	Query   string
	Sidebar *Sidebar
	// HideSearch leaves search out of the top bar, for a page that has it.
	HideSearch bool
}

func searchIconSize(hero bool) string {
	if hero {
		return "17"
	}
	return "15"
}

// Sidebar is the catalog nav: media, then franchises.
type Sidebar struct {
	Media      []MediumLink
	Franchises []GroupCard
	// OnPage links to the current page's sections.
	OnPage []Anchor
	// Orderings are the current group's orderings.
	Orderings []OrderingLink
}

// OrderingLink is an ordering in the sidebar; Meta is its kind and origin.
type OrderingLink struct {
	Name, Meta, Href string
	Active           bool
}

// GroupCard is a top-level group, with how many entries sit beneath it.
type GroupCard struct {
	Name, Kind, Href string
	Entries          int
	// Active marks the group being shown.
	Active bool
}

// MediumLink filters the browse list to one medium; the empty Label is "all".
type MediumLink struct {
	Label, Href string
	Count       int
	Active      bool
}

// Search is the search page: the query as typed, and what it found.
type Search struct {
	Query    string
	Searched bool
	Results  []Result
	// Media filters the results, when there's more than one kind to filter.
	Media []MediumLink
	// Suggestions are nearby titles, after a search that found nothing.
	Suggestions []Suggestion
	// AniList searches the same query there.
	AniList string
}

// Suggestion is a title to search for instead.
type Suggestion struct {
	Title, Href string
}

// Result is one search hit. Matched is set when an alias matched, and
// MatchedAs says which kind, eg "romanized title". Meta is the one-line
// summary the top bar shows; the search page uses Medium and Detail.
type Result struct {
	Title, Href, Meta, Matched string
	MatchedAs, Medium, Detail  string
}

// Browse is every entry, or one medium's.
type Browse struct {
	Media   []MediumLink
	Entries []Row
}

// Row is one linked line in a list.
type Row struct {
	Title, Href, Meta string
}

// count is "1 entry", "36 entries".
func count(n int, one, many string) string {
	if n == 1 {
		return fmt.Sprintf("1 %s", one)
	}
	return fmt.Sprintf("%d %s", n, many)
}
