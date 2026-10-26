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
// crop keeps the top.
func (c CoverImage) imgFit() string {
	if c.Contain {
		return "object-contain"
	}
	return "object-cover object-[50%_30%]"
}

var (
	// the entry page's rail; the top of the page, so never lazy
	entrySlot = coverSlot{Class: "aspect-[5/7] w-[210px] rounded-md", Sizes: "210px"}
	// browse's grid and shelves, coverCols across
	cardSlot = coverSlot{Class: "aspect-[5/7] w-full rounded-[9px]",
		Sizes: "(min-width: 960px) 200px, (min-width: 640px) 25vw, (min-width: 480px) 33vw, 50vw", Lazy: true}
)
