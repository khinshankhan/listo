import {
  BaseEdge,
  EdgeLabelRenderer,
  getSmoothStepPath,
  Handle,
  MarkerType,
  Position,
  ReactFlow,
  type Edge,
  type EdgeProps,
  type Node,
  type NodeProps,
} from "@xyflow/react";
import "@xyflow/react/dist/base.css";
import { useState, type ReactNode } from "react";
import * as z from "zod/mini";
import { island } from "../mount";
import "./timeline.css";

// The layout is the server's, the same one its SVG fallback draws: x is time,
// y the row each work landed on.
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
      left: z.number(),
      right: z.number(),
      y: z.number(),
      dots: z.array(z.object({ x: z.number(), label: z.string() })),
    }),
  ),
  links: z.nullable(
    z.array(
      z.object({
        source: z.number(),
        sourceDot: z.number(),
        target: z.number(),
        targetDot: z.number(),
        label: z.string(),
        sequel: z.boolean(),
        spoiler: z.boolean(),
      }),
    ),
  ),
});
type Props = z.infer<typeof Props>;
type Work = Props["works"][number];

// a card around a work: padding either side of its stretch of row, and from
// above its name to below its dots' labels
const pad = 10;
const above = 32;
const cardHeight = 58;
// links to other rows hang off this far into a card, where the cards after
// it in time can't cover them
const trunk = 8;
// a line meeting a card from above or below turns this far short of it
const gap = 10;
// a step line leaves and meets a card this straight, and turns this round
const offset = 14;
const radius = 10;

type WorkNode = Node<{ work: Work; light: (id: string | null) => void }, "work">;
type BoxNode = Node<{ name: string }, "box">;
type TickNode = Node<{ label: string }, "tick">;
type LinkEdge = Edge<{ label: string; trunkX: number | null; turn: number; spoiler: boolean; dim: boolean }, "link">;

