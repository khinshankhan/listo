package views

import (
	"regexp"
	"strings"
)

// Ported from quicksilver's Link (www at 26b6b71) until quicksilver is a
// package. Keep in step with link.tsx, link.variants.ts, headless/link/utils.ts
// and focus.variants.ts.

// LinkKind is what an href points at.
type LinkKind string

const (
	LinkHash     LinkKind = "hash"
	LinkMailto   LinkKind = "mailto"
	LinkTel      LinkKind = "tel"
	LinkInternal LinkKind = "internal"
	LinkExternal LinkKind = "external"
)

// ResolveLinkKind is quicksilver's resolveKindLite.
// NOTE: unlike quicksilver, //host is external.
func ResolveLinkKind(href string) LinkKind {
	switch {
	case strings.HasPrefix(href, "mailto:"):
		return LinkMailto
	case strings.HasPrefix(href, "tel:"):
		return LinkTel
	case strings.HasPrefix(href, "#"):
		return LinkHash
	case strings.HasPrefix(href, "//"):
		return LinkExternal
	case strings.HasPrefix(href, "/"):
		return LinkInternal
	}
	return LinkExternal
}

var extension = regexp.MustCompile(`\.([0-9a-zA-Z]+)$`)

// fileExtensions are what quicksilver treats as downloads.
var fileExtensions = map[string]bool{"pdf": true, "svg": true, "png": true}

// IsFileLink is quicksilver's isFileLite.
func IsFileLink(href string) bool {
	m := extension.FindStringSubmatch(href)
	return m != nil && fileExtensions[strings.ToLower(m[1])]
}

// LinkVariant is how a link looks.
type LinkVariant int

const (
	// LinkDefault is inline text: an accent underline that thickens on hover.
	LinkDefault LinkVariant = iota
	// LinkNav draws its underline in from the right on hover or when active.
	LinkNav
	// LinkNone leaves the look to the caller's classes.
	LinkNone
)

// focusRing is quicksilver's single focus treatment (focus.variants.ts).
const focusRing = "outline-none focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-solid focus-visible:outline-stark-contrast"

// linkClass is quicksilver's linkVariants.
func linkClass(variant LinkVariant, monochrome bool) string {
	classes := []string{
		focusRing,
		// links have no radius of their own; soften the ring without a pill
		"focus-visible:rounded-sm",
		"transition-[color] duration-500",
	}
	switch variant {
	case LinkDefault:
		classes = append(classes, "bg-linear-to-r from-accent-11 to-accent-11 subtle-underline show-underline bg-no-repeat hover:from-accent-8 hover:to-accent-8 hover:drastic-underline")
	case LinkNav:
		classes = append(classes, "bg-linear-to-r from-accent-11 to-accent-11 link-hide bg-right-bottom bg-no-repeat transition-[color,background-size] hover:link-show hover:bg-left-bottom data-[active=true]:link-show data-[active=true]:bg-left-bottom")
	}
	if !monochrome {
		classes = append(classes, "hover:text-accent-11")
	}
	return strings.Join(classes, " ")
}

// LinkProps configures a Link. Class adds to the variant's classes.
type LinkProps struct {
	Href       string
	Variant    LinkVariant
	Monochrome bool
	Class      string
	// NoIcon drops the trailing icon an external or file link would get.
	NoIcon bool
	// Label is an aria-label, for a link whose text alone isn't enough.
	Label string
	// Current sets aria-current and data-active for the current page.
	Current bool
}

// The classes the shared links add to their variant.

const breadcrumbClass = "font-mono text-11 leading-[14px] tracking-wide text-foreground-subtle uppercase"

func sideLinkClass(active bool) string {
	if active {
		return "flex items-center gap-[11px] rounded-lg bg-accent-3 px-2.5 py-[9px]"
	}
	return "flex items-center gap-[11px] rounded-lg px-2.5 py-[9px] hover:bg-background-1"
}

func pillClass(active bool) string {
	if active {
		return "rounded-md border border-accent-7 bg-accent-3 px-3 py-1.5 text-14 text-accent-12"
	}
	return "rounded-md border border-border px-3 py-1.5 text-14 text-foreground hover:border-accent-7"
}
