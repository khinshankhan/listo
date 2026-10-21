package views

// Franchise is a franchise page: what's in it, how it's split into
// continuities, and the ways in.
type Franchise struct {
	Name string
	// Trail is the breadcrumb, outermost group first, ending at this one.
	Trail []Crumb
	// Stats is the line under the title: "6 entries", "2 media".
	Stats []string
	// StartHere is where the "start here" button goes, when there's an
	// ordering to start from.
	StartHere    string
	Continuities []ContinuityCard
	Sorts        []SortLink
	Sections     []EntrySection
	Orderings    []Ordering
	// Elsewhere names the parent the siblings share, for their heading.
	Elsewhere      string
	ElsewhereAside string
	Siblings       []SiblingCard
}

// Crumb is one step of a breadcrumb. The page itself has no Href.
type Crumb struct {
	Kind, Name, Href string
}

// ContinuityCard is one continuity, with how it stands.
type ContinuityCard struct {
	Name, Href, Standing string
	Entries              int
}

// SortLink switches how the entries are grouped.
type SortLink struct {
	Label, Href string
	Active      bool
}

// EntrySection is one group of entries. A section without a Label is the
// whole list, eg by release.
type EntrySection struct {
	Label, Href, Standing string
	Rows                  []EntryRow
}

// EntryRow is one entry in a franchise list. Role is its part in its group,
// eg "source material"; Detail is when and how long, eg "spring 2011, 24
// episodes".
type EntryRow struct {
	Title, Href, Role, Medium, Detail string
}

// SiblingCard is a group beside this one under the same parent.
type SiblingCard struct {
	Name, Href, Meta, Reading string
}

// Anchor is a link to a section of the current page.
type Anchor struct {
	Label, Href string
}

type standingStyle struct {
	Bar, Badge, Dot, Text string
}

// standingStyles colour a continuity by how it stands, one accent ramp each.
var standingStyles = map[string]standingStyle{
	"primary":   {Bar: "bg-moss-9", Badge: "bg-moss-3", Dot: "bg-moss-9", Text: "text-moss-11"},
	"active":    {Bar: "bg-spring-9", Badge: "bg-spring-3", Dot: "bg-spring-9", Text: "text-spring-11"},
	"concluded": {Bar: "bg-sky-9", Badge: "bg-sky-3", Dot: "bg-sky-9", Text: "text-sky-11"},
	"disputed":  {Bar: "bg-honey-9", Badge: "bg-honey-3", Dot: "bg-honey-9", Text: "text-honey-11"},
	"abandoned": {Bar: "bg-ember-9", Badge: "bg-ember-3", Dot: "bg-ember-9", Text: "text-ember-11"},
}

// standing is a standing's colours; anything else gets a neutral bar.
func standing(s string) standingStyle {
	if style, ok := standingStyles[s]; ok {
		return style
	}
	return standingStyle{Bar: "bg-border"}
}

// standingLegend explains each standing, in the order they rank.
var standingLegend = []struct{ Standing, Meaning string }{
	{"primary", "the one to start with"},
	{"active", "still being made"},
	{"concluded", "finished, intact"},
	{"disputed", "fans disagree"},
	{"abandoned", "dropped mid-story"},
}

func sortClass(active bool) string {
	if active {
		return "flex h-[27px] items-center rounded-full bg-foreground-strong px-[11px] text-12 leading-4 font-semibold text-background-2 " + forcedCurrent
	}
	return "flex h-[27px] items-center rounded-full border border-border bg-background-2 px-[11px] text-12 leading-4 font-medium text-foreground-strong hover:border-accent-7"
}
