// The timeline's "show spoilers" box keeps ?spoilers=show in the address, and
// in the links that sort the entries, so a shared or sorted page opens the way
// it was left.
const withSpoilers = (href: string, show: boolean): URL => {
  const url = new URL(href, window.location.href);
  if (show) {
    url.searchParams.set("spoilers", "show");
  } else {
    url.searchParams.delete("spoilers");
  }
  return url;
};

document.addEventListener("change", (event: Event): void => {
  const box = event.target;
  if (!(box instanceof HTMLInputElement) || !box.classList.contains("timeline-reveal")) {
    return;
  }
  history.replaceState(history.state, "", withSpoilers(window.location.href, box.checked));
  for (const link of document.querySelectorAll<HTMLAnchorElement>('nav[aria-label="group entries"] a')) {
    const url = withSpoilers(link.href, box.checked);
    link.href = `${url.pathname}${url.search}${url.hash}`;
  }
});
