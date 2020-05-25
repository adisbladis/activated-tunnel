{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = [
    pkgs.systemfd
    pkgs.vgo2nix
    pkgs.go
  ];

  shellHook = ''
    unset GOPATH
  '';
}
