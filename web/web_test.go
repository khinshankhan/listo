package web

import (
	"reflect"
	"testing"
	"testing/fstest"
)

func TestLoadPreloadsTheFontsEveryPageUses(t *testing.T) {
	manifest := `{
	  "src/app.css": {"file": "assets/app-abc.css", "assets": [
	    "assets/manrope-latin-wght-normal-AAA.woff2",
	    "assets/manrope-latin-ext-wght-normal-BBB.woff2",
	    "assets/manrope-cyrillic-wght-normal-CCC.woff2",
	    "assets/ibm-plex-mono-latin-400-normal-DDD.woff2",
	    "assets/ibm-plex-mono-latin-400-normal-EEE.woff",
	    "assets/ibm-plex-mono-latin-500-normal-FFF.woff2",
	    "assets/ibm-plex-mono-latin-ext-400-normal-GGG.woff2"
	  ]},
	  "src/islands.ts": {"file": "assets/islands-xyz.js"}
	}`
	b, err := load(fstest.MapFS{".vite/manifest.json": {Data: []byte(manifest)}})
	if err != nil {
		t.Fatal(err)
	}

	// latin woff2 at the weights pages use, nothing else
	want := []string{"/assets/manrope-latin-wght-normal-AAA.woff2", "/assets/ibm-plex-mono-latin-400-normal-DDD.woff2"}
	if !reflect.DeepEqual(b.Tags.Fonts, want) {
		t.Errorf("Fonts = %v, want %v", b.Tags.Fonts, want)
	}
	if b.Tags.Script != "/assets/islands-xyz.js" || !reflect.DeepEqual(b.Tags.CSS, []string{"/assets/app-abc.css"}) {
		t.Errorf("tags = %+v", b.Tags)
	}
}
