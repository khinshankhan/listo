import { useState, type ReactNode } from "react";
import * as z from "zod/mini";
import { island } from "../mount";

const Props = z.object({ name: z.string() });
type Props = z.infer<typeof Props>;

// matches views.waveFrame, so mounting doesn't shift anything
const Wave = ({ name }: Props): ReactNode => {
  const [count, setCount] = useState<number>(0);

  return (
    <div className="flex h-10 items-center gap-4">
      <button
        type="button"
        className="inline-flex h-10 cursor-pointer items-center rounded-md bg-accent-4 px-4 text-14 font-semibold text-accent-12 hover:bg-accent-5 focus-visible:outline-2 focus-visible:outline-stark-contrast"
        onClick={() => setCount((c) => c + 1)}
      >
        wave
      </button>
      <span className="font-mono text-14 text-foreground-muted">
        {count === 1 ? "1 wave" : `${count} waves`}
        {count > 0 && ` at ${name}`}
      </span>
    </div>
  );
};

export const mount = island("wave", Props, Wave);
