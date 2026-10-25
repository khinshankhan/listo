// Tabbing along a sideways scroller can land on a card half out of view:
// browsers only scroll focus in when it's hidden outright. Bring it the rest
// of the way, as far as the scroller's scroll padding asks. Keyboard focus
// only, since a click lands on what's under the pointer and scrolling first
// would move it away.
document.addEventListener("focusin", (event: FocusEvent): void => {
  const el = event.target;
  if (!(el instanceof HTMLElement) || !el.parentElement?.closest('[role="region"]') || !el.matches(":focus-visible")) {
    return;
  }
  el.scrollIntoView({ block: "nearest", inline: "nearest" });
});
