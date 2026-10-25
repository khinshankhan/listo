// The cross-page fade names the top bar and sidebar so they hold still, which
// only works when both pages sit at the top. Leave from further down, or for a
// page that opens scrolled (back, forward, a #fragment), and the old and new
// bars are a screen or more apart: the fade would slide them across the page,
// over everything, as it went. Then they fade with the rest instead.
window.addEventListener("pageswap", (event: PageSwapEvent): void => {
  if (!event.viewTransition) {
    return;
  }
  const to = event.activation;
  const still = window.scrollY === 0 && (to?.navigationType === "push" || to?.navigationType === "replace") && !new URL(to.entry?.url ?? "", window.location.href).hash;
  // set either way: a page restored from the back-forward cache keeps its styles
  for (const el of document.querySelectorAll<HTMLElement>('[class*="view-transition-name:"]')) {
    el.style.viewTransitionName = still ? "" : "none";
  }
});
