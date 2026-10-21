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
}

// Alternative suggests another ordering, for another kind of viewer.
type Alternative struct {
	Title, Body, Label, Href string
}

func orderingLinkClass(active bool) string {
	if active {
		return "flex flex-col gap-0.5 rounded-lg bg-accent-3 px-2.5 py-2 font-semibold text-accent-11"
	}
	return "flex flex-col gap-0.5 rounded-lg px-2.5 py-2 font-medium text-foreground-strong hover:bg-background-1 [&>span+span]:text-foreground-muted"
}
