package views

// Cover is an entry's cover slot: a text stand-in, the medium, title and
// year, under the stored image when there is one. The stand-in is also what
// shows while the image loads, and if it never does.
type Cover struct {
	Medium, Title, Year string
	Image               *CoverImage
}

// CoverImage is a stored cover. Srcset lists every variant; Width and Height
// are the original's, for the aspect.
type CoverImage struct {
	Src, Srcset   string
	Width, Height int
	// Contain letterboxes a cover too far off 5:7 rather than cropping it.
	Contain bool
}

// coverSlot is where a cover sits: the box's size and corners, the width it
// takes for sizes, and whether it's far enough down to load lazily.
type coverSlot struct {
	Class, Sizes string
	Lazy         bool
}

// imgFit is how an image fills its slot. Titles and faces sit high, so a
// crop keeps the top. A letterboxed one paints its own bars, or the
// stand-in's words would show around it.
func (c CoverImage) imgFit() string {
	if c.Contain {
		return "object-contain bg-placeholder"
	}
	return "object-cover object-[50%_30%]"
}

var (
	// the entry page's rail; the top of the page, so never lazy
	entrySlot = coverSlot{Class: "aspect-[5/7] w-[210px] rounded-md", Sizes: "210px"}
	// a list row's thumbnail, in browse and on a franchise
	rowSlot = coverSlot{Class: "h-[46px] w-[33px] rounded-sm", Sizes: "33px", Lazy: true}
	// a search result's
	resultSlot = coverSlot{Class: "h-14 w-10 rounded-sm", Sizes: "40px", Lazy: true}
	// a related work's card, under an entry or a franchise
	relatedSlot = coverSlot{Class: "h-[52px] w-[37px] rounded-sm", Sizes: "37px", Lazy: true}
)

// aboveFold is how many covers of the first row load straight away: the
// widest grid's row, give or take.
const aboveFold = 8

// cardSlot is browse's grid and shelves, coverCols across; only the first
// row loads before it's scrolled to.
func cardSlot(eager bool) coverSlot {
	return coverSlot{Class: "aspect-[5/7] w-full rounded-[9px]",
		Sizes: "(min-width: 960px) 200px, (min-width: 640px) 25vw, (min-width: 480px) 33vw, 50vw", Lazy: !eager}
}

// stepSlot is an ordering step's thumbnail, faded with the step when it's
// one to skip.
func stepSlot(skip bool) coverSlot {
	s := coverSlot{Class: "h-[58px] w-[41px] rounded-sm", Sizes: "41px", Lazy: true}
	if skip {
		s.Class += " opacity-60"
	}
	return s
}
