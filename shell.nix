{ pkgs ? import <nixpkgs> { } }:

pkgs.mkShell {
  buildInputs = with pkgs; [
    air
    tailwindcss
    tailwindcss-language-server
  ];
}

