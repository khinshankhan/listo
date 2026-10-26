// Package views is every page, in templ.
package views

import (
	"fmt"

	"github.com/google/uuid"
)

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
	Cover                       Cover
	// MediumHref browses the entry's medium, for the breadcrumb.
	MediumHref string
	// Eyebrow is the line over the title: medium, format, when. The format
	// browses everything in it.
	Eyebrow []EyebrowPart
	// Titles are every other title the entry goes by.
	Titles    []Title
	Tags      []string
	Facts     []Fact
	Relations []Related
	Groups    []GroupRef
	Links     []SiteLink
	// Description is the entry's own summary, when there is one.
	Description string
	// Next is the next unit to air, if one's scheduled; Episodes is the
	// list, in display order.
	Next     *NextUnit
	Episodes []Episode
}

// NextUnit is what airs next and how soon: "Ep. 20", "in 3d 4h".
type NextUnit struct {
	Label, In string
}

// Episode is one unit in an entry's list; When is relative while it's
// recent, eg "2 days ago".
type Episode struct {
	Number, Title, When string
	Upcoming            bool
}

// EyebrowPart is one piece of an entry's eyebrow, a link when Href is set.
type EyebrowPart struct {
	Text, Href string
}

// Fact is one row of an entry's facts card: "episodes", "24". Href, when set,
// makes the value a link.
type Fact struct {
	Label, Value, Href string
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
	ID uuid.UUID
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
// episodes"; Medium is its format, which FormatHref browses.
type Step struct {
	Title, Href, Units, Necessity, Note string
	Medium, FormatHref, Detail          string
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
	// Wide lets a page past the content width, for a grid that uses the room.
	Wide bool
}

func searchIconSize(hero bool) string {
	if hero {
		return "17"
	}
	return "15"
}

// Sidebar is the catalog nav: media, then the biggest franchises. It's the
// same on every page; only which link is current changes.
type Sidebar struct {
	Media      []MediumLink
	Franchises []GroupCard
	// AllFranchises marks the page of every franchise as the current one.
	AllFranchises bool
}

// OrderingLink is one of a group's orderings; Meta is its kind and origin.
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
// Medium is its key, eg "visual_novel".
type MediumLink struct {
	Medium, Label, Href string
	Count               int
	Active              bool
}

// Search is the search page: the query as typed, and what it found.
type Search struct {
	Query    string
	Searched bool
	// Results are the closest of Total hits.
	Results []Result
	Total   int
	// Media filters the results, when there's more than one kind to filter.
	Media []MediumLink
	// Suggestions are nearby titles, after a search that found nothing.
	Suggestions []Suggestion
}

// Suggestion is a title to search for instead.
type Suggestion struct {
	Title, Href string
}

// Result is one search hit. Matched is set when an alias matched, and
// MatchedAs says which kind, eg "romanized title". Meta is the one-line
// summary the top bar shows; the search page uses Medium and Detail.
// FormatHref browses an entry's format.
type Result struct {
	Title, Href, Meta, Matched            string
	MatchedAs, Medium, Detail, FormatHref string
}

// count is "1 entry", "36 entries".
func count(n int, one, many string) string {
	if n == 1 {
		return fmt.Sprintf("1 %s", one)
	}
	return fmt.Sprintf("%d %s", n, many)
}

// shownEpisodes is how many episodes show before the rest fold away.
const shownEpisodes = 5
