package media_test

import (
	"context"
	"crypto/sha256"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/khinshankhan/listo/media"
)

// sha256 of no bytes
const empty = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

func TestKeys(t *testing.T) {
	sum := sha256.Sum256(nil)
	cases := []struct {
		name string
		got  string
		want string
	}{
		{"base", media.Key(media.KindCover, sum), "cover/e3/b0/" + empty},
		{"a width", media.VariantKey(media.Key(media.KindCover, sum), "320", "jpg"), "cover/e3/b0/" + empty + "-320.jpg"},
		{"the original", media.VariantKey(media.Key(media.KindBanner, sum), "orig", "png"), "banner/e3/b0/" + empty + "-orig.png"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.got != tc.want {
				t.Errorf("got %q, want %q", tc.got, tc.want)
			}
		})
	}
}

func TestOpen(t *testing.T) {
	cases := []struct {
		url    string
		ok     bool
		errHas string
	}{
		{"file://var/media", true, ""},
		{"file:///var/lib/listo/media", true, ""},
		{"file://", false, "no directory"},
		{"s3://bucket", false, "isn't built yet"},
		{"gs://bucket", false, "unknown scheme"},
		{"var/media", false, "want scheme"},
	}
	for _, tc := range cases {
		t.Run(tc.url, func(t *testing.T) {
			_, err := media.Open(tc.url)
			if tc.ok {
				if err != nil {
					t.Fatalf("err = %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tc.errHas) {
				t.Errorf("err = %v, want one saying %q", err, tc.errHas)
			}
		})
	}
}

func TestFileStoreRoundTrip(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	s := media.NewFileStore(dir, "/media/")
	key := "cover/e3/b0/" + empty + "-320.jpg"

	if err := s.Put(ctx, key, strings.NewReader("pixels"), "image/jpeg"); err != nil {
		t.Fatal(err)
	}
	// putting the same key again is fine: same key, same bytes
	if err := s.Put(ctx, key, strings.NewReader("pixels"), "image/jpeg"); err != nil {
		t.Fatal(err)
	}

	rc, err := s.Open(ctx, key)
	if err != nil {
		t.Fatal(err)
	}
	defer rc.Close()
	got, err := io.ReadAll(rc)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "pixels" {
		t.Errorf("read %q, want %q", got, "pixels")
	}
	if got := s.URL(key); got != "/media/"+key {
		t.Errorf("URL = %q", got)
	}

	// nothing left over from the writes but the image
	var files []string
	_ = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			rel, _ := filepath.Rel(dir, path)
			files = append(files, filepath.ToSlash(rel))
		}
		return err
	})
	if len(files) != 1 || files[0] != key {
		t.Errorf("files = %v, want just %s", files, key)
	}

	_, err = s.Open(ctx, "cover/e3/b0/"+empty+"-160.jpg")
	if !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("missing key: err = %v, want fs.ErrNotExist", err)
	}
}

type failing struct{}

func (failing) Read([]byte) (int, error) { return 0, errors.New("connection reset") }

func TestFileStorePutLeavesNothingOnAFailedRead(t *testing.T) {
	dir := t.TempDir()
	s := media.NewFileStore(dir, "/media/")
	key := "cover/e3/b0/" + empty + "-320.jpg"

	if err := s.Put(context.Background(), key, failing{}, "image/jpeg"); err == nil {
		t.Fatal("Put succeeded")
	}
	entries, err := os.ReadDir(filepath.Join(dir, "cover", "e3", "b0"))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Errorf("left behind %v", entries)
	}
}

func TestFileStoreRejectsBadKeys(t *testing.T) {
	ctx := context.Background()
	s := media.NewFileStore(t.TempDir(), "/media/")
	other := strings.Repeat("ab", 32)

	cases := []struct {
		name string
		key  string
	}{
		{"empty", ""},
		{"a base key, no variant", "cover/e3/b0/" + empty},
		{"traversal", "../../etc/passwd"},
		{"traversal inside the shape", "cover/e3/b0/../../" + empty + "-320.jpg"},
		{"unknown kind", "poster/e3/b0/" + empty + "-320.jpg"},
		{"a first dir that doesn't match the hash", "cover/ab/b0/" + empty + "-320.jpg"},
		{"a second dir that doesn't match the hash", "cover/e3/ab/" + empty + "-320.jpg"},
		{"short hash", "cover/e3/b0/e3b0-320.jpg"},
		{"uppercase hash", "cover/E3/B0/" + strings.ToUpper(empty) + "-320.jpg"},
		{"unknown extension", "cover/e3/b0/" + empty + "-320.gif"},
		{"word variant", "cover/e3/b0/" + empty + "-big.jpg"},
		{"leading slash", "/cover/ab/ab/" + other + "-320.jpg"},
		{"trailing junk", "cover/ab/ab/" + other + "-320.jpg/x"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := s.Put(ctx, tc.key, strings.NewReader("x"), "image/jpeg"); !errors.Is(err, media.ErrBadKey) {
				t.Errorf("Put: err = %v, want ErrBadKey", err)
			}
			if _, err := s.Open(ctx, tc.key); !errors.Is(err, media.ErrBadKey) {
				t.Errorf("Open: err = %v, want ErrBadKey", err)
			}
		})
	}
}
