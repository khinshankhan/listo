package handler

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/khinshankhan/logstox/fields"

	"github.com/khinshankhan/listo/data/sqlcgen"
	"github.com/khinshankhan/listo/platform/logx"
)

const (
	// the top bar shows a handful; the search page has the rest
	defaultSuggestions = 6
	maxSuggestions     = 10
)

type searchHit struct {
	Title   string `json:"title"`
	Href    string `json:"href"`
	Meta    string `json:"meta,omitempty"`
	Matched string `json:"matched,omitempty"`
}

type searchResponse struct {
	Query   string      `json:"query"`
	Results []searchHit `json:"results"`
	// All is the full search page for the query.
	All string `json:"all"`
}

// Search serves /api/search?q=&limit=, the top bar's suggestions.
func (a *API) Search(w http.ResponseWriter, r *http.Request) {
	q := clampQuery(r.URL.Query().Get("q"))
	limit := defaultSuggestions
	if n, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && n > 0 {
		limit = min(n, maxSuggestions)
	}

	resp := searchResponse{Query: q, Results: []searchHit{}, All: "/search?q=" + url.QueryEscape(q)}
	if q != "" {
		hits, err := a.search.SearchNodes(r.Context(), sqlcgen.SearchNodesParams{Q: q, MaxResults: int32(limit)})
		if err != nil {
			logx.Default().Named("api").Error("searching", fields.Error(err))
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "search failed"})
			return
		}
		for _, res := range resultsFor(hits) {
			resp.Results = append(resp.Results, searchHit{Title: res.Title, Href: res.Href, Meta: res.Meta, Matched: res.Matched})
		}
	}
	writeJSON(w, http.StatusOK, resp)
}
