package views

import "github.com/a-h/templ"

// Timeline is a group's works drawn by year: a row per line of the story,
// boxed by the part of the group each sits in, with the links between them.
// Everything is in SVG units, Width wide and Height tall. As JSON it's the
// timeline island's props, without what only the SVG needs.
type Timeline struct {
	Width  int            `json:"width"`
	Height int            `json:"height"`
	Ticks  []Tick         `json:"ticks"`
	Boxes  []TimelineBox  `json:"boxes"`
	Works  []TimelineWork `json:"works"`
	Links  []TimelineLink `json:"links"`
	// Spoilers is whether any link is one, so the page offers to show them.
	Spoilers bool `json:"-"`
	// Undated counts the entries left off for having no date.
	Undated int `json:"-"`
	// Lines say what the drawing shows, for screen readers.
	Lines []TimelineLine `json:"-"`
}

// Tick is a year on the axis.
type Tick struct {
	X     int    `json:"x"`
	Label string `json:"label"`
}

// TimelineBox is the dashed box around one part of the group, eg a
// continuity, with its name in the corner.
type TimelineBox struct {
	X    int    `json:"x"`
	Y    int    `json:"y"`
	W    int    `json:"w"`
	H    int    `json:"h"`
	Name string `json:"name"`
}

// TimelineWork is one work on its row: its name, and a dot per release
// joined by a bar. EndAnchored names end at NameX instead of starting there.
// ID is its first entry's, and Left and Right the stretch of row it takes.
type TimelineWork struct {
	ID          int64         `json:"id"`
	Name        string        `json:"name"`
	Href        string        `json:"href"`
	NameX       int           `json:"-"`
	NameY       int           `json:"-"`
	EndAnchored bool          `json:"-"`
	Left        int           `json:"left"`
	Right       int           `json:"right"`
	Y           int           `json:"y"`
	Dots        []TimelineDot `json:"dots"`
}

// TimelineDot is a release, or a few at once, with what they are, eg "VN,
// anime".
type TimelineDot struct {
	X     int    `json:"x"`
	Label string `json:"label"`
}

// TimelineLink is a line between works: a sequel's arrow, or a dotted line
// for anything else. Path is the line's SVG path and Head the arrowhead's.
// It runs from the Source work's dot to the Target work's, left to right.
type TimelineLink struct {
	Source    int64  `json:"source"`
	SourceDot int    `json:"sourceDot"`
	Target    int64  `json:"target"`
	TargetDot int    `json:"targetDot"`
	Path      string `json:"-"`
	Head      string `json:"-"`
	Label     string `json:"label"`
	LabelX    int    `json:"-"`
	LabelY    int    `json:"-"`
	Sequel    bool   `json:"sequel"`
	Spoiler   bool   `json:"spoiler"`
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
