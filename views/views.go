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

// Related is one relation from this entry's side. Href is empty for groups,
// which have no page yet.
type Related struct {
	Reading, Title, Href, Meta, Note string
}

// GroupRef is a group the entry sits in, nearest first.
type GroupRef struct {
	Name, Position string
}

// Link is a jump-out link to another site's page for the same work.
type Link struct {
	Site, Href, ID string
}

// Island props: the JSON each island's mount() validates before rendering.

type WaveProps struct {
	Name string `json:"name"`
}
