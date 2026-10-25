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
	// Spoilers is whether any link is one, so the page offers to show them;
	// Revealed, whether the address asked to.
	Spoilers bool `json:"-"`
	Revealed bool `json:"-"`
	// Undated counts the entries left off for having no date.
	Undated int `json:"-"`
	// Lines say what the drawing shows, for screen readers.
	Lines []TimelineLine `json:"-"`
	// Outline is the drawing for a screen too narrow for it.
	Outline []TimelinePart `json:"-"`
}

// TimelinePart is one box's works, or the loose ones when Name is "", row by
// row.
type TimelinePart struct {
	Name  string
	Items []TimelineItem
}

// TimelineItem is a work in the outline, Depth rows in from the line it
// hangs off, which Via names, eg "side story of Devil May Cry". Spoiler is
// whether saying so is one.
type TimelineItem struct {
	Name    string
	Href    string
	Formats string
	Years   string
	Via     string
	Depth   int
	Spoiler bool
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

// TimelineWork is one work on its row: a card with its name over a dot per
// release, joined by a bar, the dots on the row at Y. ID is its first
// entry's.
type TimelineWork struct {
	ID    int64         `json:"id"`
	Name  string        `json:"name"`
	Href  string        `json:"href"`
	Card  TimelineCard  `json:"card"`
	NameX int           `json:"-"`
	NameY int           `json:"-"`
	Y     int           `json:"y"`
	Dots  []TimelineDot `json:"dots"`
}

// TimelineCard is the box a work is drawn in.
type TimelineCard struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

// TimelineDot is a release, or a few at once, with what they are, eg "VN,
// anime", centred on LabelX. Dots too close to label one by one share a
// label, eg "game ×9", on the first; the rest have none.
type TimelineDot struct {
	X      int    `json:"x"`
	Label  string `json:"label"`
	LabelX int    `json:"labelX"`
}

// TimelineLink is a step line from the Source work to the Target, left to
// right: an arrow for a sequel, side story or spin-off, dotted for anything
// else. Path is its SVG path. Its label is a pill LabelW wide centred on
// LabelX, LabelY, or ending at LabelX when LabelEnd.
type TimelineLink struct {
	Source   int64  `json:"source"`
	Target   int64  `json:"target"`
	Path     string `json:"path"`
	Label    string `json:"label"`
	LabelX   int    `json:"labelX"`
	LabelY   int    `json:"labelY"`
	LabelW   int    `json:"labelW"`
	LabelEnd bool   `json:"labelEnd"`
	Sequel   bool   `json:"sequel"`
	Spoiler  bool   `json:"spoiler"`
}

// LabelLeft is where a link's label pill starts.
func (l TimelineLink) LabelLeft() int {
	if l.LabelEnd {
		return l.LabelX - l.LabelW
	}
	return l.LabelX - l.LabelW/2
}

// TimelineLine is one sentence of the drawing in words.
type TimelineLine struct {
	Text    string
	Spoiler bool
}

// indent sets an outline item in under the row it hangs off, a guide line
// down its left; past a few deep they line up, so a name keeps room to wrap.
func indent(depth int) string {
	switch depth {
	case 0:
		return ""
	case 1:
		return "border-l border-border pl-3"
	case 2:
		return "ml-3 border-l border-border pl-3"
	default:
		return "ml-6 border-l border-border pl-3"
	}
}

// halo outlines a timeline's text in the ground colour, under the fill.
var halo = templ.Attributes{"stroke": "var(--color-background-2)", "stroke-width": "4", "stroke-linejoin": "round", "paint-order": "stroke"}
