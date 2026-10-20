// The top bar's search suggests as you type. Without this the form still
// works: Enter goes to the full search page.

type Hit = { title: string; href: string; meta?: string; matched?: string };
type Suggestions = { query: string; results: Hit[]; all: string };

const DEBOUNCE_MS = 120;
const LIMIT = 6;

const input = document.getElementById("site-search");
const panel = document.getElementById("site-search-panel");
const list = document.getElementById("site-search-results");
const all = document.getElementById("site-search-all");

if (input instanceof HTMLInputElement && panel && list && all instanceof HTMLAnchorElement) {
  suggest(input, panel, list, all);
}

function suggest(input: HTMLInputElement, panel: HTMLElement, list: HTMLElement, all: HTMLAnchorElement): void {
  const form = input.form;
  let hits: Hit[] = [];
  let active = -1;
  let timer: ReturnType<typeof setTimeout> | undefined;
  let pending: AbortController | undefined;

  const open = (): void => {
    panel.hidden = false;
    input.setAttribute("aria-expanded", "true");
  };

  const close = (): void => {
    panel.hidden = true;
    input.setAttribute("aria-expanded", "false");
    highlight(-1);
  };

  const highlight = (i: number): void => {
    active = i;
    list.querySelectorAll<HTMLElement>("[role=option]").forEach((option, j) => {
      const on = j === i;
      option.setAttribute("aria-selected", String(on));
      option.dataset.active = String(on);
      if (on) {
        option.scrollIntoView({ block: "nearest" });
      }
    });
    if (i < 0) {
      input.removeAttribute("aria-activedescendant");
    } else {
      input.setAttribute("aria-activedescendant", `site-search-option-${i}`);
    }
  };

  const render = (data: Suggestions): void => {
    hits = data.results;
    const rows = hits.map((hit, i) => {
      const option = document.createElement("li");
      option.id = `site-search-option-${i}`;
      option.setAttribute("role", "option");
      option.setAttribute("aria-selected", "false");
      option.className =
        "flex cursor-pointer items-baseline gap-3 px-3 py-2 text-14 hover:bg-background-1 data-[active=true]:bg-accent-3";

      // the title gets the width; an alias match sits under it
      const names = document.createElement("span");
      names.className = "flex min-w-0 flex-col";
      const title = document.createElement("span");
      title.className = "truncate font-semibold text-foreground-strong";
      title.textContent = hit.title;
      names.append(title);
      if (hit.matched) {
        const matched = document.createElement("span");
        matched.className = "truncate text-12 text-foreground-muted";
        matched.textContent = `as ${hit.matched}`;
        names.append(matched);
      }
      option.append(names);
      if (hit.meta) {
        const meta = document.createElement("span");
        meta.className = "ml-auto shrink-0 font-mono text-12 text-foreground-muted";
        meta.textContent = hit.meta;
        option.append(meta);
      }

      // keep focus in the field, so the keyboard still drives the list
      option.addEventListener("pointerdown", (event) => event.preventDefault());
      option.addEventListener("click", () => window.location.assign(hit.href));
      return option;
    });

    if (rows.length === 0) {
      const empty = document.createElement("li");
      empty.setAttribute("role", "presentation");
      empty.className = "px-3 py-2.5 text-13 text-foreground-muted";
      empty.textContent = `nothing matches “${data.query}” yet`;
      rows.push(empty);
    }
    list.replaceChildren(...rows);
    all.href = data.all;
    all.textContent = `see all results for “${data.query}”`;
    active = -1;
    open();
  };

  const fetchSuggestions = async (): Promise<void> => {
    const q = input.value.trim();
    pending?.abort();
    if (q === "") {
      close();
      return;
    }
    const request = new AbortController();
    pending = request;
    try {
      const res = await fetch(`/api/search?q=${encodeURIComponent(q)}&limit=${LIMIT}`, { signal: request.signal });
      if (!res.ok) {
        close();
        return;
      }
      const data = (await res.json()) as Suggestions;
      // a reply for what the field no longer says is stale
      if (!request.signal.aborted && input.value.trim() === q) {
        render(data);
      }
    } catch {
      // aborted, or offline: the form still works
    }
  };

  input.addEventListener("input", () => {
    clearTimeout(timer);
    timer = setTimeout(() => void fetchSuggestions(), DEBOUNCE_MS);
  });

  input.addEventListener("focus", () => {
    if (input.value.trim() !== "" && list.childElementCount > 0) {
      open();
    }
  });

  input.addEventListener("keydown", (event) => {
    if (event.key === "Escape" && !panel.hidden) {
      event.preventDefault();
      close();
      return;
    }
    if (hits.length === 0) {
      return;
    }
    switch (event.key) {
      case "ArrowDown":
        event.preventDefault();
        open();
        highlight((active + 1) % hits.length);
        break;
      case "ArrowUp":
        event.preventDefault();
        open();
        highlight(active <= 0 ? hits.length - 1 : active - 1);
        break;
      case "Enter":
        // with nothing highlighted, Enter submits the form: the full search
        if (!panel.hidden && active >= 0) {
          event.preventDefault();
          window.location.assign(hits[active].href);
        }
        break;
    }
  });

  document.addEventListener("pointerdown", (event) => {
    if (form && event.target instanceof Node && !form.contains(event.target)) {
      close();
    }
  });
  form?.addEventListener("submit", close);
}
