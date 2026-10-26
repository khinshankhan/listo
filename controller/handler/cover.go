package handler

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/khinshankhan/listo/media"
	"github.com/khinshankhan/listo/views"
)

// srcWidth is the variant a browser without srcset gets: enough for any
// slot at 1x.
const srcWidth = 320

// coverOf is an entry's cover from its row: the stand-in from its medium,
// title and start, and the stored image over it when the row has one.
func coverOf(store media.Store, medium, title string, start *time.Time, key *string, width, height *int32, fit *string) views.Cover {
	c := views.Cover{Medium: medium, Title: title}
	if start != nil {
		c.Year = start.Format("2006")
	}
	if key == nil || width == nil || height == nil || fit == nil {
		return c
	}

	url := func(w int) string { return store.URL(media.VariantKey(*key, strconv.Itoa(w), "jpg")) }
	widths := media.Widths(media.KindCover, int(*width))
	src := widths[0]
	set := make([]string, 0, len(widths))
	for _, w := range widths {
		set = append(set, fmt.Sprintf("%s %dw", url(w), w))
		if w <= srcWidth {
			src = w
		}
	}
	c.Image = &views.CoverImage{Src: url(src), Srcset: strings.Join(set, ", "),
		Width: int(*width), Height: int(*height), Contain: *fit == "contain"}
	return c
}
