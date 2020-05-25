{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = [
    pkgs.systemfd
    pkgs.go
  ];

  shellHook = ''
    unset GOPATH
  '';
}
