{ pkgs ? import <nixpkgs> { } }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    air
    tailwindcss
    tailwindcss-language-server
  ];

  shellHook = ''
    set -a
    source env.sh
    set +a
  '';
}

