// Package web embeds the Vite bundle; `make web` builds it.
package web

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
)

// "all:" because the manifest lives in a dot directory (.vite/), which embed
// skips by default.
//
//go:embed all:dist
var dist embed.FS

// Vite's entries, as named in vite.config.ts.
const (
	entryCSS    = "src/app.css"
	entryLoader = "src/islands.ts"
)

// Tags are the hashed files every page links, found through the manifest.
type Tags struct {
	CSS     []string
	Script  string
	Preload []string
}

// Bundle is a loaded Vite build.
type Bundle struct {
	// Files is the build output, served under /assets/.
	Files fs.FS
	Tags  Tags
}

// chunk is one entry of Vite's .vite/manifest.json.
type chunk struct {
	File    string   `json:"file"`
	CSS     []string `json:"css"`
	Imports []string `json:"imports"`
}

// Load works out the tags every page carries from the manifest.
func Load() (*Bundle, error) {
	files, err := fs.Sub(dist, "dist/bundle")
	if err != nil {
		return nil, err
	}

	raw, err := fs.ReadFile(files, ".vite/manifest.json")
	if err != nil {
		return nil, fmt.Errorf("reading manifest (run make web): %w", err)
	}

	var manifest map[string]chunk
	if err := json.Unmarshal(raw, &manifest); err != nil {
		return nil, fmt.Errorf("parsing manifest: %w", err)
	}

	style, ok := manifest[entryCSS]
	if !ok {
		return nil, fmt.Errorf("manifest has no %s", entryCSS)
	}
	loader, ok := manifest[entryLoader]
	if !ok {
		return nil, fmt.Errorf("manifest has no %s", entryLoader)
	}

	tags := Tags{CSS: []string{"/" + style.File}, Script: "/" + loader.File}
	for _, css := range loader.CSS {
		tags.CSS = append(tags.CSS, "/"+css)
	}
	for _, key := range loader.Imports {
		tags.Preload = append(tags.Preload, "/"+manifest[key].File)
	}

	return &Bundle{Files: files, Tags: tags}, nil
}
