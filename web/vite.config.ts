import tailwindcss from "@tailwindcss/vite";
import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";

// Backend integration, no index.html: Go finds the two entries' hashed files
// through the manifest.
export default defineConfig({
  plugins: [react(), tailwindcss()],
  build: {
    // a subdirectory, so emptying it keeps dist/.gitkeep for the Go embed
    outDir: "dist/bundle",
    manifest: true,
    rollupOptions: {
      input: ["src/app.css", "src/islands.ts"],
    },
  },
});
