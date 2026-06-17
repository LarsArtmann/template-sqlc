{
  description = "SQLC project template for Go";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-parts = {
      url = "github:hercules-ci/flake-parts";
      inputs.nixpkgs-lib.follows = "nixpkgs";
    };
    systems.url = "github:nix-systems/default";
    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    inputs@{
      self,
      flake-parts,
      systems,
      treefmt-nix,
      ...
    }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = import systems;

      imports = [
        treefmt-nix.flakeModule
      ];

      perSystem =
        {
          config,
          pkgs,
          ...
        }:
        {
          treefmt = {
            projectRootFile = "go.mod";
            programs = {
              gofumpt.enable = true;
              goimports.enable = true;
              nixfmt.enable = true;
            };
          };

          checks.format = config.treefmt.build.check self;

          devShells = {
            default = pkgs.mkShell {
              name = "template-sqlc-dev";

              packages = [
                pkgs.go_1_26
                pkgs.golangci-lint
                pkgs.gopls
                pkgs.delve
                pkgs.gotools
                pkgs.gofumpt
                pkgs.sqlc
              ];

              GOWORK = "off";
              GOPRIVATE = "github.com/LarsArtmann/*,github.com/larsartmann/*";
            };

            ci = pkgs.mkShellNoCC {
              packages = [
                pkgs.go_1_26
                pkgs.golangci-lint
              ];

              GOWORK = "off";
              GOPRIVATE = "github.com/LarsArtmann/*,github.com/larsartmann/*";
            };
          };
        };
    };
}
