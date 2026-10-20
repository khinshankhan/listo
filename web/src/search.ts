// The top bar's search suggests as you type. Without this the form still
// works: Enter goes to the full search page.

type Hit = { title: string; href: string; meta?: string; matched?: string };
type Suggestions = { query: string; results: Hit[]; all: string };

const DEBOUNCE_MS = 120;
// a reply slower than this shows "searching", faster ones never flicker it
const SLOW_MS = 250;
const LIMIT = 6;

const input = document.getElementById("site-search");
const panel = document.getElementById("site-search-panel");
const list = document.getElementById("site-search-results");
const all = document.getElementById("site-search-all");
const announcer = document.getElementById("site-search-status");
const note = document.getElementById("site-search-note");

if (input instanceof HTMLInputElement && panel && list && all instanceof HTMLAnchorElement && announcer && note) {
  suggest(input, panel, list, all, announcer, note);
}

function suggest(
  input: HTMLInputElement,
  panel: HTMLElement,
  list: HTMLElement,
  all: HTMLAnchorElement,
  announcer: HTMLElement,
  note: HTMLElement,
): void {
  const form = input.form;
  let hits: Hit[] = [];
  let active = -1;
  let timer: ReturnType<typeof setTimeout> | undefined;
  let pending: AbortController | undefined;

  // the combobox is expanded only while it has options to offer; a note alone
  // isn't a list
  const open = (): void => {
    panel.hidden = false;
    input.setAttribute("aria-expanded", String(hits.length > 0));
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

  // message shows a note in place of the list: nothing found, still
  // searching, or a failure
  const message = (text: string): void => {
    hits = [];
    list.replaceChildren();
    list.hidden = true;
    note.textContent = text;
    note.hidden = false;
    highlight(-1);
    open();
  };

  const render = (data: Suggestions): void => {
    hits = data.results;
    const rows = hits.map((hit, i) => {
      const option = document.createElement("li");
      option.id = `site-search-option-${i}`;
      option.setAttribute("role", "option");
      option.setAttribute("aria-selected", "false");
      option.className =
        "flex cursor-pointer items-baseline gap-3 px-3 py-2 text-14 hover:bg-background-1 data-[active=true]:bg-accent-3 data-[active=true]:shadow-[inset_3px_0_0_var(--color-accent-9)]";

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
      message(`nothing matches “${data.query}” yet`);
    } else {
      note.hidden = true;
      list.hidden = false;
      list.replaceChildren(...rows);
    }
    // naming the query makes each announcement differ: screen readers skip a
    // live region that repeats itself, eg two searches with six hits each
    announcer.textContent =
      hits.length === 0
        ? `no suggestions for “${data.query}”`
        : `${hits.length} suggestion${hits.length === 1 ? "" : "s"} for “${data.query}”, up and down to choose`;
    active = -1;
    open();
  };

  const fetchSuggestions = async (): Promise<void> => {
    const q = input.value.trim();
    pending?.abort();
    if (q === "") {
      announcer.textContent = "";
      close();
      return;
    }
    // the full search follows the field, whatever state the panel is in
    all.href = `/search?q=${encodeURIComponent(q)}`;
    all.textContent = `see all results for “${q}”`;
    const request = new AbortController();
    pending = request;
    // keep showing the last results while the next ones load; only an empty
    // panel needs a "searching" line
    const slow = setTimeout(() => {
      list.setAttribute("aria-busy", "true");
      if (hits.length === 0) {
        message("searching…");
      }
    }, SLOW_MS);
    // a reply for what the field no longer says is stale
    const current = (): boolean => !request.signal.aborted && input.value.trim() === q;
    try {
      const res = await fetch(`/api/search?q=${encodeURIComponent(q)}&limit=${LIMIT}`, { signal: request.signal });
      if (!res.ok) {
        throw new Error(`search: ${res.status}`);
      }
      const data = (await res.json()) as Suggestions;
      if (current()) {
        render(data);
      }
    } catch {
      // cancelled because the field changed: the newer request answers
      if (current()) {
        message("couldn't load suggestions, enter still searches");
        announcer.textContent = `couldn't load suggestions for “${q}”, enter still searches`;
      }
    } finally {
      clearTimeout(slow);
      if (pending === request) {
        list.removeAttribute("aria-busy");
      }
    }
  };

  input.addEventListener("input", () => {
    clearTimeout(timer);
    timer = setTimeout(() => void fetchSuggestions(), DEBOUNCE_MS);
  });

  input.addEventListener("focus", () => {
    if (input.value.trim() !== "" && (hits.length > 0 || !note.hidden)) {
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
  // tabbing out of the search closes it; clicks on options keep focus in the
  // field, and the "see all" link is inside the form
  form?.addEventListener("focusout", (event) => {
    if (!(event.relatedTarget instanceof Node && form.contains(event.relatedTarget))) {
      close();
    }
  });
}
