package views

// Browse is the catalog: filters, then either a shelf per medium or, once
// something's filtered, everything that matches.
type Browse struct {
	Filters  []Filter
	Filtered bool
	Shelves  []Shelf
	// Cards are the matches, when Filtered.
	Cards []Card
}

// Filter is one dropdown; its first option is "any".
type Filter struct {
	Name, Label string
	Options     []Option
}

// Option is one choice in a filter. An empty Value is "any".
type Option struct {
	Value, Label string
	Selected     bool
}

// Shelf is one medium's newest entries, and where to see the rest. ID is the
// medium, for the heading's id.
type Shelf struct {
	ID, Title, Href string
	Cards           []Card
}

// Card is an entry with its cover; Meta is what it is and how long, eg
// "anime, 2011, 24 episodes".
type Card struct {
	Title, Href, Meta string
}

// filterClass greys a select still on "any".
func filterClass(f Filter) string {
	for _, o := range f.Options {
		if o.Selected && o.Value != "" {
			return "font-medium text-foreground-strong"
		}
	}
	return "text-foreground-muted"
}
