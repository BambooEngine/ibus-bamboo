{
  description = "A flake for ibus-bamboo";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
  };

  outputs = { self, nixpkgs }:
    let
      version = "v0.8.4";

      supportedSystems = [ "x86_64-linux" "aarch64-linux" ];

      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;

      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });
    in
    {
      packages = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          default = pkgs.stdenv.mkDerivation {
            pname = "ibus-bamboo";
            inherit version;

            src = ./.;

            nativeBuildInputs = [
              pkgs.pkg-config
              pkgs.wrapGAppsHook
              pkgs.git
            ];

            buildInputs = [
              pkgs.go
              pkgs.xorg.libXtst
		          pkgs.gtk4
		          pkgs.glib
		          pkgs.librsvg
		          pkgs.gdk-pixbuf
		          pkgs.gobject-introspection
		          pkgs.hicolor-icon-theme
            ];

            preConfigure = ''
              export GOCACHE="$TMPDIR/go-cache"
              sed -i "s,/usr,$out," data/bamboo.xml
            '';

            makeFlags = [
              "PREFIX=${placeholder "out"}"
            ];

            meta = {
              isIbusEngine = true;
            };
          };
        }
      );

      devShells = forAllSystems (system:
        let
            pkgs = nixpkgsFor.${system};
        in
        {
          default = pkgs.mkShell {
            nativeBuildInputs = [
              pkgs.pkg-config
              pkgs.wrapGAppsHook3
              pkgs.go
              pkgs.gopls
            ];

            buildInputs = [
              pkgs.xorg.libXtst
              pkgs.gtk4
              pkgs.graphene
              pkgs.glib
              pkgs.gobject-introspection
            ];
          };
        }
      );
    };
}
