import { BaseEdge, EdgeLabelRenderer, Handle, MarkerType, Position, ReactFlow, type Edge, type EdgeProps, type Node, type NodeProps } from "@xyflow/react";
import "@xyflow/react/dist/base.css";
import { useState, type ReactNode } from "react";
import * as z from "zod/mini";
import { island } from "../mount";
import "./timeline.css";

// The layout is the server's, lines and all, the same one its SVG fallback
// draws: x is time, y the row each work landed on.
const Props = z.object({
  width: z.number(),
  height: z.number(),
  ticks: z.array(z.object({ x: z.number(), label: z.string() })),
  boxes: z.nullable(z.array(z.object({ x: z.number(), y: z.number(), w: z.number(), h: z.number(), name: z.string() }))),
  works: z.array(
    z.object({
      id: z.number(),
      name: z.string(),
      href: z.string(),
      card: z.object({ x: z.number(), y: z.number(), w: z.number(), h: z.number() }),
      y: z.number(),
      dots: z.array(z.object({ x: z.number(), label: z.string() })),
    }),
  ),
  links: z.nullable(
    z.array(
      z.object({
        source: z.number(),
        target: z.number(),
        path: z.string(),
        label: z.string(),
        labelX: z.number(),
        labelY: z.number(),
        labelW: z.number(),
        labelEnd: z.boolean(),
        sequel: z.boolean(),
        spoiler: z.boolean(),
      }),
    ),
  ),
});
type Props = z.infer<typeof Props>;
type Work = Props["works"][number];

type WorkNode = Node<{ work: Work; light: (id: string | null) => void }, "work">;
type BoxNode = Node<{ name: string }, "box">;
type TickNode = Node<{ label: string }, "tick">;
type Line = NonNullable<Props["links"]>[number];
type LinkEdge = Edge<{ link: Line; dim: boolean }, "link">;

const WorkCard = ({ id, data: { work, light } }: NodeProps<WorkNode>): ReactNode => {
  const { card } = work;
  const row = work.y - card.y;
  const first = work.dots[0].x - card.x;
  const last = work.dots[work.dots.length - 1].x - card.x;
  return (
    <div className="timeline-card relative rounded-lg border border-border bg-background-2" style={{ width: card.w, height: card.h }}>
      <a
        href={work.href}
        onFocus={() => light(id)}
        onBlur={() => light(null)}
        className="absolute top-[7px] right-[10px] left-[10px] truncate text-13 leading-4 font-semibold text-foreground-strong after:absolute after:inset-[-7px_-10px_-27px] hover:underline focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-stark-contrast"
      >
        {work.name}
      </a>
      {last > first && <span className="absolute h-[3px] -translate-y-1/2 rounded-full bg-foreground-strong" style={{ left: first, width: last - first, top: row }} aria-hidden="true" />}
      {work.dots.map((d, i) => (
        <span key={i} aria-hidden="true">
          <span className="absolute size-[9px] -translate-1/2 rounded-full border-[1.5px] border-foreground-strong bg-background-2" style={{ left: d.x - card.x, top: row }} />
          <span className="absolute -translate-x-1/2 text-11 leading-[14px] whitespace-nowrap text-foreground-muted" style={{ left: d.x - card.x, top: row + 7 }}>
            {d.label}
          </span>
        </span>
      ))}
      {/* the server draws the lines; these only satisfy xyflow */}
      <Handle type="source" position={Position.Right} isConnectable={false} />
      <Handle type="target" position={Position.Left} isConnectable={false} />
    </div>
  );
};

const Box = ({ data, width, height }: NodeProps<BoxNode>): ReactNode => (
  <div className="rounded-xl border border-dashed border-border" style={{ width, height }}>
    <span className="absolute top-[13px] left-[11px] bg-background-2 px-1.5 text-14 leading-[18px] font-bold text-foreground-strong">{data.name}</span>
  </div>
);

const Tick = ({ data, height }: NodeProps<TickNode>): ReactNode => (
  <div className="flex flex-col items-center" style={{ height }} aria-hidden="true">
    <span className="w-px flex-1 border-l border-dotted border-border-subtle" />
    <span className="pt-1 font-mono text-[10px] leading-3 text-foreground-muted">{data.label}</span>
  </div>
);

// Link draws the server's line and its label.
const Link = ({ id, markerEnd, style, data }: EdgeProps<LinkEdge>): ReactNode => {
  if (!data) {
    return null;
  }
  const { link, dim } = data;
  const left = link.labelEnd ? link.labelX - link.labelW : link.labelX - link.labelW / 2;
  return (
    <>
      <BaseEdge id={id} path={link.path} markerEnd={markerEnd} style={style} />
      {link.label && (
        <EdgeLabelRenderer>
          <span
            className={`timeline-label absolute rounded-full bg-background-2 text-center text-11 leading-[14px] whitespace-nowrap text-foreground-muted ${link.spoiler ? "timeline-spoiler" : ""} ${dim ? "is-dim" : ""}`}
            style={{ width: link.labelW, transform: `translate(${left}px, ${link.labelY - 7}px)` }}
            aria-hidden="true"
          >
            {link.label}
          </span>
        </EdgeLabelRenderer>
      )}
    </>
  );
};

