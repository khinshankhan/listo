// Package media stores images. Callers see keys, never paths, so the bytes
// can move from disk to a bucket without touching them.
package media

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// Store holds images by key.
type Store interface {
	// Put writes r under key; a reader that fails leaves nothing behind.
	Put(ctx context.Context, key string, r io.Reader, contentType string) error
	// Open reads key. A missing key is an error matching fs.ErrNotExist.
	Open(ctx context.Context, key string) (io.ReadCloser, error)
	// URL is where a browser fetches key.
	URL(key string) string
}

// ErrBadKey is a key not shaped like VariantKey's.
var ErrBadKey = errors.New("media: bad key")

// Kinds an image can be.
const (
	KindCover      = "cover"
	KindBanner     = "banner"
	KindArt        = "art"
	KindScreenshot = "screenshot"
)

// Open returns the store for url: file://<dir> for now, s3:// later.
func Open(url string) (Store, error) {
	scheme, rest, ok := strings.Cut(url, "://")
	if !ok {
		return nil, fmt.Errorf("storage url %q: want scheme://...", url)
	}
	switch scheme {
	case "file":
		if rest == "" {
			return nil, fmt.Errorf("storage url %q: no directory", url)
		}
		return NewFileStore(rest, "/media/"), nil
	case "s3":
		return nil, fmt.Errorf("storage url %q: s3 isn't built yet, use file://", url)
	default:
		return nil, fmt.Errorf("storage url %q: unknown scheme %q", url, scheme)
	}
}

// Key is the base key for an image of kind whose original bytes hash to sum,
// eg cover/3f/a1/3fa1...; every variant's key starts with it.
func Key(kind string, sum [32]byte) string {
	h := hex.EncodeToString(sum[:])
	return fmt.Sprintf("%s/%s/%s/%s", kind, h[0:2], h[2:4], h)
}

// VariantKey is one stored file of base, eg base-320.jpg or base-orig.png.
func VariantKey(base, variant, ext string) string {
	return fmt.Sprintf("%s-%s.%s", base, variant, ext)
}

var keyShape = regexp.MustCompile(`^(cover|banner|art|screenshot)/([0-9a-f]{2})/([0-9a-f]{2})/([0-9a-f]{64})-([0-9]+|orig)\.(jpg|png|webp)$`)

// checkKey rejects anything not shaped like VariantKey's output, which also
// keeps keys from walking out of a directory.
func checkKey(key string) error {
	m := keyShape.FindStringSubmatch(key)
	if m == nil || m[2] != m[4][0:2] || m[3] != m[4][2:4] {
		return fmt.Errorf("%w: %q", ErrBadKey, key)
	}
	return nil
}
