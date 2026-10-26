package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/khinshankhan/listo/data"
	"github.com/khinshankhan/listo/data/sqlcgen"
	"github.com/khinshankhan/listo/media"
	"github.com/khinshankhan/listo/platform/config"
)

// imageCmd is listo image: add one image, or import a manifest of them.
func imageCmd(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return errors.New("listo image: want add or import")
	}
	switch args[0] {
	case "add":
		return imageAdd(ctx, args[1:])
	case "import":
		return imageImport(ctx, args[1:])
	default:
		return fmt.Errorf("listo image: unknown subcommand %q", args[0])
	}
}

// imageSource is one image to add to an entry, and who to credit for it.
type imageSource struct {
	Slug, Kind, From, Source, SourceURL string
}

func imageAdd(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("image add", flag.ContinueOnError)
	kind := fs.String("kind", media.KindCover, "cover, banner, art or screenshot")
	source := fs.String("source", "staff", "where the image came from")
	sourceURL := fs.String("source-url", "", "the page it came from")
	pos, err := parseInterleaved(fs, args)
	if err != nil {
		return err
	}
	if len(pos) != 2 {
		return errors.New("usage: listo image add <entry slug> <file|url> [--kind cover] [--source x --source-url y]")
	}

	pool, store, err := openImages(ctx)
	if err != nil {
		return err
	}
	defer pool.Close()

	src := imageSource{Slug: pos[0], From: pos[1], Kind: *kind, Source: *source, SourceURL: *sourceURL}
	res, err := addFrom(ctx, pool, store, src)
	if err != nil {
		return err
	}
	fmt.Println(res.describe(src))
	return nil
}

func imageImport(ctx context.Context, args []string) error {
	if len(args) != 1 {
		return errors.New("usage: listo image import <manifest.tsv>")
	}
	f, err := os.Open(args[0])
	if err != nil {
		return err
	}
	defer f.Close()
	// files are named relative to raw/ beside the manifest
	sources, err := readManifest(f, filepath.Join(filepath.Dir(args[0]), "raw"))
	if err != nil {
		return err
	}

	pool, store, err := openImages(ctx)
	if err != nil {
		return err
	}
	defer pool.Close()

	var added, had, failed int
	for _, src := range sources {
		res, err := addFrom(ctx, pool, store, src)
		if err != nil {
			failed++
			fmt.Fprintf(os.Stderr, "%s %s: %v\n", src.Slug, src.From, err)
			continue
		}
		if res.New {
			added++
		} else {
			had++
		}
		fmt.Println(res.describe(src))
	}
	fmt.Printf("%d added, %d already there, %d failed\n", added, had, failed)
	if failed > 0 {
		return fmt.Errorf("%d of %d images failed", failed, len(sources))
	}
	return nil
}

func openImages(ctx context.Context) (*pgxpool.Pool, media.Store, error) {
	store, err := media.Open(config.StorageURL())
	if err != nil {
		return nil, nil, err
	}
	pool, err := data.Open(ctx, config.DatabaseURL())
	if err != nil {
		return nil, nil, err
	}
	return pool, store, nil
}

