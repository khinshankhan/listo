package media

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	_ "image/png" // decoders register themselves
	"io"
	"math"
	"strconv"

	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp"
)

const (
	// MaxBytes is the largest original Import takes.
	MaxBytes = 25 << 20
	// maxPixels stops a small file that decodes to a huge image.
	maxPixels = 50_000_000
	// jpegQuality is every variant's; covers keep their text sharp at it.
	jpegQuality = 80
)

// Image is an imported image, as its row records it.
type Image struct {
	Kind string
	// Key is the base key; Widths gives the variants under it.
	Key           string
	Width, Height int
	// Fit is "cover" or "contain", how it sits in its slot.
	Fit string
	// Ext is the original's format, kept unserved under -orig.
	Ext string
}

// Import stores an image of kind from r: the original as it came, then a
// JPEG per width in Widths. Variants are upright, sRGB and metadata free.
// The same bytes always land on the same keys, so a repeat is harmless.
//
// NOTE: ICC profiles are ignored, CMYK is converted naively; fine for covers,
// not for print.
func Import(ctx context.Context, store Store, r io.Reader, kind string) (Image, error) {
	if !knownKind(kind) {
		return Image{}, fmt.Errorf("unknown image kind %q", kind)
	}
	raw, err := io.ReadAll(io.LimitReader(r, MaxBytes+1))
	if err != nil {
		return Image{}, fmt.Errorf("reading image: %w", err)
	}
	if len(raw) > MaxBytes {
		return Image{}, fmt.Errorf("image is over %d MB", MaxBytes>>20)
	}

	cfg, format, err := image.DecodeConfig(bytes.NewReader(raw))
	if err != nil {
		return Image{}, fmt.Errorf("not an image we read (jpeg, png, webp): %w", err)
	}
	if cfg.Width*cfg.Height > maxPixels {
		return Image{}, fmt.Errorf("image is %dx%d, too many pixels", cfg.Width, cfg.Height)
	}
	src, _, err := image.Decode(bytes.NewReader(raw))
	if err != nil {
		return Image{}, fmt.Errorf("decoding %s: %w", format, err)
	}

	pic := flatten(src)
	if format == "jpeg" {
		pic = orient(pic, jpegOrientation(raw))
	}
	w, h := pic.Bounds().Dx(), pic.Bounds().Dy()
	img := Image{Kind: kind, Key: Key(kind, sha256.Sum256(raw)), Width: w, Height: h, Fit: Fit(kind, w, h), Ext: extOf(format)}
	if err := store.Put(ctx, VariantKey(img.Key, "orig", img.Ext), bytes.NewReader(raw), "image/"+format); err != nil {
		return Image{}, err
	}
	for _, vw := range Widths(kind, w) {
		var buf bytes.Buffer
		if err := jpeg.Encode(&buf, resize(pic, vw), &jpeg.Options{Quality: jpegQuality}); err != nil {
			return Image{}, fmt.Errorf("encoding %dpx: %w", vw, err)
		}
		if err := store.Put(ctx, VariantKey(img.Key, strconv.Itoa(vw), "jpg"), &buf, "image/jpeg"); err != nil {
			return Image{}, err
		}
	}
	return img, nil
}

// variantWidths are the widths each kind is shrunk to, smallest first. The
// widest cover slot is 210px, so 640 covers it at 3x and no screen needs the
// original.
var variantWidths = map[string][]int{
	KindCover:      {160, 320, 640},
	KindBanner:     {640, 1280, 1920},
	KindArt:        {640, 1280},
	KindScreenshot: {640, 1280},
}

// Widths are the variants an image of kind and width has: each of the kind's
// widths it can shrink to, then its own. Nothing is upscaled.
func Widths(kind string, width int) []int {
	var out []int
	for _, w := range variantWidths[kind] {
		if w < width {
			out = append(out, w)
		}
	}
	return append(out, width)
}

// coverAspect is every cover slot's width over height.
const coverAspect = 5.0 / 7.0

// coverSlack is how far off 5:7 a cover can be and still fill its slot.
const coverSlack = 0.10

// Fit decides how an image of kind sits in its slot. A cover near 5:7 fills
// it and any other is letterboxed, landscape boxes included; the kind is the
// caller's to decide. Other kinds fill their slots and are steered by a focal
// point.
func Fit(kind string, width, height int) string {
	if kind != KindCover || aspectOff(kind, width, height) <= coverSlack {
		return "cover"
	}
	return "contain"
}

// aspectOff is how far an image is from its kind's slot, as a fraction; kinds
// without a fixed slot are never off.
func aspectOff(kind string, width, height int) float64 {
	if kind != KindCover {
		return 0
	}
	return math.Abs(float64(width)/float64(height)/coverAspect - 1)
}

// coverFloor is the narrowest cover that can become canon on its own.
const coverFloor = 200

// Size is an image's width and height.
type Size struct {
	Width, Height int
}