const nodeTypes = { work: WorkCard, box: Box, tick: Tick };

// lineage is a work and every work before or after it along the links, as
// the works up to it and the works down from it, each including it
const lineage = (id: string, links: Props["links"], spoilers: boolean): [Set<string>, Set<string>] => {
  const shown = (links ?? []).filter((l) => spoilers || !l.spoiler);
  const walk = (from: "source" | "target", to: "source" | "target"): Set<string> => {
    const seen = new Set([id]);
    const next = [id];
    for (let at = next.pop(); at !== undefined; at = next.pop()) {
      for (const l of shown) {
        const other = String(l[to]);
        if (String(l[from]) === at && !seen.has(other)) {
          seen.add(other);
          next.push(other);
        }
      }
    }
    return seen;
  };
  return [walk("target", "source"), walk("source", "target")];
};
const edgeTypes = { link: Link };

const Timeline = ({ width, height, ticks, boxes, works, links }: Props): ReactNode => {
  const [lit, setLit] = useState<{ id: string; spoilers: boolean } | null>(null);
  const [root, setRoot] = useState<HTMLElement | null>(null);
  // spoilers count only once they're shown, so a highlight can't give one away
  const light = (id: string | null): void => {
    const reveal = root?.closest(".timeline")?.querySelector<HTMLInputElement>(".timeline-reveal");
    setLit(id === null ? null : { id, spoilers: reveal?.checked ?? true });
  };
  const [up, down] = lit === null ? [null, null] : lineage(lit.id, links, lit.spoilers);
  const onLine = (a: string, b: string): boolean => up === null || down === null || (up.has(a) && up.has(b)) || (down.has(a) && down.has(b));

  const nodes: Node[] = [
    ...ticks.map((t): TickNode => ({ id: `tick${t.x}`, type: "tick", position: { x: t.x - 20, y: 0 }, width: 40, height: height - 4, data: { label: t.label }, zIndex: -2 })),
    ...(boxes ?? []).map((b, i): BoxNode => ({ id: `box${i}`, type: "box", position: { x: b.x, y: b.y }, width: b.w, height: b.h, data: { name: b.name } })),
    ...works.map(
      (w): WorkNode => ({
        id: String(w.id),
        type: "work",
        position: { x: w.card.x, y: w.card.y },
        width: w.card.w,
        height: w.card.h,
        className: onLine(String(w.id), String(w.id)) ? "" : "is-dim",
        data: { work: w, light },
      }),
    ),
  ];
  for (const n of nodes) {
    n.draggable = false;
    n.selectable = false;
    n.focusable = false;
  }

  const edges = (links ?? []).map((l, i): LinkEdge => {
    const id = `link${i}`;
    const dim = !onLine(String(l.source), String(l.target));
    return {
      id,
      type: "link",
      source: String(l.source),
      target: String(l.target),
      className: [l.sequel ? "timeline-sequel" : "timeline-other", l.spoiler ? "timeline-spoiler" : "", dim ? "is-dim" : ""].join(" "),
      markerEnd: l.sequel ? { type: MarkerType.ArrowClosed, width: 14, height: 14, color: dim ? "var(--color-border)" : "var(--color-foreground-muted)" } : undefined,
      // lit lines draw over dim ones where they share a trunk, and both under
      // the box names and labels
      zIndex: dim ? -1 : 0,
      data: { link: l, dim },
      focusable: false,
      selectable: false,
    };
  });

  return (
    <div ref={setRoot} className="overflow-x-auto" tabIndex={0} role="region" aria-label="timeline, scrolls sideways">
      <div className="mx-auto" style={{ width, height }}>
        <ReactFlow
          nodes={nodes}
          edges={edges}
          nodeTypes={nodeTypes}
          edgeTypes={edgeTypes}
          defaultViewport={{ x: 0, y: 0, zoom: 1 }}
          minZoom={1}
          maxZoom={1}
          panOnDrag={false}
          panOnScroll={false}
          zoomOnScroll={false}
          zoomOnPinch={false}
          zoomOnDoubleClick={false}
          preventScrolling={false}
          nodesDraggable={false}
          nodesConnectable={false}
          nodesFocusable={false}
          edgesFocusable={false}
          elementsSelectable={false}
          onNodeMouseEnter={(_, n) => n.type === "work" && light(n.id)}
          onNodeMouseLeave={() => light(null)}
          proOptions={{ hideAttribution: true }}
        />
      </div>
    </div>
  );
};

export const mount = island("timeline", Props, Timeline);
