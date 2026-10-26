package views

// OrderingDetail is one ordering on its own page. Shown is how many positions
// show before the rest fold away.
type OrderingDetail struct {
	Trail             []Crumb
	Name              string
	Stats             []string
	Origin, Rationale string
	Positions         []Position
	Shown             int
	Alternative       *Alternative
	// Orderings are the group's orderings, this one included, when it has
	// more than one.
	Orderings []OrderingLink
}

// Alternative suggests another ordering, for another kind of viewer.
type Alternative struct {
	Title, Body, Label, Href string
}

// orderingLinkClass is one of the switcher's orderings, tinted when it's this
// one.
func orderingLinkClass(active bool) string {
	if active {
		return "flex flex-col gap-0.5 rounded-lg border border-accent-7 bg-accent-3 px-3 py-2 font-semibold text-accent-11 " + forcedCurrent
	}
	return "flex flex-col gap-0.5 rounded-lg border border-border bg-background-2 px-3 py-2 font-medium text-foreground-strong hover:border-accent-7 [&>span+span]:text-foreground-muted"
}
