{
  description = "listo -- the 1 list for all your content";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";

  outputs = { nixpkgs, ... }:
    let
      systems = [ "x86_64-linux" "aarch64-linux" "aarch64-darwin" ];
      forAll = f: nixpkgs.lib.genAttrs systems (system: f nixpkgs.legacyPackages.${system});
    in {
      # go, node and pnpm come from the machine; this pins the rest
      devShells = forAll (pkgs: {
        default = pkgs.mkShellNoCC {
          packages = [
            # 18 for uuidv7(); ships pg_trgm and unaccent
            pkgs.postgresql_18
            # needs cgo, so it can't be a go.mod tool like templ
            pkgs.sqlc
          ];
        };
      });
    };
}
