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

// Related is one relation from this node's side. Anchor is where it attaches,
// eg "at Ep. 23".
type Related struct {
	Reading, Title, Href, Meta, Anchor, Note string
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
	Positions             []Position
}

// Position is one stop on the path; more than one choice is a fork. Section
// is set where a new section starts.
type Position struct {
	Number, Section string
	Choices         []Choice
}

func (p Position) PickOne() bool {
	return len(p.Choices) > 1
}

// Choice is one way through a position: its entries, in order.
type Choice struct {
	Steps []Step
}

// Step is one entry in a choice. Necessity is empty when it's required.
type Step struct {
	Title, Href, Units, Necessity, Note string
}

// Link is a jump-out link to another site's page for the same work.
type Link struct {
	Site, Href, ID string
}

// Island props: the JSON each island's mount() validates before rendering.

type WaveProps struct {
	Name string `json:"name"`
}
