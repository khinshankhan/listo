// Browse's format filter only has formats once there's a type. Picking one
// swaps in that type's formats from the form's templates straight away, rather
// than leaving the filter on "pick a type first" until the form is applied.
document.addEventListener("change", (event: Event): void => {
  const type = event.target;
  if (!(type instanceof HTMLSelectElement) || type.id !== "filter-medium") {
    return;
  }
  const format = document.getElementById("filter-format");
  const options = type.form?.querySelector<HTMLTemplateElement>(`template[data-formats="${CSS.escape(type.value)}"]`);
  if (!(format instanceof HTMLSelectElement) || !options) {
    return;
  }
  format.replaceChildren(options.content.cloneNode(true));
  format.disabled = type.value === "";
});
