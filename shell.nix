{ pkgs ? import <nixpkgs> {} }:
pkgs.mkShell {
  buildInputs = with pkgs; [
    go 
    gotools 
    gopls 
    delve # Go debugger 

    xorg.libX11
    xorg.libXi
    xorg.libXext
  ];

  shellHook = ''
    export GOPATH=$HOME/go
    export PATH=$GOPATH/bin:$PATH 
    mkdir -p $GOPATH/src $GOPATH/bin $GOPATH/pkgs

    echo $(go version)

    zsh 
  '';
}

