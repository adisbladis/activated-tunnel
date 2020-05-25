{ pkgs ? import <nixpkgs> {} }:

let
  inherit (pkgs) lib;

in pkgs.buildGoPackage {
  pname = "activated-tunnel";
  version = "0.1";

  goPackagePath = "github.com/adisbladis/activated-tunnel";

  src = lib.cleanSource ./.;

  goDeps = ./deps.nix;

  meta = {
    description = "Socket activated SSH tunnels";
    homepage = https://github.com/adisbladis/activated-tunnel;
    license = lib.licenses.mit;
    maintainers =  [ lib.maintainers.adisbladis ];
  };
}
