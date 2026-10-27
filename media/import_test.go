package media_test

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/khinshankhan/listo/media"
)

func TestWidths(t *testing.T) {
	cases := []struct {
		kind  string
		width int
		want  []int
	}{
		{media.KindCover, 460, []int{160, 320, 460}},
		{media.KindCover, 320, []int{160, 320}},
		{media.KindCover, 150, []int{150}},
		{media.KindCover, 640, []int{160, 320, 640}},
		{media.KindCover, 1000, []int{160, 320, 640, 1000}},
		{media.KindCover, 3840, []int{160, 320, 640, 3840}},
		{media.KindBanner, 1920, []int{640, 1280, 1920}},
		{media.KindBanner, 2560, []int{640, 1280, 1920, 2560}},
		{media.KindArt, 900, []int{640, 900}},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("%s %dpx", tc.kind, tc.width), func(t *testing.T) {
			if got := media.Widths(tc.kind, tc.width); !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestFit(t *testing.T) {
	cases := []struct {
		name          string
		kind          string
		width, height int
		want          string
	}{
		{"exactly 5:7", media.KindCover, 500, 700, "cover"},
		{"an anilist-ish cover", media.KindCover, 460, 650, "cover"},
		{"a dvd case, 9% narrow", media.KindCover, 1300, 2000, "cover"},
		{"2:3, 7% narrow", media.KindCover, 400, 600, "cover"},
		{"a tall paperback, 20% narrow", media.KindCover, 400, 700, "contain"},
		{"square, 40% wide", media.KindCover, 600, 600, "contain"},
		{"a landscape box is still a cover", media.KindCover, 640, 480, "contain"},
		{"16:9 too, if the caller calls it a cover", media.KindCover, 1920, 1080, "contain"},
		{"banners always fill", media.KindBanner, 1920, 400, "cover"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := media.Fit(tc.kind, tc.width, tc.height); got != tc.want {
				t.Errorf("got %s, want %s", got, tc.want)
			}
		})
	}
}

func TestAutoCanon(t *testing.T) {
	size := func(w, h int) *media.Size { return &media.Size{Width: w, Height: h} }
	cases := []struct {
		name string
		kind string
		next media.Size
		cur  *media.Size
		want bool
	}{
		{"the first cover", media.KindCover, media.Size{Width: 460, Height: 644}, nil, true},
		{"the first cover, but tiny", media.KindCover, media.Size{Width: 150, Height: 210}, nil, false},
		{"just over the floor", media.KindCover, media.Size{Width: 200, Height: 280}, nil, true},
		{"closer to 5:7 beats bigger", media.KindCover, media.Size{Width: 300, Height: 420}, size(1000, 1500), true},
		{"bigger doesn't beat closer", media.KindCover, media.Size{Width: 1000, Height: 1500}, size(300, 420), false},
		{"same shape, bigger wins", media.KindCover, media.Size{Width: 920, Height: 1288}, size(460, 644), true},
		{"same shape, smaller loses", media.KindCover, media.Size{Width: 460, Height: 644}, size(920, 1288), false},
		{"half a percent off is a tie, so bigger wins", media.KindCover, media.Size{Width: 920, Height: 1295}, size(460, 644), true},
		{"a tie in size keeps what's there", media.KindCover, media.Size{Width: 460, Height: 644}, size(460, 644), false},
		{"banners: just bigger", media.KindBanner, media.Size{Width: 1920, Height: 1080}, size(1280, 720), true},
		{"a small banner isn't held to the cover floor", media.KindBanner, media.Size{Width: 150, Height: 50}, nil, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := media.AutoCanon(tc.kind, tc.next, tc.cur); got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

// importFile imports a testdata file into a fresh store.
func importFile(t *testing.T, name, kind string) (media.Store, media.Image) {
	t.Helper()
	raw, err := os.ReadFile("testdata/" + name)
	if err != nil {
		t.Fatal(err)
	}
	store := media.NewFileStore(t.TempDir(), "/media/")
	img, err := media.Import(context.Background(), store, bytes.NewReader(raw), kind)
	if err != nil {
		t.Fatal(err)
	}
	return store, img
}

// variant decodes one stored variant.
func variant(t *testing.T, store media.Store, img media.Image, width int) image.Image {
	t.Helper()
	rc, err := store.Open(context.Background(), media.VariantKey(img.Key, fmt.Sprint(width), "jpg"))
	if err != nil {
		t.Fatal(err)
	}
	defer rc.Close()
	pic, err := jpeg.Decode(rc)
	if err != nil {
		t.Fatal(err)
	}
	return pic
}

// near is a JPEG's take on want: close, never exact.
func near(c color.Color, want color.RGBA) bool {
	r, g, b, _ := c.RGBA()
	d := func(a uint32, b uint8) bool { return int(a>>8)-int(b) < 40 && int(b)-int(a>>8) < 40 }
	return d(r, want.R) && d(g, want.G) && d(b, want.B)
}

var (
	red   = color.RGBA{R: 255, A: 255}
	blue  = color.RGBA{B: 255, A: 255}
	green = color.RGBA{G: 255, A: 255}
)

func TestImportTurnsPhotosUpright(t *testing.T) {
	// both are 40x20, red then blue, tagged to turn a quarter one way or the
	// other, in either byte order
	cases := []struct {
		file string
		top  color.RGBA
	}{
		{"orient6-mm.jpg", red},
		{"orient8-ii.jpg", blue},
	}
	for _, tc := range cases {
		t.Run(tc.file, func(t *testing.T) {
			store, img := importFile(t, tc.file, media.KindCover)
			// upright it's tall, so still a cover, though too thin to fill the slot
			if img.Width != 20 || img.Height != 40 || img.Kind != media.KindCover || img.Fit != "contain" {
				t.Errorf("got %+v, want a 20x40 cover that's contained", img)
			}
			pic := variant(t, store, img, 20)
			if b := pic.Bounds(); b.Dx() != 20 || b.Dy() != 40 {
				t.Fatalf("variant is %v", b)
			}
			if c := pic.At(10, 5); !near(c, tc.top) {
				t.Errorf("top is %v, want %v", c, tc.top)
			}
		})
	}
}

func TestImportConvertsCMYK(t *testing.T) {
	store, img := importFile(t, "cmyk.jpg", media.KindCover)
	if c := variant(t, store, img, 20).At(10, 10); !near(c, red) {
		t.Errorf("got %v, want red", c)
	}
}

func TestImportFormats(t *testing.T) {
	cases := []struct {
		file, ext string
	}{
		{"png.png", "png"},
		{"webp.webp", "webp"},
		{"cmyk.jpg", "jpg"},
	}
	for _, tc := range cases {
		t.Run(tc.file, func(t *testing.T) {
			store, img := importFile(t, tc.file, media.KindCover)
			if img.Ext != tc.ext {
				t.Errorf("ext = %q, want %q", img.Ext, tc.ext)
			}
			// the original is kept as it came
			raw, _ := os.ReadFile("testdata/" + tc.file)
			rc, err := store.Open(context.Background(), media.VariantKey(img.Key, "orig", tc.ext))
			if err != nil {
				t.Fatal(err)
			}
			defer rc.Close()
			var got bytes.Buffer
			_, _ = got.ReadFrom(rc)
			if !bytes.Equal(got.Bytes(), raw) {
				t.Error("original changed on the way in")
			}
		})
	}
	store, img := importFile(t, "webp.webp", media.KindCover)
	if c := variant(t, store, img, 50).At(25, 35); !near(c, green) {
		t.Errorf("webp decoded as %v, want green", c)
	}
}

func TestImportShrinks(t *testing.T) {
	var buf bytes.Buffer
	src := image.NewRGBA(image.Rect(0, 0, 460, 644))
	for y := range 644 {
		for x := range 460 {
			src.SetRGBA(x, y, blue)
		}
	}
	if err := png.Encode(&buf, src); err != nil {
		t.Fatal(err)
	}
	store := media.NewFileStore(t.TempDir(), "/media/")
	img, err := media.Import(context.Background(), store, &buf, media.KindCover)
	if err != nil {
		t.Fatal(err)
	}
	if img.Width != 460 || img.Height != 644 || img.Fit != "cover" || !strings.HasPrefix(img.Key, "cover/") {
		t.Errorf("got %+v", img)
	}

	for _, w := range []struct{ width, height int }{{160, 224}, {320, 448}, {460, 644}} {
		pic := variant(t, store, img, w.width)
		if b := pic.Bounds(); b.Dx() != w.width || b.Dy() != w.height {
			t.Errorf("%dpx variant is %dx%d, want %dx%d", w.width, b.Dx(), b.Dy(), w.width, w.height)
		}
		if c := pic.At(w.width/2, w.height/2); !near(c, blue) {
			t.Errorf("%dpx variant is %v in the middle, want blue", w.width, c)
		}
	}
	if _, err := store.Open(context.Background(), media.VariantKey(img.Key, "640", "jpg")); err == nil {
		t.Error("a 640 variant was made from a 460 original")
	}
}

func TestImportRejects(t *testing.T) {
	store := media.NewFileStore(t.TempDir(), "/media/")
	cases := []struct {
		name, kind string
		body       string
	}{
		{"not an image", media.KindCover, "<html>404</html>"},
		{"a kind that isn't", "poster", "whatever"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := media.Import(context.Background(), store, strings.NewReader(tc.body), tc.kind); err == nil {
				t.Error("imported")
			}
		})
	}
}
