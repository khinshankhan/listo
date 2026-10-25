// The one script every page loads: it mounts each island placeholder when its
// load strategy says.
//
// NOTE: import only types from mount, or React lands in this chunk.
import type { Mount } from "./mount";
import "./shortcuts";
import "./search";
import "./spoilers";

type Loader = () => Promise<Mount>;

// One line per island.
const islands: Record<string, Loader> = {
  timeline: () => import("./islands/timeline").then((m) => m.mount),
};

const start = (el: HTMLElement): void => {
  const name = el.dataset.island ?? "";
  const load = islands[name];
  if (!load) {
    console.error(`unknown island "${name}"; keeping its fallback`);
    return;
  }
  let props: unknown;
  try {
    props = JSON.parse(el.dataset.props ?? "null");
  } catch (err: unknown) {
    console.error(`island "${name}" has unreadable props; keeping its fallback`, err);
    return;
  }
  load()
    .then((mount) => mount(el, props))
    .catch((err: unknown) => console.error(`island "${name}" failed to load; keeping its fallback`, err));
};

const onIdle = (run: () => void): void => {
  if ("requestIdleCallback" in window) {
    window.requestIdleCallback(run, { timeout: 2000 });
  } else {
    setTimeout(run, 200);
  }
};

const onVisible = (el: HTMLElement, run: () => void): void => {
  const observer = new IntersectionObserver(
    (seen) => {
      for (const entry of seen) {
        if (entry.isIntersecting) {
          observer.disconnect();
          run();
          return;
        }
      }
    },
    { rootMargin: "200px" },
  );
  observer.observe(el);
};

for (const el of document.querySelectorAll<HTMLElement>("[data-island]")) {
  const run = (): void => start(el);
  switch (el.dataset.load) {
    case "idle":
      onIdle(run);
      break;
    case "visible":
      onVisible(el, run);
      break;
    default:
      run();
  }
}