// addFrom reads src and adds it in a transaction of its own.
func addFrom(ctx context.Context, pool *pgxpool.Pool, store media.Store, src imageSource) (added, error) {
	r, err := openSource(ctx, src.From)
	if err != nil {
		return added{}, err
	}
	defer r.Close()

	tx, err := pool.Begin(ctx)
	if err != nil {
		return added{}, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()
	res, err := addImage(ctx, tx, store, src, r)
	if err != nil {
		return added{}, err
	}
	return res, tx.Commit(ctx)
}

// added is an image that's been through addImage. New is false when the
// entry already had it.
type added struct {
	media.Image
	New, Canon bool
}

func (a added) describe(src imageSource) string {
	what := "added"
	if !a.New {
		what = "already had"
	}
	canon := ""
	if a.Canon {
		canon = ", canon"
	}
	return fmt.Sprintf("%s %s %s: %dx%d %s%s", what, src.Slug, a.Kind, a.Width, a.Height, a.Fit, canon)
}

// addImage stores the image in r and records it against the entry, canon if
// it beats the one there. The rows go through tx; the files land in store
// either way, harmlessly, since keys are content hashes.
func addImage(ctx context.Context, tx pgx.Tx, store media.Store, src imageSource, r io.Reader) (added, error) {
	q := sqlcgen.New(tx)
	entry, err := q.EntryIDBySlug(ctx, src.Slug)
	if errors.Is(err, pgx.ErrNoRows) {
		return added{}, fmt.Errorf("no entry %q", src.Slug)
	}
	if err != nil {
		return added{}, err
	}

	img, err := media.Import(ctx, store, r, src.Kind)
	if err != nil {
		return added{}, err
	}
	var sourceURL *string
	if src.SourceURL != "" {
		sourceURL = &src.SourceURL
	}
	n, err := q.AddEntryImage(ctx, sqlcgen.AddEntryImageParams{
		EntryID: entry, Kind: img.Kind, StorageKey: img.Key,
		Width: int32(img.Width), Height: int32(img.Height), Fit: img.Fit,
		Source: src.Source, SourceUrl: sourceURL,
	})
	if err != nil {
		return added{}, err
	}
	if n == 0 {
		return added{Image: img}, nil
	}

	var cur *media.Size
	switch c, err := q.CanonImage(ctx, sqlcgen.CanonImageParams{EntryID: entry, Kind: img.Kind}); {
	case err == nil:
		cur = &media.Size{Width: int(c.Width), Height: int(c.Height)}
	case !errors.Is(err, pgx.ErrNoRows):
		return added{}, err
	}
	if !media.AutoCanon(img.Kind, media.Size{Width: img.Width, Height: img.Height}, cur) {
		return added{Image: img, New: true}, nil
	}
	if err := q.ClearCanonImage(ctx, sqlcgen.ClearCanonImageParams{EntryID: entry, Kind: img.Kind}); err != nil {
		return added{}, err
	}
	if err := q.SetCanonImage(ctx, sqlcgen.SetCanonImageParams{EntryID: entry, StorageKey: img.Key}); err != nil {
		return added{}, err
	}
	return added{Image: img, New: true, Canon: true}, nil
}

// fetchClient gives up on a slow source rather than hanging an import.
var fetchClient = &http.Client{Timeout: 30 * time.Second}

// openSource opens a file, or fetches a URL that serves an image.
func openSource(ctx context.Context, from string) (io.ReadCloser, error) {
	if !strings.HasPrefix(from, "http://") && !strings.HasPrefix(from, "https://") {
		return os.Open(from)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, from, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "listo-image-import (+https://github.com/khinshankhan/listo)")
	res, err := fetchClient.Do(req)
	if err != nil {
		return nil, err
	}
	fail := func(format string, args ...any) (io.ReadCloser, error) {
		res.Body.Close()
		return nil, fmt.Errorf("fetching %s: %s", from, fmt.Sprintf(format, args...))
	}
	if res.StatusCode != http.StatusOK {
		return fail("status %d", res.StatusCode)
	}
	switch mt, _, _ := mime.ParseMediaType(res.Header.Get("Content-Type")); mt {
	case "image/jpeg", "image/png", "image/webp":
	default:
		return fail("content type %q isn't an image we read", mt)
	}
	if res.ContentLength > media.MaxBytes {
		return fail("%d bytes is too big", res.ContentLength)
	}
	return res.Body, nil
}

// readManifest reads a tab-separated manifest with a header row naming its
// columns: slug, kind, file, source and source_url. Files are relative to dir
// unless they're URLs; kind defaults to cover and source to staff.
func readManifest(r io.Reader, dir string) ([]imageSource, error) {
	sc := bufio.NewScanner(r)
	var cols map[string]int
	var out []imageSource
	for line := 1; sc.Scan(); line++ {
		text := sc.Text()
		if strings.TrimSpace(text) == "" || strings.HasPrefix(text, "#") {
			continue
		}
		fields := strings.Split(text, "\t")
		if cols == nil {
			cols = map[string]int{}
			for i, f := range fields {
				cols[strings.TrimSpace(f)] = i
			}
			for _, need := range []string{"slug", "file"} {
				if _, ok := cols[need]; !ok {
					return nil, fmt.Errorf("manifest header has no %s column", need)
				}
			}
			continue
		}
		get := func(name string) string {
			if i, ok := cols[name]; ok && i < len(fields) {
				return strings.TrimSpace(fields[i])
			}
			return ""
		}
		src := imageSource{Slug: get("slug"), Kind: get("kind"), From: get("file"), Source: get("source"), SourceURL: get("source_url")}
		if src.Slug == "" || src.From == "" {
			return nil, fmt.Errorf("manifest line %d: wants a slug and a file", line)
		}
		if src.Kind == "" {
			src.Kind = media.KindCover
		}
		if src.Source == "" {
			src.Source = "staff"
		}
		if !strings.Contains(src.From, "://") {
			src.From = filepath.Join(dir, src.From)
		}
		out = append(out, src)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// parseInterleaved parses flags wherever they sit among the positional
// arguments, which flag.Parse stops at.
func parseInterleaved(fs *flag.FlagSet, args []string) ([]string, error) {
	var pos []string
	for {
		if err := fs.Parse(args); err != nil {
			return nil, err
		}
		args = fs.Args()
		if len(args) == 0 {
			return pos, nil
		}
		pos = append(pos, args[0])
		args = args[1:]
	}
}
