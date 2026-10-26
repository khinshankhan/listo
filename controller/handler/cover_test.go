package handler

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/khinshankhan/listo/media"
	"github.com/khinshankhan/listo/views"
)

func TestCoverOf(t *testing.T) {
	store := media.NewFileStore(t.TempDir(), "/media/")
	key := "cover/3f/a1/" + strings.Repeat("3fa1", 16)
	v := func(w string) string { return "/media/" + key + "-" + w + ".jpg" }
	ptr := func(n int32) *int32 { return &n }
	str := func(s string) *string { return &s }
	start := time.Date(1988, 12, 24, 0, 0, 0, 0, time.UTC)

	cases := []struct {
		name          string
		key           *string
		width, height *int32
		fit           *string
		want          *views.CoverImage
	}{
		{"no image, just the stand-in", nil, nil, nil, nil, nil},
		{"wide enough for every variant", str(key), ptr(600), ptr(900), str("cover"), &views.CoverImage{
			Src: v("320"), Srcset: v("160") + " 160w, " + v("320") + " 320w, " + v("600") + " 600w", Width: 600, Height: 900}},
		{"a small box, letterboxed", str(key), ptr(256), ptr(324), str("contain"), &views.CoverImage{
			Src: v("256"), Srcset: v("160") + " 160w, " + v("256") + " 256w", Width: 256, Height: 324, Contain: true}},
		{"too small to shrink", str(key), ptr(150), ptr(210), str("cover"), &views.CoverImage{
			Src: v("150"), Srcset: v("150") + " 150w", Width: 150, Height: 210}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := coverOf(store, "game", "Mega Man 2", &start, tc.key, tc.width, tc.height, tc.fit)
			want := views.Cover{Medium: "game", Title: "Mega Man 2", Year: "1988", Image: tc.want}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("got  %+v %+v\nwant %+v %+v", got, got.Image, want, want.Image)
			}
		})
	}
}
