// Package views is every page, in templ.
package views

// Head is what a page puts in <head>.
type Head struct {
	Title, Description, Canonical string
}

// Entry is an entry page, in display form. Format is empty when it would
// repeat the medium.
type Entry struct {
	Title, Medium, Format, Country, Year, Released, Units string
	// Native is the first native title, shown under the heading.
	Native, NativeLang string
	Titles             []Title
	Relations          []Related
	Groups             []GroupRef
	Links              []Link
}

type Title struct {
	Kind, Title string
}

// Related is one relation from this node's side.
type Related struct {
	Reading, Title, Href, Meta, Note string
}

// GroupRef is a group a node sits in, nearest first.
type GroupRef struct {
	Name, Href, Position string
}

// Group is what a group page shows, already in display form.
type Group struct {
	Name, Kind, Standing string
	Native, NativeLang   string
	Members              []Member
	Orderings            []Ordering
	Titles               []Title
	Relations            []Related
	Groups               []GroupRef
	Links                []Link
}

// Positioned reports whether any member has a position.
func (g Group) Positioned() bool {
	for _, m := range g.Members {
		if m.Position != "" {
			return true
		}
	}
	return false
}

// Member is one direct member of a group, an entry or a subgroup.
type Member struct {
	Position, Title, Href, Meta string
}

// Ordering is a curated path through a group. Names render as typed.
type Ordering struct {
	Kind, Name, Rationale string
	Steps                 []Step
}

// Step is one row of an ordering, numbered on a position's first branch.
type Step struct {
	Number, Title, Href, Units, Necessity, Note string
}

// Link is a jump-out link to another site's page for the same work.
type Link struct {
	Site, Href, ID string
}

// Island props: the JSON each island's mount() validates before rendering.

type WaveProps struct {
	Name string `json:"name"`
}
