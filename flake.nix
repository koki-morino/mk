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
              plan9port
            ];

            shellHook = ''
              export PS1="\[\033[1;32m\][nix-dev:\w]\$\[\033[0m\] "
              alias mk='${pkgs.plan9port}/plan9/bin/mk'
              echo "Plan 9 user space loaded at '${pkgs.plan9port}'"
            '';
          };
        });
    };
}