const WorkCard = ({ id, data: { work, light } }: NodeProps<WorkNode>): ReactNode => {
  const left = work.left - pad;
  const first = work.dots[0].x - left;
  const last = work.dots[work.dots.length - 1].x - left;
  return (
    <div className="timeline-card relative rounded-lg border border-border bg-background-2" style={{ width: work.right - work.left + 2 * pad, height: cardHeight }}>
      <a
        href={work.href}
        onFocus={() => light(id)}
        onBlur={() => light(null)}
        className="absolute top-[7px] right-[10px] left-[10px] truncate text-13 leading-4 font-semibold text-foreground-strong after:absolute after:inset-[-7px_-10px_-27px] hover:underline focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-stark-contrast"
      >
        {work.name}
      </a>
      {last > first && <span className="absolute h-[3px] -translate-y-1/2 rounded-full bg-foreground-strong" style={{ left: first, width: last - first, top: above }} aria-hidden="true" />}
      {work.dots.map((d, i) => {
        const x = d.x - left;
        return (
          <span key={i} aria-hidden="true">
            <span className="absolute size-[9px] -translate-1/2 rounded-full border-[1.5px] border-foreground-strong bg-background-2" style={{ left: x, top: above }} />
            <span className="absolute -translate-x-1/2 text-11 leading-[14px] whitespace-nowrap text-foreground-muted" style={{ left: x, top: above + 7 }}>
              {d.label}
            </span>
            <Handle type="source" id={`up${i}`} position={Position.Top} style={{ left: x }} isConnectable={false} />
            <Handle type="source" id={`down${i}`} position={Position.Bottom} style={{ left: x }} isConnectable={false} />
            <Handle type="target" id={`top${i}`} position={Position.Top} style={{ left: x }} isConnectable={false} />
            <Handle type="target" id={`bottom${i}`} position={Position.Bottom} style={{ left: x }} isConnectable={false} />
          </span>
        );
      })}
      <Handle type="source" id="down" position={Position.Bottom} style={{ left: trunk }} isConnectable={false} />
      <Handle type="source" id="rise" position={Position.Top} style={{ left: trunk }} isConnectable={false} />
      <Handle type="target" id="in" position={Position.Left} style={{ top: above }} isConnectable={false} />
      <Handle type="source" id="out" position={Position.Right} style={{ top: above }} isConnectable={false} />
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

// labels are about this wide a character, plus their padding
const labelWidth = (s: string): number => s.length * 6 + 12;

// Link is a step line between cards. Off a trunk, its label sits on the run
// into its target if it fits, else just left of the trunk; otherwise on the
// line's run across, turn pixels from the target's top or bottom.
const Link = ({ id, sourceX, sourceY, sourcePosition, targetX, targetY, targetPosition, markerEnd, style, data }: EdgeProps<LinkEdge>): ReactNode => {
  const turn = data?.turn ?? 0;
  const centerY = turn === 0 ? undefined : targetY + turn;
  const [path, midX, midY] = getSmoothStepPath({ sourceX, sourceY, sourcePosition, targetX, targetY, targetPosition, borderRadius: radius, offset, centerY });
  let at = `translate(-50%, -50%) translate(${centerY === undefined ? midX : (sourceX + targetX) / 2}px, ${centerY ?? midY}px)`;
  if (data?.trunkX != null) {
    const right = targetX - 14 - labelWidth(data.label) > data.trunkX + 8 ? targetX - 14 : data.trunkX - 6;
    at = `translate(-100%, -50%) translate(${right}px, ${centerY ?? targetY}px)`;
  }
  return (
    <>
      <BaseEdge id={id} path={path} markerEnd={markerEnd} style={style} />
      {data?.label && (
        <EdgeLabelRenderer>
          <span
            className={`timeline-label absolute rounded-full bg-background-2 px-1.5 text-11 leading-[14px] whitespace-nowrap text-foreground-muted ${data.spoiler ? "timeline-spoiler" : ""} ${data.dim ? "is-dim" : ""}`}
            style={{ transform: at }}
            aria-hidden="true"
          >
            {data.label}
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
  const byID = new Map(works.map((w) => [w.id, w]));

  const nodes: Node[] = [
    ...ticks.map((t): TickNode => ({ id: `tick${t.x}`, type: "tick", position: { x: t.x - 20, y: 0 }, width: 40, height: height - 4, data: { label: t.label }, zIndex: -2 })),
    ...(boxes ?? []).map((b, i): BoxNode => ({ id: `box${i}`, type: "box", position: { x: b.x, y: b.y }, width: b.w, height: b.h, data: { name: b.name } })),
    ...works.map(
      (w): WorkNode => ({
        id: String(w.id),
        type: "work",
        position: { x: w.left - pad, y: w.y - above },
        width: w.right - w.left + 2 * pad,
        height: cardHeight,
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

  // along a row a link runs card to card, arcing over the row unless it's a
  // sequel. To another row, a sequel, side story or spin-off hangs off its
  // card's trunk and runs in, or meets the card from above or below when it
  // starts too near the trunk; anything else leaves its release that way.
  const edges = (links ?? []).map((l, i): LinkEdge => {
    const from = byID.get(l.source);
    const to = byID.get(l.target);
    if (!from || !to) {
      throw new Error(`link ${i} joins a work that isn't drawn`);
    }
    const trunkX = from.left - pad + trunk;
    const down = from.y < to.y;
    let handles: [string, string];
    let turn = 0;
    if (from.y === to.y) {
      handles = l.sequel ? ["out", "in"] : [`up${l.sourceDot}`, `top${l.targetDot}`];
    } else if (l.sequel && to.left - pad > trunkX + offset + radius) {
      handles = [down ? "down" : "rise", "in"];
    } else {
      const source = l.sequel ? (down ? "down" : "rise") : `${down ? "down" : "up"}${l.sourceDot}`;
      handles = [source, `${down ? "top" : "bottom"}${l.targetDot}`];
      turn = down ? -gap : gap;
    }
    const id = `link${i}`;
    const dim = !onLine(String(l.source), String(l.target));
    return {
      id,
      type: "link",
      source: String(l.source),
      sourceHandle: handles[0],
      target: String(l.target),
      targetHandle: handles[1],
      className: [l.sequel ? "timeline-sequel" : "timeline-other", l.spoiler ? "timeline-spoiler" : "", dim ? "is-dim" : ""].join(" "),
      markerEnd: l.sequel ? { type: MarkerType.ArrowClosed, width: 14, height: 14, color: dim ? "var(--color-border)" : "var(--color-foreground-muted)" } : undefined,
      // lit lines draw over dim ones where they share a trunk, and both under
      // the box names and labels
      zIndex: dim ? -1 : 0,
      data: { label: l.label, trunkX: l.sequel && from.y !== to.y ? trunkX : null, turn, spoiler: l.spoiler, dim },
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
