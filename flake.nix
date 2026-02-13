{
  description = "Go dev with Plan 9 mk";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
  };

  outputs = { self, nixpkgs }:
    let
      supportedSystems = [ "x86_64-linux" "aarch64-linux" ];
      forAllSystems = f: nixpkgs.lib.genAttrs supportedSystems (system: f {
        pkgs = import nixpkgs { inherit system; };
      });
    in
    {
      devShells = forAllSystems ({ pkgs, ... }:
        {
          default = pkgs.mkShell {
            buildInputs = with pkgs; [
              go
              gopls
              gotools
              golangci-lint
              plan9port
            ];

            shellHook = ''
              export PLAN9="${pkgs.plan9port}/plan9"
              export PATH="$PATH:$PLAN9/bin"
              echo "Plan 9: user space loaded at '$PLAN9'"
              echo "Go: $(go version)"
            '';
          };
        });
    };
}
