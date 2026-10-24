package views

import "github.com/a-h/templ"

// Timeline is a group's works drawn by year: a row per line of the story,
// boxed by the part of the group each sits in, with the links between them.
// Everything is in SVG units, Width wide and Height tall.
type Timeline struct {
	Width, Height int
	Ticks         []Tick
	Boxes         []TimelineBox
	Works         []TimelineWork
	Links         []TimelineLink
	// Spoilers is whether any link is one, so the page offers to show them.
	Spoilers bool
	// Undated counts the entries left off for having no date.
	Undated int
	// Lines say what the drawing shows, for screen readers.
	Lines []TimelineLine
}

// Tick is a year on the axis.
type Tick struct {
	X     int
	Label string
}

// TimelineBox is the dashed box around one part of the group, eg a
// continuity, with its name in the corner.
type TimelineBox struct {
	X, Y, W, H int
	Name       string
}

// TimelineWork is one work on its row: its name, and a dot per release
// joined by a bar. EndAnchored names end at NameX instead of starting there.
type TimelineWork struct {
	Name         string
	NameX, NameY int
	EndAnchored  bool
	Y            int
	Dots         []TimelineDot
}

// TimelineDot is a release, or a few at once, with what they are, eg "VN,
// anime".
type TimelineDot struct {
	X     int
	Label string
}

// TimelineLink is a line between works: a sequel's arrow, or a dotted line
// for anything else. Path is the line's SVG path and Head the arrowhead's.
type TimelineLink struct {
	Path, Head     string
	Label          string
	LabelX, LabelY int
	Sequel         bool
	Spoiler        bool
}

// TimelineLine is one sentence of the drawing in words.
type TimelineLine struct {
	Text    string
	Spoiler bool
}

// halo outlines a timeline's text in the ground colour, under the fill.
var halo = templ.Attributes{"stroke": "var(--color-background-2)", "stroke-width": "4", "stroke-linejoin": "round", "paint-order": "stroke"}

// labelWidth is about how wide a link's label draws, with a little to spare,
// for the backing that keeps lines from running through it.
func labelWidth(s string) int {
	return len(s)*6 + 8
}
