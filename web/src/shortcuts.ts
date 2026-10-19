// "/" focuses the top bar search unless something else has focus for typing.
document.addEventListener("keydown", (event: KeyboardEvent): void => {
  if (event.key !== "/" || event.metaKey || event.ctrlKey || event.altKey) {
    return;
  }
  const target = event.target;
  if (target instanceof HTMLElement && target.closest("input, textarea, select, [contenteditable]")) {
    return;
  }
  const search = document.getElementById("site-search");
  if (search instanceof HTMLInputElement) {
    event.preventDefault();
    search.focus();
  }
});