// AutoCanon reports whether a new image of kind should replace cur, the
// current canon or nil, without a curator: the one closer to its slot's
// aspect wins, within a percent counting as a tie, then the larger. A cover
// under coverFloor never wins on its own.
func AutoCanon(kind string, next Size, cur *Size) bool {
	if kind == KindCover && next.Width < coverFloor {
		return false
	}
	if cur == nil {
		return true
	}
	dn, dc := aspectOff(kind, next.Width, next.Height), aspectOff(kind, cur.Width, cur.Height)
	if math.Abs(dn-dc) > 0.01 {
		return dn < dc
	}
	return next.Width*next.Height > cur.Width*cur.Height
}

func knownKind(kind string) bool {
	_, ok := variantWidths[kind]
	return ok
}

func extOf(format string) string {
	if format == "jpeg" {
		return "jpg"
	}
	return format
}

// flatten draws src as plain RGB: CMYK and grey convert, transparency lands
// on white since JPEG has none.
func flatten(src image.Image) *image.RGBA {
	b := src.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(dst, dst.Bounds(), image.NewUniform(color.White), image.Point{}, draw.Src)
	draw.Draw(dst, dst.Bounds(), src, b.Min, draw.Over)
	return dst
}

// lanczos3 keeps cover text sharper than CatmullRom when shrinking.
var lanczos3 = &draw.Kernel{Support: 3, At: func(t float64) float64 {
	if t == 0 {
		return 1
	}
	x := math.Pi * t
	return 3 * math.Sin(x) * math.Sin(x/3) / (x * x)
}}

// resize scales pic to width, keeping its aspect.
func resize(pic *image.RGBA, width int) image.Image {
	b := pic.Bounds()
	if width == b.Dx() {
		return pic
	}
	height := max(1, int(math.Round(float64(b.Dy())*float64(width)/float64(b.Dx()))))
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	lanczos3.Scale(dst, dst.Bounds(), pic, b, draw.Src, nil)
	return dst
}

// orient turns pic upright from its EXIF orientation, 1 to 8.
func orient(pic *image.RGBA, o int) *image.RGBA {
	if o < 2 || o > 8 {
		return pic
	}
	w, h := pic.Bounds().Dx(), pic.Bounds().Dy()
	dw, dh := w, h
	if o >= 5 {
		dw, dh = h, w
	}
	dst := image.NewRGBA(image.Rect(0, 0, dw, dh))
	for y := range dh {
		for x := range dw {
			var sx, sy int
			switch o {
			case 2:
				sx, sy = w-1-x, y
			case 3:
				sx, sy = w-1-x, h-1-y
			case 4:
				sx, sy = x, h-1-y
			case 5:
				sx, sy = y, x
			case 6:
				sx, sy = y, h-1-x
			case 7:
				sx, sy = w-1-y, h-1-x
			case 8:
				sx, sy = w-1-y, x
			}
			dst.SetRGBA(x, y, pic.RGBAAt(sx, sy))
		}
	}
	return dst
}

// jpegOrientation reads the EXIF orientation tag and nothing else; 1,
// upright, when there's none.
func jpegOrientation(b []byte) int {
	if len(b) < 4 || b[0] != 0xFF || b[1] != 0xD8 {
		return 1
	}
	for i := 2; i+4 <= len(b); {
		if b[i] != 0xFF {
			return 1
		}
		switch marker := b[i+1]; {
		case marker == 0xFF:
			// fill byte
			i++
			continue
		case marker == 0xDA || marker == 0xD9:
			// the image data starts, metadata is over
			return 1
		case marker == 0x01 || marker >= 0xD0 && marker <= 0xD8:
			// no length
			i += 2
			continue
		}
		n := int(binary.BigEndian.Uint16(b[i+2:]))
		if n < 2 || i+2+n > len(b) {
			return 1
		}
		if seg := b[i+4 : i+2+n]; b[i+1] == 0xE1 && bytes.HasPrefix(seg, []byte("Exif\x00\x00")) {
			return tiffOrientation(seg[6:])
		}
		i += 2 + n
	}
	return 1
}

// tiffOrientation finds tag 0x0112 in the first IFD of a TIFF header.
func tiffOrientation(t []byte) int {
	if len(t) < 8 {
		return 1
	}
	var order binary.ByteOrder
	switch string(t[:2]) {
	case "II":
		order = binary.LittleEndian
	case "MM":
		order = binary.BigEndian
	default:
		return 1
	}
	ifd := int(order.Uint32(t[4:]))
	if ifd < 8 || ifd+2 > len(t) {
		return 1
	}
	for k := range int(order.Uint16(t[ifd:])) {
		e := ifd + 2 + 12*k
		if e+12 > len(t) {
			return 1
		}
		// a SHORT, 1 to 8
		if order.Uint16(t[e:]) == 0x0112 && order.Uint16(t[e+2:]) == 3 {
			if o := int(order.Uint16(t[e+8:])); o >= 1 && o <= 8 {
				return o
			}
			return 1
		}
	}
	return 1
}
