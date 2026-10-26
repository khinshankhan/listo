package views

// Browse is the catalog: filters, then either a shelf per medium or, once
// something's filtered, everything that matches.
type Browse struct {
	Filters  []Filter
	Filtered bool
	Shelves  []Shelf
	// Cards are a page of the matches, when Filtered, out of Total.
	Cards []Card
	Total int
	// Page is which of Pages this is; Prev and Next link either side.
	Page, Pages int
	Prev, Next  string
	// View is "cards", "list" or "text"; Views switch between them.
	// KeepView is the view from the address, for the form to send on.
	View, KeepView string
	Views          []SortLink
	// Formats are the format filter's options for each type, and for none,
	// so picking a type can fill it in before the form is applied.
	Formats []FormatOptions
}

// FormatOptions is the format filter's options when Medium is the type.
type FormatOptions struct {
	Medium  string
	Options []Option
}

// Filter is one dropdown; its first option is "any". Disabled is a filter
// with nothing to pick yet.
type Filter struct {
	Name, Label string
	Options     []Option
	Disabled    bool
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

// listShelf is how much of a shelf the list views show; a list doesn't
// scroll sideways, so it stops sooner than the covers do.
const listShelf = 6

// Card is an entry with its cover: its format, which FormatHref browses, and
// Meta, when and how long, eg "2011 · 24 episodes".
type Card struct {
	Title, Href, Format, FormatHref, Meta string
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
