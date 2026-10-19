// Package views is every page, in templ.
package views

// Head is what a page puts in <head>.
type Head struct {
	Title, Description, Canonical string
}

// Entry is an entry page, in display form.
type Entry struct {
	Title, Format, Released, Units, Country string
}

// Island props: the JSON each island's mount() validates before rendering.

type WaveProps struct {
	Name string `json:"name"`
}
