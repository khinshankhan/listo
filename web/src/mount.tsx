import { Component, StrictMode, type FunctionComponent, type ReactNode } from "react";
import { createRoot } from "react-dom/client";

// What every island module exports.
export type Mount = (el: HTMLElement, props: unknown) => void;

// Anything that can turn unknown JSON into P or throw: a zod/mini schema here.
type Schema<P> = { parse: (raw: unknown) => P };

type BoundaryProps = { name: string; fallbackHTML: string; children: ReactNode };
type BoundaryState = { failed: boolean };

// If an island throws while rendering, put the server's fallback back.
class Boundary extends Component<BoundaryProps, BoundaryState> {
  state: BoundaryState = { failed: false };

  static getDerivedStateFromError(): BoundaryState {
    return { failed: true };
  }

  componentDidCatch(error: unknown): void {
    console.error(`island "${this.props.name}" crashed; showing its fallback`, error);
  }

  render(): ReactNode {
    if (this.state.failed) {
      // Markup this page's server rendered and escaped, not user input.
      return <div style={{ display: "contents" }} dangerouslySetInnerHTML={{ __html: this.props.fallbackHTML }} />;
    }
    return this.props.children;
  }
}

// island turns a component into a Mount: validate props, then replace the
// fallback.
export const island =
  <P extends object>(name: string, schema: Schema<P>, View: FunctionComponent<P>): Mount =>
  (el, raw) => {
    const props = schema.parse(raw);
    const fallbackHTML = el.innerHTML;
    createRoot(el).render(
      <StrictMode>
        <Boundary name={name} fallbackHTML={fallbackHTML}>
          <View {...props} />
        </Boundary>
      </StrictMode>,
    );
  };
