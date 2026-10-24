package handler

import (
	"math"
	"testing"
)

// Widths a browser drew these names at, to the pixel.
func TestNameWidth(t *testing.T) {
	cases := []struct {
		name string
		want float64
	}{
		{"Anonymous;Code", 111},
		{"Robotics;Notes", 97},
		{"Ben 10", 40},
		{"Railgun: Misaka-san wa Ima Chuumoku no Mato desukara", 352},
		{"とある", 39},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := nameFont.width(c.name); math.Abs(got-c.want) > 1 {
				t.Errorf("width = %.1f, want %.0f", got, c.want)
			}
		})
	}
}
