// Covers fade in over their stand-in once they've loaded. Ones already
// loaded show at once; without this, every cover just appears. A cover that
// fails stays hidden, so the stand-in shows instead of a broken image.
const loaded = (img: HTMLImageElement): void => {
  img.dataset.loaded = "";
};

for (const img of document.querySelectorAll<HTMLImageElement>("img[data-cover]")) {
  if (img.complete && img.naturalWidth > 0) {
    loaded(img);
  }
}

// load doesn't bubble, so listen on the way down; this catches covers added
// later too
document.addEventListener(
  "load",
  (event: Event): void => {
    if (event.target instanceof HTMLImageElement && event.target.hasAttribute("data-cover")) {
      loaded(event.target);
    }
  },
  true,
);

// last, in the same task as the checks above, so no loaded cover blinks out
document.documentElement.classList.add("covers-fade");
