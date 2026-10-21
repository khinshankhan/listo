package views

// Home is the front page: what listo is, search, and previews of the guides.
type Home struct {
	Media     []MediumChip
	Franchise *FranchiseGuide
	Ordering  *OrderingGuide
}

// MediumChip is a medium under the hero search. One without an Href has no
// entries yet.
type MediumChip struct {
	Label, Href string
}

// FranchiseGuide previews a franchise page: why it's hard, and its first
// continuities.
type FranchiseGuide struct {
	Name, Href, Blurb string
	Stats             []string
	Continuities      []GuideColumn
}

// GuideColumn is one continuity's first few entries.
type GuideColumn struct {
	Name, Standing string
	Rows           []GuideRow
}

// GuideRow is an entry in a preview; Meta is its medium and year.
type GuideRow struct {
	Title, Href, Meta string
}

// OrderingGuide previews an ordering's first steps. Total is how many steps
// it has in all.
type OrderingGuide struct {
	Name, Href string
	Stats      []string
	Steps      []GuideStep
	Total      int
}

// GuideStep is one step of a preview. Necessity is "required",
// "recommended", "optional" or "skip the first time".
type GuideStep struct {
	Number                       int
	Title, Href, Note, Necessity string
}

// necessityClass colours a necessity badge: required in the accent,
// recommended in sky, the rest plain or dashed.
func necessityClass(n string) string {
	switch n {
	case "required":
		return "bg-moss-3 text-moss-11"
	case "recommended":
		return "bg-sky-3 text-sky-11"
	case "skip the first time":
		return "border border-dashed border-border text-foreground-muted"
	}
	return "border border-border text-foreground-muted"
}

func hasSoon(media []MediumChip) bool {
	for _, m := range media {
		if m.Href == "" {
			return true
		}
	}
	return false
}
