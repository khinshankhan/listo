package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/khinshankhan/logstox/fields"

	"github.com/khinshankhan/listo/data/sqlcgen"
	"github.com/khinshankhan/listo/platform/logx"
	"github.com/khinshankhan/listo/views"
)

// Catalog is the slice of the generated queries the catalog pages need.
type Catalog interface {
	NodeByPublicID(ctx context.Context, publicID int64) (sqlcgen.NodeByPublicIDRow, error)
	EntryByPublicID(ctx context.Context, publicID int64) (sqlcgen.EntryByPublicIDRow, error)
	GroupByPublicID(ctx context.Context, publicID int64) (sqlcgen.GroupByPublicIDRow, error)
	GroupOrderings(ctx context.Context, groupID uuid.UUID) ([]sqlcgen.GroupOrderingsRow, error)
	NodeTitles(ctx context.Context, nodeID uuid.UUID) ([]sqlcgen.NodeTitlesRow, error)
	NodeLinks(ctx context.Context, nodeID uuid.UUID) ([]sqlcgen.NodeLinksRow, error)
	NodeRelations(ctx context.Context, nodeID uuid.UUID) ([]sqlcgen.NodeRelationsRow, error)
	NodeGroups(ctx context.Context, memberID uuid.UUID) ([]sqlcgen.NodeGroupsRow, error)
	NodeTags(ctx context.Context, nodeID uuid.UUID) ([]sqlcgen.NodeTagsRow, error)
	GroupEntries(ctx context.Context, groupID uuid.UUID) ([]sqlcgen.GroupEntriesRow, error)
	GroupSiblings(ctx context.Context, memberID uuid.UUID) ([]sqlcgen.GroupSiblingsRow, error)
	SearchNodes(ctx context.Context, arg sqlcgen.SearchNodesParams) ([]sqlcgen.SearchNodesRow, error)
	TitleSuggestions(ctx context.Context, arg sqlcgen.TitleSuggestionsParams) ([]string, error)
	TopGroups(ctx context.Context) ([]sqlcgen.TopGroupsRow, error)
	MediumCounts(ctx context.Context) ([]sqlcgen.MediumCountsRow, error)
	BrowseEntries(ctx context.Context, arg sqlcgen.BrowseEntriesParams) ([]sqlcgen.BrowseEntriesRow, error)
	BrowseShelves(ctx context.Context, perShelf int32) ([]sqlcgen.BrowseShelvesRow, error)
	BrowseGenres(ctx context.Context) ([]string, error)
	BrowseYears(ctx context.Context) ([]int32, error)
}

// EntryPath is an entry's canonical URL: the id is the key, the slug decoration.
func EntryPath(publicID int64, slug string) string {
	return fmt.Sprintf("/entries/%d/%s", publicID, slug)
}

// GroupPath is a group's canonical URL, shaped like EntryPath.
func GroupPath(publicID int64, slug string) string {
	return fmt.Sprintf("/groups/%d/%s", publicID, slug)
}

// OrderingPath is an ordering's URL, under its group's: orderings have no id
// of their own, and their names are unique within a group.
func OrderingPath(group, name string) string {
	return group + "/orderings/" + slugify(name)
}

// slugify is "first-timer order" as "first-timer-order".
func slugify(s string) string {
	var b strings.Builder
	dash := false
	for _, r := range strings.ToLower(s) {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			if dash && b.Len() > 0 {
				b.WriteByte('-')
			}
			b.WriteRune(r)
			dash = false
			continue
		}
		dash = true
	}
	return b.String()
}

// NodePath is the canonical URL of a node of either type.
func NodePath(typ string, publicID int64, slug string) string {
	if typ == "group" {
		return GroupPath(publicID, slug)
	}
	return EntryPath(publicID, slug)
}

// publicID parses the {id} path value.
func publicID(r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	return id, err == nil && id > 0
}

// redirectCanonical redirects a non-canonical slug and reports whether it did.
func redirectCanonical(w http.ResponseWriter, r *http.Request, canonical, slug string) bool {
	if r.PathValue("slug") == slug {
		return false
	}
	moved(w, r, canonical)
	return true
}

// moved 301s to to, keeping the query.
func moved(w http.ResponseWriter, r *http.Request, to string) {
	if r.URL.RawQuery != "" {
		to += "?" + r.URL.RawQuery
	}
	http.Redirect(w, r, to, http.StatusMovedPermanently)
}

// elsewhere redirects an id of the other node type to its page, else 404s.
func (p *Pages) elsewhere(w http.ResponseWriter, r *http.Request, id int64) {
	n, err := p.catalog.NodeByPublicID(r.Context(), id)
	if errors.Is(err, pgx.ErrNoRows) {
		p.NotFound(w, r)
		return
	}
	if err != nil {
		p.serverError(w, r, err)
		return
	}
	moved(w, r, NodePath(n.Type, n.PublicID, n.Slug))
}

func (p *Pages) serverError(w http.ResponseWriter, r *http.Request, err error) {
	logx.Default().Named("pages").Error("serving page",
		fields.String("path", r.URL.Path),
		fields.Error(err))
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

// nodeParts are the sections every node page has, entry or group.
type nodeParts struct {
	native, nativeLang string
	titles             []views.Title
	relations          []views.Related
	groups             []views.GroupRef
	links              []views.SiteLink
}

func (p *Pages) nodeParts(ctx context.Context, id uuid.UUID, country string) (nodeParts, error) {
	var parts nodeParts

	titles, err := p.catalog.NodeTitles(ctx, id)
	if err != nil {
		return parts, fmt.Errorf("titles: %w", err)
	}
	for _, t := range titles {
		title := views.Title{Kind: t.Kind, Title: t.Title}
		if t.Kind == "native" {
			title.Lang = lang(t.Language, country)
			if parts.native == "" {
				parts.native, parts.nativeLang = t.Title, title.Lang
			}
		}
		parts.titles = append(parts.titles, title)
	}

	relations, err := p.catalog.NodeRelations(ctx, id)
	if err != nil {
		return parts, fmt.Errorf("relations: %w", err)
	}
	for _, rel := range relations {
		related := views.Related{
			Reading: rel.Reading,
			Title:   rel.Title,
			Href:    NodePath(rel.Type, rel.PublicID, rel.Slug),
			Meta:    "group",
			Anchor:  anchor(rel.MyUnit, rel.TheirUnit),
			Note:    deref(rel.Note),
		}
		if rel.Type == "entry" {
			related.Meta = withYear(deref(rel.MediumLabel), rel.StartDate)
		}
		parts.relations = append(parts.relations, related)
	}

	groups, err := p.catalog.NodeGroups(ctx, id)
	if err != nil {
		return parts, fmt.Errorf("groups: %w", err)
	}
	for _, g := range groups {
		parts.groups = append(parts.groups, views.GroupRef{
			Name:     g.Name,
			Href:     GroupPath(g.PublicID, g.Slug),
			Position: deref(g.PositionLabel),
		})
	}

	links, err := p.catalog.NodeLinks(ctx, id)
	if err != nil {
		return parts, fmt.Errorf("links: %w", err)
	}
	for _, l := range links {
		parts.links = append(parts.links, views.SiteLink{Site: l.Site, Href: l.Url, ID: l.ExternalID})
	}

	return parts, nil
}

// anchor phrases where a relation attaches, read from this node's side.
func anchor(mine, theirs *string) string {
	switch {
	case mine != nil && theirs != nil:
		return fmt.Sprintf("here at %s, there at %s", *mine, *theirs)
	case mine != nil:
		return "here at " + *mine
	case theirs != nil:
		return "at " + *theirs
	}
	return ""
}
